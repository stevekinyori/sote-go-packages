package sError

import (
	"testing"
)

func TestErrNotNumeric(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotNumeric(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotNumeric() constructor returned nil error")
			}
		})
	}
}

func TestErrNotString(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotString(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotString() constructor returned nil error")
			}
		})
	}
}

func TestErrTooSmall(t *testing.T) {
	type args struct {
		fieldName interface{}
		minSize   interface{}
		operation string
		Err       error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrTooSmall(tt.args.fieldName, tt.args.minSize, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrTooSmall() constructor returned nil error")
			}
		})
	}
}

func TestErrNotFloat(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotFloat(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotFloat() constructor returned nil error")
			}
		})
	}
}

func TestErrNotArray(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotArray(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotArray() constructor returned nil error")
			}
		})
	}
}

func TestErrNotJSONString(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotJSONString(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotJSONString() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidEmail(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidEmail(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidEmail() constructor returned nil error")
			}
		})
	}
}

func TestErrNotDate(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotDate(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotDate() constructor returned nil error")
			}
		})
	}
}

func TestErrNotTimestamp(t *testing.T) {
	type args struct {
		fieldName  interface{}
		fieldValue interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotTimestamp(tt.args.fieldName, tt.args.fieldValue, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotTimestamp() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidSize(t *testing.T) {
	type args struct {
		fieldName    interface{}
		fieldValue   interface{}
		relativeSize interface{}
		expectedSize interface{}
		actualSize   interface{}
		operation    string
		Err          error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidSize(tt.args.fieldName, tt.args.fieldValue, tt.args.relativeSize, tt.args.expectedSize, tt.args.actualSize, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidSize() constructor returned nil error")
			}
		})
	}
}

func TestErrJsonConversionError(t *testing.T) {
	type args struct {
		structName interface{}
		structType interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrJsonConversionError(tt.args.structName, tt.args.structType, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrJsonConversionError() constructor returned nil error")
			}
		})
	}
}

func TestErrNotMap(t *testing.T) {
	type args struct {
		paramName interface{}
		operation string
		Err       error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotMap(tt.args.paramName, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotMap() constructor returned nil error")
			}
		})
	}
}

func TestErrMissingErrorNumber(t *testing.T) {
	type args struct {
		errorMsgNumber interface{}
		operation      string
		Err            error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrMissingErrorNumber(tt.args.errorMsgNumber, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrMissingErrorNumber() constructor returned nil error")
			}
		})
	}
}
