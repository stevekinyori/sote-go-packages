package sError

import "testing"

func TestErrUnexpected(t *testing.T) {
	type args struct {
		logData   interface{}
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
			err := ErrUnexpected(tt.args.logData, tt.args.operation, tt.args.Err)
			if err == nil {
				t.Errorf("ErrUnexpected() constructor returned nil error")
			}
		})
	}
}
