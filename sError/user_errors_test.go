package sError

import (
	"testing"
)

// TODO: Positive and Negative Test cases - Add to Wiki

func TestErrItemAlreadyExists(t *testing.T) {
	type args struct {
		item      interface{}
		operation string
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
			err := ErrItemAlreadyExists(tt.args.item, tt.args.operation)
			if err == nil {
				t.Errorf("ErrItemAlreadyExists() constructor returned nil error")
			}
		})
	}
}

func TestErrNotAuthorized(t *testing.T) {
	type args struct {
		userRoles   interface{}
		permissions interface{}
		operation   string
		Err         error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive case: error successfully constructed",
			args: args{
				userRoles:   []string{"test-role"},
				permissions: []string{"test-permission"},
				operation:   "test-operation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNotAuthorized(tt.args.userRoles, tt.args.permissions, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNotAuthorized() constructor returned nil error")
			}
		})
	}
}

func TestErrItemNotFound(t *testing.T) {
	type args struct {
		item      interface{}
		operation string
		Err       error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
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
			err := ErrItemNotFound(tt.args.item, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrItemNotFound() constructor returned nil error")
			}
		})
	}
}

// func BenchmarkRandInt(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		ErrDBConnectionFailed("db connection failed", "operation", nil)
// 	}
// }
