/*
General information about the consumers. (Gathered from https://github.com/nats-io/jetstream)
CONSUMERS:
	Consumers come in two flavors, push and pull.  Push is expecting that the service is available at all times.  Pull holds the
	messages until the process is active. There are many setting for these two types,which effect the behavior of the
	consumer.  Only a limited set of push and pull consumers will are supported below. Consumers can either be push based where
	JetStream will deliver the messages as fast as possible to a subject of your choice or pull based for typical work queue like
	behavior.

	The Consumer subject filter must be a subset of the Stream subject. Error Code: 336100
	Here are some examples:
		Stream subject: IMAGES  Consumer subject filter: IMAGES -> works
		Stream subject: IMAGES  Consumer subject filter: IMAGES.cat -> DOES NOT works
		--
		Stream subject: IMAGES.*  Consumer subject filter: IMAGES -> works
		Stream subject: IMAGES.*  Consumer subject filter: IMAGES.cat -> works
		--
		Stream subject: IMAGES.CATS  Consumer subject filter: IMAGES -> DOES NOT works
		Stream subject: IMAGES.cat  Consumer subject filter: IMAGES.* -> works
*/
package sMessage

import (
	"log"
	"strings"

	"github.com/nats-io/jsm.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	ACKPOLICYNONE     = "none"
	ACKPOLICYALL      = "all"
	ACKPOLICYEXPLICIT = "explicit"
	//
	DELIVERPOLICYALL  = "all"
	DELIVERPOLICYLAST = "last"
	DELIVERPOLICYNEXT = "next"
	//
	DELIVERYSUBJECTPULL = ""
	//
	REPLAYPOLICYINSTANT  = "instant"
	REPLAYPOLICYORIGINAL = "original"
)

/*
CreatePullConsumerWithReplayInstant will create a consumer. If the consumer exists, it will load.
	This will all read all the message in the stream without concern of the order.  This is go when transaction
	order doesn't matter.
	Required parameters:
		streamName
		consumerName
		durableName (Required for a pull consumer)
		subjectFilter
		maxDeliver
			Sote defaults value: 1 (Sote Max value is 3)
		nc (pointer to a Jetstream connection)

	Set values:
		AckPolicy: explicit (required for a pull consumer)
		DeliverPolicy: all
		DeliverySubject: "" (required for a pull consumer)
		ReplayPolicy: instant
*/
func CreatePullConsumerWithReplayInstant(streamName, consumerName, durableName, subjectFilter string, maxDeliveries int, jsmManager *JSMManager) (pConsumer *jsm.Consumer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err        error
		errDetails = make(map[string]string)
	)

	if soteErr = validateConsumerParams(streamName, consumerName, durableName, DELIVERYSUBJECTPULL, subjectFilter, jsmManager); soteErr.ErrCode == nil {
		pConsumer, err = jsmManager.Manager.LoadOrNewConsumer(streamName, consumerName, jsm.DurableName(durableName), jsm.FilterStreamBySubject(subjectFilter), jsm.ReplayInstantly(),
			jsm.MaxDeliveryAttempts(setMaxDeliver(maxDeliveries)), jsm.AcknowledgeExplicit())
		if err != nil {
			errDetails["NATS ERROR:"] = err.Error()
			if strings.Contains("consumer filter subject is not a valid subset of the interest subjects", err.Error()) {
				soteErr = sError.GetSError(336100, sError.BuildParams([]string{streamName, subjectFilter}), errDetails)
			} else {
				soteErr = sError.GetSError(805000, nil, errDetails)
				log.Fatal(soteErr.FmtErrMsg)
			}
		}
	}

	return
}

/*
CreatePullConsumerWithReplayOriginal will create a consumer. If the consumer exists, it will load.
	This will read all the messages from the stream as they were received.  It is good for playback of actions.
	Required parameters:
		streamName
		consumerName
		durableName (Required for a pull consumer)
		subjectFilter
		maxDeliver
			Sote defaults value: 1 (Sote Max value is 3)
		nc (pointer to a Jetstream connection)

	Set values:
		AckPolicy: explicit (required for a pull consumer)
		DeliverPolicy: all
		DeliverySubject: "" (required for a pull consumer)
		ReplayPolicy: original
*/
func CreatePullConsumerWithReplayOriginal(streamName, consumerName, durableName, subjectFilter string, maxDeliveries int, jsmManager *JSMManager) (pConsumer *jsm.Consumer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err        error
		errDetails = make(map[string]string)
	)

	if soteErr = validateConsumerParams(streamName, consumerName, durableName, DELIVERYSUBJECTPULL, subjectFilter, jsmManager); soteErr.ErrCode == nil {
		pConsumer, err = jsmManager.Manager.LoadOrNewConsumer(streamName, consumerName, jsm.DurableName(durableName), jsm.FilterStreamBySubject(subjectFilter), jsm.ReplayAsReceived(),
			jsm.MaxDeliveryAttempts(setMaxDeliver(maxDeliveries)), jsm.AcknowledgeExplicit())
		if err != nil {
			errDetails["NATS ERROR:"] = err.Error()
			if strings.Contains("consumer filter subject is not a valid subset of the interest subjects", err.Error()) {
				soteErr = sError.GetSError(336100, sError.BuildParams([]string{streamName, subjectFilter}), errDetails)
			} else {
				soteErr = sError.GetSError(805000, nil, errDetails)
				log.Fatal(soteErr.FmtErrMsg)
			}
		}
	}

	return
}

func DeleteConsumer(pConsumer *jsm.Consumer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateConsumer(pConsumer); soteErr.ErrCode == nil {
		err := pConsumer.Delete()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

func validateConsumerParams(streamName, consumerName, durableName, deliverySubject, subjectFilter string, jsmmManager *JSMManager) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = validateStreamName(streamName)

	if soteErr.ErrCode == nil {
		soteErr = validateConsumerName(consumerName)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateDurableName(durableName, deliverySubject)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateDeliverySubject(deliverySubject, deliverySubject)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateSubjectFilter(subjectFilter)
	}

	if soteErr.ErrCode == nil && jsmmManager == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"JSMManager"}), nil)
	}

	return
}

func validateConsumerName(consumerName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(consumerName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"consumerName"}), nil)
	}

	return
}

func validateDurableName(durableName, deliverySubject string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	// Durable name is required when the consumer is a pull.  Otherwise, it is optional.
	if len(durableName) == 0 && deliverySubject == DELIVERYSUBJECTPULL {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"durableName"}), nil)
	}

	return
}

func validateDeliverySubject(subjectFilter, deliverySubject string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	// Subject Filter must be empty when the consumer is a pull.  Otherwise, it is required
	if len(subjectFilter) > 0 && deliverySubject == DELIVERYSUBJECTPULL {
		soteErr = sError.GetSError(200515, sError.BuildParams([]string{"subjectFilter", "deliverySubject"}), nil)
	} else if len(subjectFilter) == 0 && deliverySubject != DELIVERYSUBJECTPULL {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
	}

	return
}

func validateSubjectFilter(deliverySubject string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(deliverySubject) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"deliverySubject"}), nil)
	}

	return
}

func validateConsumer(pConsumer *jsm.Consumer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if pConsumer == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Consumer"}), nil)
	}

	return
}
/*
	setMaxDeliver forces the value be between 1 and 10.  If it below 1, then 1 and if greater than 10, its set to 3.
 */
func setMaxDeliver(tMaxDeliver int) (maxDeliver int) {
	sLogger.DebugMethod()

	if tMaxDeliver <= 0 {
		maxDeliver = 1
	} else if tMaxDeliver > 10 {
		maxDeliver = 3
	} else {
		maxDeliver = tMaxDeliver
	}

	return
}
