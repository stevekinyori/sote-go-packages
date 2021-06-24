/*
This will process a set of tables that belong together

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).

    * A fully qualified table group json file must be provided and exist.

NOTES:
    None
*/
package sDatabase

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	LOGMESSAGEPREFIX = "packages"
)

var (
	tableGroupTableColumnInfo []TableColumnInfo
)

// Input file structure
type TableGroup struct {
	TableGroupName string   `json:"table-group-name"`
	Schema         string   `json:"schema"`
	Tables         []string `json:"tables"`
}

type TableColumnInfo struct {
	TableName  string
	ColumnInfo []SColumnInfo
}

var (
// Add Variables here for the file (Remember, they are global)
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

func GetTableGroupInfo(fileName string, testMode bool) (tableGroupTableInfo TableGroup, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if fileName == "" {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{fileName}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	if soteErr.ErrCode == nil {
		if tableGroupFileHandle, err := ioutil.ReadFile(fileName); err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				soteErr = sError.GetSError(109999, sError.BuildParams([]string{fileName + "/" + err.Error()}), sError.EmptyMap)
				sLogger.Info(soteErr.FmtErrMsg)
			} else {
				soteErr = sError.GetSError(199999, sError.BuildParams([]string{fileName + "/" + err.Error()}), sError.EmptyMap)
				sLogger.Info(soteErr.FmtErrMsg)
			}
		} else {
			if err = json.Unmarshal(tableGroupFileHandle, &tableGroupTableInfo); err != nil {
				soteErr = sError.GetSError(207110, sError.BuildParams([]string{fileName}), sError.EmptyMap)
				sLogger.Info(soteErr.FmtErrMsg)
			}
		}
	}

	if soteErr.ErrCode != nil && testMode == false {
		panic(soteErr.FmtErrMsg)
	}

	return
}
