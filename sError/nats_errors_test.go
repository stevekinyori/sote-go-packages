package sError

import (
	"testing"
)

func TestErrJetStreamError(t *testing.T) {
	type args struct {
		err error
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
			err := ErrJetStreamError(tt.args.err)
			if err == nil {
				t.Errorf("ErrJetStreamError() constructor returned nil error")
			}
		})
	}
}

func TestErrNATSSubscriptionError(t *testing.T) {
	type args struct {
		subscriptionName interface{}
		subject          interface{}
		err              error
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
			err := ErrNATSSubscriptionError(tt.args.subscriptionName, tt.args.subject, tt.args.err)
			if err == nil {
				t.Errorf("ErrNATSSubscriptionError() constructor returned nil error")
			}
		})
	}
}

func TestErrNATSMissingStreamPointer(t *testing.T) {
	type args struct {
		err error
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
			err := ErrNATSMissingStreamPointer(tt.args.err)
			if err == nil {
				t.Errorf("ErrNATSMissingStreamPointer() constructor returned nil error")
			}
		})
	}
}

func TestErrNATSStreamCreationError(t *testing.T) {
	type args struct {
		streamName interface{}
		err        error
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
			err := ErrNATSStreamCreationError(tt.args.streamName, tt.args.err)
			if err == nil {
				t.Errorf("ErrNATSStreamCreationError() constructor returned nil error")
			}
		})
	}
}

func TestErrNATSConsumerCreationError(t *testing.T) {
	type args struct {
		streamName   interface{}
		consumerName interface{}
		err          error
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
			err := ErrNATSConsumerCreationError(tt.args.streamName, tt.args.consumerName, tt.args.err)
			if err == nil {
				t.Errorf("ErrNATSConsumerCreationError() constructor returned nil error")
			}
		})
	}
}

func TestErrNATSInvalidConsumerSubjectFilter(t *testing.T) {
	type args struct {
		streamName            interface{}
		consumerSubjectFilter interface{}
		err                   error
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
			err := ErrNATSInvalidConsumerSubjectFilter(tt.args.streamName, tt.args.consumerSubjectFilter, tt.args.err)
			if err == nil {
				t.Errorf("ErrNATSInvalidConsumerSubjectFilter() constructor returned nil error")
			}
		})
	}
}
