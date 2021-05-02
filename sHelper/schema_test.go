package sHelper

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type TestSchema struct {
	BaseSchema
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
	Field4 string
	Items  []string `json:"items"`
}

type TestComplexSchema struct {
	TestSchema
	Nested TestSchema `json:"nested"`
}

type TestInvalidSchema struct {
	TestSchema
	Parent TestSchema
}

func TestSchemaCreate(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "")
}

func TestSchemaInvalidJson(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.go",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "207110: schema_test.go couldn't be parsed - Invalid JSON error")
}

func TestSchemaInvalidPath(t *testing.T) {
	absPath, _ := filepath.Abs("foo.json")
	schema := Schema{
		FileName:  "foo.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "209010: foo.json file was not found. Message return: "+absPath)
}

func TestSchemaReqFields(t *testing.T) {
	type TestSchema struct {
		Field1 string `json:"field1"`
	}
	schema := Schema{
		StructRef: &TestSchema{},
	}
	json.Unmarshal([]byte("{\"required\": [\"field1\", \"field2\"], \"properties\": {\"field1\": {\"$id\": \"#/properties/field1\"}}}"), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "109999: #/properties/field2 was/were not found")
}

func TestSchemaMissingFields(t *testing.T) {
	type TestSchema struct {
		Field1 string `json:"field1"`
	}
	schema := Schema{
		StructRef: &TestSchema{},
	}
	json.Unmarshal([]byte("{\"properties\": {\"field1\": {\"$id\": \"#/properties/field1\"},\"field2\": {\"$id\": \"#/properties/field2\"}}}"), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "200513: field2 (#/properties/field2) must be populated")
}

func TestSchemaIsKind(t *testing.T) {
	AssertEqual(t, isKind(jsonKinds["integer"], reflect.ValueOf(1).Type()), true)
	AssertEqual(t, isKind(jsonKinds["integer"], reflect.ValueOf(true).Type()), false)
}

func TestSchemaMissingParameters(t *testing.T) {
	schema := Schema{
		StructRef: &TestSchema{},
	}
	json.Unmarshal([]byte("{\"properties\": {\"field0\": {\"$id\": \"#/properties/field0\"}}}"), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "200513: field0 (#/properties/field0) must be populated")
}

func TestSchemaInvalidTypes(t *testing.T) {
	schema := Schema{
		StructRef: &TestSchema{},
	}
	json.Unmarshal([]byte("{\"properties\": {\"field1\": {\"$id\": \"#/properties/field1\", \"type\": \"int\"}}}"), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "200200: [Field1('#/properties/field1')] must be of type [int]")
}

func TestSchemaParseInvalidType(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "")
	soteErr = schema.Parse([]byte("{\"field1\": \"Hello\"}"), &TestComplexSchema{})
	AssertEqual(t, soteErr.FmtErrMsg, "200200: *sHelper.TestComplexSchema must be of type *sHelper.TestSchema")
}

func TestSchemaParseInvalidJson(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "")
	body := TestSchema{}
	soteErr = schema.Parse([]byte("{int}"), &body)
	AssertEqual(t, soteErr.FmtErrMsg, "207110: Body couldn't be parsed - Invalid JSON error")
}

func TestSchemaParse(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "")
	body := TestSchema{}
	soteErr = schema.Parse([]byte("{\"field1\": \"Hello\", \"field2\": \"World\"}"), &body)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, body.Field1, "Hello")
	AssertEqual(t, body.Field2, "World")
	AssertEqual(t, body.Field3, "VALUE1")
}

func TestSchemaParseReqFields(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "")
	body := TestSchema{}
	soteErr = schema.Parse([]byte("{\"field1\": \"Hello\"}"), &body)
	AssertEqual(t, soteErr.FmtErrMsg, "206200: Message doesn't match signature. Sender must provide the following parameter names: #/properties/field2")
	AssertEqual(t, body.Field1, "Hello")
	AssertEqual(t, body.Field2, "")
	AssertEqual(t, body.Field3, "VALUE1")
}

func TestSchemaParseInvalidEnum(t *testing.T) {
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := schema.Validate()
	AssertEqual(t, soteErr.FmtErrMsg, "")
	body := TestSchema{}
	soteErr = schema.Parse([]byte("{\"field1\": \"Test\", \"field2\": \"Hello\", \"field3\": \"World\"}"), &body)
	AssertEqual(t, soteErr.FmtErrMsg, "200250: #/properties/field3 (World) must contain one of these values: [VALUE1 VALUE2]")
	AssertEqual(t, body.Field1, "Test")
	AssertEqual(t, body.Field2, "Hello")
	AssertEqual(t, body.Field3, "World")
}

func TestSchemaParseInvalidEnumParam(t *testing.T) {
	schema := Schema{
		StructRef: &TestSchema{},
	}
	json.Unmarshal([]byte(`
	{
		"properties": {
			"items": {
				"$id": "#/properties/items",
				"type": "array",
				"enum": ["VALUE1", "VALUE2"]
			}
		}
	}
	`), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "")
	body := TestSchema{}
	soteErr := schema.Parse([]byte("{\"items\": [\"VALUE3\", \"VALUE1\", \"VALUE2\"]}"), &body)
	AssertEqual(t, soteErr.FmtErrMsg, "200250: #/properties/items ([VALUE3 VALUE1 VALUE2]) must contain one of these values: [VALUE1 VALUE2]")
}

func TestSchemaParseNestedStruct(t *testing.T) {
	schema := Schema{
		StructRef: &TestComplexSchema{},
	}
	json.Unmarshal([]byte(`
	{
		"required": [
			"field1",
			"nested"
		],
		"properties": {
			"field1": {
				"$id": "#/properties/field1",
				"type": "string"
			},
			"nested": {
				"$id": "#/properties/nested",
				"type": "object",
				"required": [
					"field2"
				],
				"properties": {
					"field2" : {
						"$id": "#/properties/nested/properties/field2",
						"type": "string",
						"default": "MyValue"
					},
					"items": {
						"$id": "#/properties/nested/properties/items",
						"type": "array",
						"enum": ["VALUE1", "VALUE2"]
					}
				}
			}
		}
	}
	`), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "")

	body := TestComplexSchema{}
	soteErr := schema.Parse([]byte(`
	{
		"field1": "Hello",
		"nested": {
			"field1": "World",
			"items": ["VALUE1", "VALUE2"]
		}
	}
	`), &body)

	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, body.Field1, "Hello")
	AssertEqual(t, body.Field2, "")
	AssertEqual(t, body.Nested.Field1, "World")
	AssertEqual(t, body.Nested.Field2, "MyValue")
	AssertEqual(t, body.Nested.Field3, "")
	AssertEqual(t, strings.Join(body.Nested.Items, ", "), strings.Join([]string{"VALUE1", "VALUE2"}, ", "))
}

func TestSchemaParseMissingField(t *testing.T) {
	schema := Schema{
		StructRef: &TestInvalidSchema{},
	}
	json.Unmarshal([]byte(`
	{
		"properties": {
			"parent": {
				"$id": "#/properties/parent",
				"type": "object",
				"properties": {
					"field4": {
						"$id": "#/properties/nested/properties/field4",
						"type": "string"
					}
				}
			}
		}
	}
	`), &schema.jsonSchema)
	AssertEqual(t, schema.validateSchema().FmtErrMsg, "200513: parent (#/properties/parent), field4 (#/properties/nested/properties/field4) must be populated")
	body := TestInvalidSchema{}
	soteErr := schema.Parse([]byte("{\"parent\": {\"field4\": \"Hello\"}}"), &body)

	val := reflect.ValueOf(body)
	f, _ := val.Type().FieldByName("Field4")
	v := findField(val, &f, schema.jsonSchema.Properties["parent"].Properties["field4"], 1)

	AssertEqual(t, v == nil, true)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, body.Parent.Field4, "Hello")
}
