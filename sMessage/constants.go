package sMessage

import (
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

const (
	ORGGANIZATIONSBSID        = "5"
	USERMANAGEMENTBSID        = "6"
	SHIPMENTSBSID             = "7"
	TRIPSBSID                 = "8"
	SHIPMENTSFINTRANSBSID     = "9"
	TRIPFINTRANSBSID          = "10"
	DOCUMENTSBSID             = "11"
	NOTESBSID                 = "12"
	LISTOFVALUESBSID          = "13"
	COGNITOBSID               = "14"
	COGNITOUSERMANAGEMENTBSID = "15"
	QUICKBOOKSBSID            = "16"
)

func GetBusinessServiceIds() []string {
	sLogger.DebugMethod()

	return []string{
		ORGGANIZATIONSBSID,
		USERMANAGEMENTBSID,
		SHIPMENTSBSID,
		TRIPSBSID,
		SHIPMENTSFINTRANSBSID,
		TRIPFINTRANSBSID,
		DOCUMENTSBSID,
		NOTESBSID,
		LISTOFVALUESBSID,
		COGNITOBSID,
		COGNITOUSERMANAGEMENTBSID,
		QUICKBOOKSBSID,
	}
}
