package sError

import "fmt"

const (
	// Error Category
	NATS_ERROR = "nats-error"

	JET_STREAM_ERROR                     = 206000
	NATS_SUBSCRIPTION_ERROR              = 206050
	NATS_MISSING_STREAM_POINTER          = 206300
	NATS_STREAM_CREATE_ERROR             = 206400
	NATS_CONSUMER_CREATE_ERROR           = 206600
	NATS_INVALID_CONSUMER_SUBJECT_FILTER = 206700
)

func ErrJetStreamError(err error) error {
	return &SoteError{
		ErrCode:   JET_STREAM_ERROR,
		ErrType:   NATS_ERROR,
		FmtErrMsg: ": Jetstream is not enabled",
		Err:       err,
	}
}

func ErrNATSSubscriptionError(subscriptionName, subject interface{}, err error) error {
	return &SoteError{
		ErrCode: NATS_SUBSCRIPTION_ERROR,
		ErrType: NATS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": (%v) is an invalid subscription. Subject: %v",
			subscriptionName, subject),
		Err: err,
	}
}

func ErrNATSMissingStreamPointer(err error) error {
	return &SoteError{
		ErrCode:   NATS_SUBSCRIPTION_ERROR,
		ErrType:   NATS_ERROR,
		FmtErrMsg: ": Stream pointer is nil. Must be a validate pointer to a stream.",
		Err:       err,
	}
}

func ErrNATSStreamCreationError(streamName interface{}, err error) error {
	return &SoteError{
		ErrCode: NATS_STREAM_CREATE_ERROR,
		ErrType: NATS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": Stream creation encountered an error that is not expected. Stream Name: %v",
			streamName),
		Err: err,
	}
}

func ErrNATSConsumerCreationError(streamName, consumerName interface{}, err error) error {
	return &SoteError{
		ErrCode: NATS_CONSUMER_CREATE_ERROR,
		ErrType: NATS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": Consumer creation encountered an error that is not expected. Stream Name: %v Consumer Name: %v",
			streamName, consumerName),
		Err: err,
	}
}

func ErrNATSInvalidConsumerSubjectFilter(streamName, consumerSubjectFilter interface{}, err error) error {
	return &SoteError{
		ErrCode: NATS_CONSUMER_CREATE_ERROR,
		ErrType: NATS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": The consumer subject filter must be a subset of the stream subject. Stream Name: %v Consumer Subject Filter: %v",
			streamName, consumerSubjectFilter),
		Err: err,
	}
}
