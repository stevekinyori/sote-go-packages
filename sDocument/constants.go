package sDocument

import (
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

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

const (
	PROCESSED                             = "processed"
	INBOUND                               = "inbound"
	PROCESSEDFOLDER                       = PROCESSED
	INBOUNDFOLDER                         = INBOUND
	PROCESSEDLINKKEY                      = PROCESSED
	INBOUNDLINKKEY                        = INBOUND
	BUCKETDOCSUBFOLDER                    = "documents-management-system"
	DOCUMENTSMOUNTPOINTENVIRONMENTVARNAME = "DOCUMENTS_MOUNT_POINT"
	TEMPLATESMOUNTPOINTENVIRONMENTVARNAME = "TEMPLATES_MOUNT_POINT"
)

// List type's here

var (
// Add Variables here for the file (Remember, they are global)
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}
