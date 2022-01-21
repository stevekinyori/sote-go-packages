package sError

import (
	"testing"
)

func TestErrDirtyRead(t *testing.T) {
	type args struct {
		operation string
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				operation: "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrDirtyRead(tt.args.operation, tt.args.err)
			if err == nil {
				t.Errorf("ErrDirtyRead() constructor returned nil error")
			}
		})
	}
}

func TestErrCancelledComplete(t *testing.T) {
	type args struct {
		item      interface{}
		operation string
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				item:      "test-item",
				operation: "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrCancelledComplete(tt.args.item, tt.args.operation, tt.args.err)
			if err == nil {
				t.Errorf("ErrCancelledComplete() constructor returned nil error")
			}
		})
	}
}

func TestErrInactiveItem(t *testing.T) {
	type args struct {
		item      interface{}
		operation string
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				item:      "test-item",
				operation: "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInactiveItem(tt.args.item, tt.args.operation, tt.args.err)
			if err == nil {
				t.Errorf("ErrInactiveItem() constructor returned nil error")
			}
		})
	}
}

func TestErrTimeOut(t *testing.T) {
	type args struct {
		service   interface{}
		operation string
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				service:   "test-service",
				operation: "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrTimeOut(tt.args.service, tt.args.operation, tt.args.err)
			if err == nil {
				t.Errorf("ErrTimeOut() constructor returned nil error")
			}
		})
	}
}

func TestErrTableDoesNotExist(t *testing.T) {
	type args struct {
		operation string
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				operation: "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrTableDoesNotExist(tt.args.operation, tt.args.err)
			if err == nil {
				t.Errorf("ErrTableDoesNotExist() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidDataType(t *testing.T) {
	type args struct {
		expectedType interface{}
		param        interface{}
		operation    string
		err          error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				expectedType: "test-expected-type",
				param:        "test-param",
				operation:    "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidDataType(tt.args.expectedType, tt.args.param, tt.args.operation, tt.args.err)
			if err == nil {
				t.Errorf("ErrInvalidDataType() constructor returned nil error")
			}
		})
	}
}

func TestErrMissingRequiredValue(t *testing.T) {
	type args struct {
		paramName           interface{}
		paramValue          interface{}
		listofValuesAllowed interface{}
		err                 error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				paramName:           "test-param-name",
				paramValue:          "test-param-value",
				listofValuesAllowed: []string{"list-of-values-allowed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrMissingRequiredValue(
				tt.args.paramName, tt.args.paramValue, tt.args.listofValuesAllowed, tt.args.err)
			if err == nil {
				t.Errorf("ErrMissingRequiredValue() constructor returned nil error")
			}
		})
	}
}

func TestErrLinkedParameterValueMissing(t *testing.T) {
	type args struct {
		requiredParam    interface{}
		linkedParam      interface{}
		linkedParamValue interface{}
		err              error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				requiredParam:    "test-required-param",
				linkedParam:      "test-linked-param",
				linkedParamValue: "test-linked-param-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrLinkedParameterValueMissing(
				tt.args.requiredParam, tt.args.linkedParam, tt.args.linkedParamValue, tt.args.err)
			if err == nil {
				t.Errorf("ErrLinkedParameterValueMissing() constructor returned nil error")
			}
		})
	}
}

func TestErrParameterLockOtherParameterSet(t *testing.T) {
	type args struct {
		lockedParam     interface{}
		otherParam      interface{}
		otherParamValue interface{}
		err             error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				lockedParam:     "test-locked-param",
				otherParam:      "test-other-param",
				otherParamValue: "test-other-param-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrParameterLockOtherParameterSet(
				tt.args.lockedParam, tt.args.otherParam, tt.args.otherParamValue, tt.args.err)
			if err == nil {
				t.Errorf("ErrParameterLockOtherParameterSet() constructor returned nil error")
			}
		})
	}
}

func TestErrLinkedParamsMustBothBeSetOrNull(t *testing.T) {
	type args struct {
		paramName   interface{}
		linkedParam interface{}
		err         error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				paramName:   "test-param-name",
				linkedParam: "test-linked-param",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrLinkedParamsMustBothBeSetOrNull(tt.args.paramName, tt.args.linkedParam, tt.args.err)
			if err == nil {
				t.Errorf("ErrLinkedParamsMustBothBeSetOrNull() constructor returned nil error")
			}
		})
	}
}

func TestErrParamsMustBeProvided(t *testing.T) {
	type args struct {
		params interface{}
		err    error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				params: []string{"test-param"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrParamsMustBeProvided(tt.args.params, tt.args.err)
			if err == nil {
				t.Errorf("ErrParamsMustBeProvided() constructor returned nil error")
			}
		})
	}
}

func TestErrParamMustBeSet(t *testing.T) {
	type args struct {
		param interface{}
		err   error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				param: "test-param",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrParamMustBeSet(tt.args.param, tt.args.err)
			if err == nil {
				t.Errorf("ErrParamMustBeSet() constructor returned nil error")
			}
		})
	}
}

func TestErrParameterMustBeEmptyWhenParameterSet(t *testing.T) {
	type args struct {
		param       interface{}
		linkedParam interface{}
		err         error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				param:       "test-param",
				linkedParam: "test-linked-param",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrParameterMustBeEmptyWhenParameterSet(tt.args.param, tt.args.linkedParam, tt.args.err)
			if err == nil {
				t.Errorf("ErrParameterMustBeEmptyWhenParameterSet() constructor returned nil error")
			}
		})
	}
}

func TestErrBadHTTPRequest(t *testing.T) {
	type args struct {
		req interface{}
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				req: "test-http-request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrBadHTTPRequest(tt.args.req, tt.args.err)
			if err == nil {
				t.Errorf("ErrBadHTTPRequest() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidAPIEnvironment(t *testing.T) {
	type args struct {
		environmentName interface{}
		err             error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				environmentName: "test-environment-name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidAPIEnvironment(tt.args.environmentName, tt.args.err)
			if err == nil {
				t.Errorf("ErrInvalidAPIEnvironment() constructor returned nil error")
			}
		})
	}
}

func TestErrQuickSight(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrQuickSight(tt.args.err)
			if err == nil {
				t.Errorf("ErrQuickSight() constructor returned nil error")
			}
		})
	}
}

func TestErrDatabaseError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrDatabaseError(tt.args.err)
			if err == nil {
				t.Errorf("ErrDatabaseError() constructor returned nil error")
			}
		})
	}
}

func TestErrSQLError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrSQLError(tt.args.err)
			if err == nil {
				t.Errorf("ErrSQLError() constructor returned nil error")
			}
		})
	}
}

func TestErrCognitoError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrCognitoError(tt.args.err)
			if err == nil {
				t.Errorf("ErrCognitoError() constructor returned nil error")
			}
		})
	}
}

func TestErrAwsSESError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrAwsSESError(tt.args.err)
			if err == nil {
				t.Errorf("ErrAwsSESError() constructor returned nil error")
			}
		})
	}
}

func TestErrAwsSTSError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrAwsSTSError(tt.args.err)
			if err == nil {
				t.Errorf("ErrAwsSTSError() constructor returned nil error")
			}
		})
	}
}
