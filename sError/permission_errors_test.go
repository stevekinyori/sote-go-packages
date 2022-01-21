package sError

import (
	"testing"
)

func TestErrInvalidISS(t *testing.T) {
	type args struct {
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidISS(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidISS() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidSubject(t *testing.T) {
	type args struct {
		subject   interface{}
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidSubject(tt.args.subject, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidSubject() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidToken(t *testing.T) {
	type args struct {
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidToken(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidToken() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidAppClientID(t *testing.T) {
	type args struct {
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidAppClientID(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidAppClientID() constructor returned nil error")
			}
		})
	}
}

func TestErrTokenExpired(t *testing.T) {
	type args struct {
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrTokenExpired(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrTokenExpired() constructor returned nil error")
			}
		})
	}
}

func TestErrSegmentsCountInvalid(t *testing.T) {
	type args struct {
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrSegmentsCountInvalid(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrSegmentsCountInvalid() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidClaim(t *testing.T) {
	type args struct {
		claimNames interface{}
		operation  string
		Err        error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidClaim(tt.args.claimNames, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidClaim() constructor returned nil error")
			}
		})
	}
}

func TestErrMissingClaim(t *testing.T) {
	type args struct {
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrMissingClaim(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrMissingClaim() constructor returned nil error")
			}
		})
	}
}
