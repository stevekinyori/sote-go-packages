/*

*/
package sMessage

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*
	PPublish will send a persist message to the stream
*/
func (mm *MessageManager) PPublish(subject, message string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mm.NatsConnection.JetStream()
	if err != nil {
		mm.natsErrorHandle(err, subject, "", "", []byte(message))
	}
	js.Publish(subject, []byte(message))

	return
}
