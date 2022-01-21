package sError

import (
	"testing"
)

func TestErrDBConnectionFailed(t *testing.T) {
	type args struct {
		message   string
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
			err := ErrDBConnectionFailed(tt.args.message, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrDBConnectionFailed() constructor returned nil error")
			}
		})
	}
}

func TestErrEnvFileMissing(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrEnvFileMissing(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrEnvFileMissing() constructor returned nil error")
			}
		})
	}
}

func TestErrFileNotFound(t *testing.T) {
	type args struct {
		fileName  interface{}
		message   interface{}
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
			err := ErrFileNotFound(tt.args.fileName, tt.args.message, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrFileNotFound() constructor returned nil error")
			}
		})
	}
}

func TestErrEnvironmentMissing(t *testing.T) {
	type args struct {
		envName   interface{}
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
			err := ErrEnvironmentMissing(tt.args.envName, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrEnvironmentMissing() constructor returned nil error")
			}
		})
	}
}

func TestErrEnvironmentInvalid(t *testing.T) {
	type args struct {
		envName   interface{}
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
			err := ErrEnvironmentInvalid(tt.args.envName, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrEnvironmentInvalid() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidDBConnection(t *testing.T) {
	type args struct {
		DBName       interface{}
		DBDriverName interface{}
		Port         interface{}
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
			err := ErrInvalidDBConnection(tt.args.DBName, tt.args.DBDriverName, tt.args.Port, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidDBConnection() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidDBAuthentication(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidDBAuthentication(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidDBAuthentication() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidDBSSLMode(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidDBSSLMode(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidDBSSLMode() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidConnectionType(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrInvalidConnectionType(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidConnectionType() constructor returned nil error")
			}
		})
	}
}

func TestErrNATSNKeyMissing(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNATSNKeyMissing(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNATSNKeyMissing() constructor returned nil error")
			}
		})
	}
}

func TestErrNoNATSConnection(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNoNATSConnection(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrNoNATSConnection() constructor returned nil error")
			}
		})
	}
}

func TestErrUnexpectedSign(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrUnexpectedSign(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrUnexpectedSign() constructor returned nil error")
			}
		})
	}
}

func TestErrKIDNotFound(t *testing.T) {
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
			name: "good case: error successfully constructed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrKIDNotFound(tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrKIDNotFound() constructor returned nil error")
			}
		})
	}
}

func TestErrKIDMissingFromToken(t *testing.T) {
	type args struct {
		kID       string
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
			err := ErrKIDMissingFromToken(tt.args.kID, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrKIDMissingFromToken() constructor returned nil error")
			}
		})
	}
}

func TestErrKIDDoesNotMatchPublicKeySet(t *testing.T) {
	type args struct {
		kID       string
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
			err := ErrKIDDoesNotMatchPublicKeySet(tt.args.kID, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrKIDDoesNotMatchPublicKeySet() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidRegion(t *testing.T) {
	type args struct {
		region      string
		environment string
		operation   string
		Err         error
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
			err := ErrInvalidRegion(tt.args.region, tt.args.environment, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidRegion() constructor returned nil error")
			}
		})
	}
}

func TestErrInvalidURL(t *testing.T) {
	type args struct {
		param     string
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
			err := ErrInvalidURL(tt.args.param, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrInvalidURL() constructor returned nil error")
			}
		})
	}
}

func TestErrOutofValidRange(t *testing.T) {
	type args struct {
		param     string
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
			err := ErrOutofValidRange(tt.args.param, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrOutofValidRange() constructor returned nil error")
			}
		})
	}
}
