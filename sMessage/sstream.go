package sMessage

import (
	"log"
	"strconv"
	"strings"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// F      = "f"
	// FILE   = "file"
	M      = "m"
	MEMORY = "memory"
)

/*
	CreateLimitsStream will create a stream.  If it exists, it will load the stream
		Name: Required
			value is set using: no spaces and case sensitive
			Sote defaults value: Required
			Sote immutable: no
		Replicas: Default value: 1
			value is set using: >0
			Sote defaults value: 1
			Sote immutable: no
		Storage: Required
			value is set using: (f)ile, (m)emory
			example: f, file, m, memory
			Sote defaults value: f
			Sote immutable: no
		Subjects: Required
			format of string: <subject>,<subject>,...
			example: images,animals,cars,BILLS
 			Sote defaults value: None
			Sote immutable: no
		--
		Retention: Default value: Limits
			value is set using: Limits, Interest, work queue, workq, work
			Sote defaults value: Limits
			Sote immutable: yes
*/
func CreateOrLoadLimitsStream(streamName, subjects, storage string, replicas int, nc *nats.Conn) (stream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamParams(streamName, subjects, nc); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			stream, err = jsm.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.MemoryStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)),
				jsm.LimitsRetention())
		} else {
			stream, err = jsm.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.FileStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)), jsm.LimitsRetention())
		}
		if err != nil {
			soteErr = sError.GetSError(805000, sError.BuildParams([]string{streamName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	CreateOrLoadWorkStream will create a stream.  If it exists, it will load the stream.
		Name: Required
			value is set using: no spaces and case sensitive
			Sote defaults value: Required
			Sote immutable: no
		Replicas: Default value: 1
			value is set using: >0
			Sote defaults value: 1
			Sote immutable: no
		Storage: Required
			value is set using: (f)ile, (m)emory
			example: f, file, m, memory
			Sote defaults value: f
			Sote immutable: no
		Subjects: Required
			format of string: <subject>,<subject>,...
			example: images,animals,cars,BILLS
 			Sote defaults value: None
			Sote immutable: no
		--
		Retention: Default value: Limits
			value is set using: Limits, Interest, work queue, workq, work
			Sote defaults value: Limits
			Sote immutable: yes
*/
func CreateOrLoadWorkStream(streamName, subjects, storage string, replicas int, nc *nats.Conn) (stream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamParams(streamName, subjects, nc); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			stream, err = jsm.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.MemoryStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)),
				jsm.LimitsRetention())
		} else {
			stream, err = jsm.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.FileStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)), jsm.LimitsRetention())
		}
		if err != nil {
			soteErr = sError.GetSError(805000, sError.BuildParams([]string{streamName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	DeleteStream will destroy the stream
*/
func DeleteStream(pStream *jsm.Stream) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
		err := pStream.Delete()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	StreamInfo loads an existing stream information
*/
func StreamInfo(pStream *jsm.Stream) (info *api.StreamInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
		info, err = pStream.LatestInformation()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	PurgeStream will remove all messages from the stream
*/
func PurgeStream(pStream *jsm.Stream) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
		err := pStream.Purge()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	DeleteMessageFromStream will remove a message from the stream
*/
func DeleteMessageFromStream(pStream *jsm.Stream, sequenceNumber int) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if sequenceNumber <= 0 {
		soteErr = sError.GetSError(400005, sError.BuildParams([]string{"sequenceNumber", "0"}), nil)
	} else {
		// TODO Test what happens if there is no message with the sequence number
		if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
			err := pStream.DeleteMessage(sequenceNumber)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"sequenceNumber(" + strconv.Itoa(sequenceNumber) + ")"}), nil)
				} else {
					soteErr = sError.GetSError(805000, nil, nil)
					log.Fatal(soteErr.FmtErrMsg)
				}
			}
		}
	}

	return
}

func validateStreamParams(streamName, subjects string, nc *nats.Conn) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil && len(subjects) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjects"}), nil)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateConnection(nc)
	}

	return
}

func validateConnection(nc *nats.Conn) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if nc == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Connect"}), nil)
	}

	return

}

func validateStreamName(streamName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamName"}), nil)
	}

	return
}

func validateStream(pStream *jsm.Stream) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if pStream == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Stream"}), nil)
	}

	return
}

func setReplicas(tReplicas int) (replicas int) {
	sLogger.DebugMethod()

	if tReplicas <= 0 {
		replicas = 1
	} else {
		replicas = tReplicas
	}

	return
}
