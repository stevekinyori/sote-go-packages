package sDocument

/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS}

    {Other categories of restrictions}
    * AWS S3 Bucket must be mounted.
    * Mount point environment variable must be defined (Check constant.go for the name).

NOTES:
    {Enter any additional notes that you believe will help the next developer.}
*/

import (
	// Add imports here
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	INBOUNDOBJECTKEYNAME FieldName = iota + 1
)

const (
	PARENTDOCUMENTKEY DocumentKey = iota + 1
	SUPPORTINGDOCUMENTSKEY
)

type FieldName int
type DocumentKey int

var (
	InboundObjectKeyFieldName = INBOUNDOBJECTKEYNAME
	ParentDocumentKey         = PARENTDOCUMENTKEY
	SupportingDocumentsKey    = SUPPORTINGDOCUMENTSKEY
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

func (fN FieldName) String() (fieldName string) {
	sLogger.DebugMethod()

	fieldName = [...]string{"inbound-object-key"}[fN-1]

	return
}

func (d DocumentKey) String() (documentKey string) {
	sLogger.DebugMethod()

	documentKey = [...]string{"parent-document", "supporting-documents"}[d-1]

	return
}
