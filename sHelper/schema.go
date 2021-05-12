package sHelper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type BaseSchema struct {
	JsonWebToken    string   `json:"json-web-token"`
	MessageId       string   `json:"message-id"`
	AwsUserName     string   `json:"aws-user-name"`
	OrganizationId  int      `json:"organizations-id"`
	ClientCompanyId int      `json:"client-company-id"` //optional
	RoleList        []string `json:"role-list"`         //optional
}

type Schema struct {
	FileName       string
	StructRef      interface{}
	structType     reflect.Type
	defaultFields  map[string]*jsonProperty
	enumFields     map[string]*jsonProperty
	requiredFields map[string]*jsonProperty
	jsonFields     map[string]*reflect.StructField
	jsonSchema     jsonSchema
}

type jsonSchema struct {
	Required   []string
	Properties map[string]*jsonProperty
}

type jsonProperty struct {
	Id         string `json:"$id"`
	Default    interface{}
	Enum       []interface{}
	Type       string
	Required   []string
	Properties map[string]*jsonProperty
}

var (
	missingParameters []string
	requiredFields    []string
	invalidFileds     []string
	invalidTypes      []string
	jsonKinds         = map[string][]reflect.Kind{
		"array":   {reflect.Slice, reflect.Array},
		"boolean": {reflect.Bool},
		"string":  {reflect.String},
		"object":  {reflect.Interface, reflect.Map, reflect.Struct},
		"number":  {reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128},
		"integer": {reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	}
)

func isKind(kinds []reflect.Kind, kind reflect.Type) bool {
	for _, k := range kinds {
		if k == kind.Kind() {
			return true
		}
	}
	return false
}

func find(s *Schema, val reflect.Value, propLevel string) {
	e := val.Elem()
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		var jsonTag string
		if jsonTag = f.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			s.jsonFields[propLevel+"/"+jsonTag] = &f
		}

		if f.Type.Kind() == reflect.Struct {
			v := e.Field(i).Addr().Interface()
			if jsonTag != "" && jsonTag != "-" {
				find(s, reflect.ValueOf(v), fmt.Sprintf("%v/%v/properties", propLevel, jsonTag))
			} else {
				find(s, reflect.ValueOf(v), propLevel) // injected struct
			}
		}
	}
}

func propValidation(s *Schema, propLevel string, props map[string]*jsonProperty, required []string) {
	for _, n := range required {
		id := propLevel + "/" + n
		f := s.jsonFields[id]
		if f == nil || props[n] == nil {
			requiredFields = append(requiredFields, id)
		} else {
			s.requiredFields[id] = props[n]
		}
	}
	for id, prop := range props {
		if !(prop.Default == nil || prop.Default == "") {
			s.defaultFields[propLevel+"/"+id] = prop
		}
		if s.jsonFields[prop.Id] == nil {
			missingParameters = append(missingParameters, fmt.Sprintf("%v (%v)", id, prop.Id))
		} else {
			f := s.jsonFields[prop.Id] //stuct field type
			k := jsonKinds[prop.Type]  // kind(s)
			if !isKind(k, f.Type) {
				invalidFileds = append(invalidFileds, fmt.Sprintf("%s('%s')", f.Name, prop.Id))
				invalidTypes = append(invalidTypes, prop.Type)
			}
			// save enum types
			if prop.Enum != nil && len(prop.Enum) > 0 {
				s.enumFields[prop.Id] = prop
			}
		}
		propValidation(s, prop.Id+"/properties", prop.Properties, prop.Required)
	}
}

func findField(v reflect.Value, f *reflect.StructField, prop *jsonProperty, level int) *reflect.Value {
	levels := strings.Split(prop.Id, "/properties/")
	l := len(levels) - 1 // skip prefix '#'
	n := levels[level]

	t := v.Type()
	if l == level {
		fn := v.FieldByName(f.Name)
		return &fn
	} else {
		if t.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				fv := v.Field(i)
				ft := v.Type().Field(i)
				if fv.Kind() == reflect.Struct && (ft.Tag.Get("json") == n) {
					return findField(fv, f, prop, level+1)
				}
			}
		}
	}
	return nil
}

func (s *Schema) Validate() (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	absPath, _ := filepath.Abs(s.FileName)
	//Parse file (find required fields)
	if _, err := os.Stat(absPath); err != nil {
		soteErr = NewError().FileNotFound(s.FileName, absPath)
	} else {
		plan, _ := ioutil.ReadFile(s.FileName)
		err := json.Unmarshal(plan, &s.jsonSchema)
		if err != nil {
			soteErr = NewError().InvalidJson(s.FileName)
		} else {
			soteErr = s.validateSchema()
		}
	}
	return
}

func (s *Schema) validateSchema() (soteErr sError.SoteError) {
	missingParameters = []string{}
	requiredFields = []string{}
	invalidFileds = []string{}
	invalidTypes = []string{}

	s.jsonFields = make(map[string]*reflect.StructField)
	s.requiredFields = make(map[string]*jsonProperty)
	s.defaultFields = make(map[string]*jsonProperty)
	s.enumFields = make(map[string]*jsonProperty)

	val := reflect.ValueOf(s.StructRef)
	s.structType = val.Type()

	//Parse Struct (validate required fields)
	find(s, val, "#/properties")

	// Validate missing/required fields and datatype
	propValidation(s, "#/properties", s.jsonSchema.Properties, s.jsonSchema.Required)

	if len(requiredFields) > 0 {
		soteErr = NewError().ItemNotFound(strings.Join(requiredFields, ", "))
	} else if len(missingParameters) > 0 {
		soteErr = NewError().MustBePopulated(strings.Join(missingParameters, ", "))
	} else if len(invalidTypes) > 0 {
		soteErr = NewError().MustBeType(
			fmt.Sprintf("[%s]", strings.Join(invalidFileds, ", ")),
			fmt.Sprintf("[%s]", strings.Join(invalidTypes, ", ")))
	}

	return
}

func (s *Schema) Parse(data []byte, body interface{}) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	b := reflect.ValueOf(body)
	if s.structType != b.Type() {
		soteErr = NewError().MustBeType(b.Type().String(), s.structType)
	} else {
		err := json.Unmarshal(data, &body)
		emptyFields := []string{}
		if err != nil {
			soteErr = NewError().InvalidJson("Body")
		} else {
			elem := b.Elem()

			//assign default values
			for id, prop := range s.defaultFields {
				f := s.jsonFields[id]
				nf := findField(elem, f, prop, 1)
				if nf != nil && nf.IsZero() {
					v := reflect.ValueOf(prop.Default)
					nf.Set(v.Convert(nf.Type()))
				}
			}

			//validate required fields
			for id, prop := range s.requiredFields {
				f := s.jsonFields[id]
				nf := findField(elem, f, prop, 1)
				if nf == nil || nf.IsZero() {
					emptyFields = append(emptyFields, id)
				}
			}
			if soteErr.ErrCode == nil && len(emptyFields) > 0 {
				soteErr = NewError().InvalidParameters(strings.Join(emptyFields, ", "))
			}

			if soteErr.ErrCode == nil {
				//validate enum fields
				for id, prop := range s.enumFields {
					f := s.jsonFields[id]
					nf := findField(elem, f, prop, 1)
					v := nf.Interface()
					found := false
					if prop.Type == "array" {
						for j := 0; j < nf.Len(); j++ {
							found = false
							for _, i := range prop.Enum {
								if nf.Index(j).Interface() == i {
									found = true
									break
								}
							}
							if !found {
								break
							}
						}
					} else {
						for _, i := range prop.Enum {
							if i == v {
								found = true
								break
							}
						}
					}
					if !found {
						soteErr = NewError().AllowValues(id, v, prop.Enum)
						break
					}
				}
			}
		}
	}
	return
}
