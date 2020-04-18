package packages

import (
	"strings"
	"testing"

	"gitlab.com/getsote/packages/serror"
)

func TestIncorrectParams(t *testing.T) {
	var paramValues = []string{""}
	if x := serror.GetSError(100100, paramValues); x.ErrCode != 230000 {
		t.Errorf("The wrong error code (%v) was returned.  230000 should have been returned.", x.ErrCode)
	}
}

func TestErrorCodeNotFound(t *testing.T) {
	errCode := 999999999
	if x := serror.GetSError(errCode, nil); x.ErrCode != 109999 {
		t.Errorf("%v should have returned an error of 109999", errCode)
	}
}

func Test100100Error(t *testing.T) {
	var errCode = 100100
	var paramValues = []string{"SUPER_USER, EXECUTIVE", "DELETE"}
	validateReply(t, errCode, paramValues, serror.GetSError(errCode, paramValues))
}
func Test200200Error(t *testing.T) {
	var errCode = 200200
	var paramValues = []string{"PARAMETER_NAME", "DATA_TYPE"}
	validateReply(t, errCode, paramValues, serror.GetSError(200200, paramValues))
}
func Test200250Error(t *testing.T) {
	var errCode = 200250
	var paramValues = []string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"}
	validateReply(t, errCode, paramValues, serror.GetSError(200250, paramValues))
}
func Test200260Error(t *testing.T) {
	var errCode = 200260
	var paramValues = []string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(200260, paramValues))
}
func Test200500Error(t *testing.T) {
	var errCode = 200500
	var paramValues = []string{"OBJECT_NAME"}
	validateReply(t, errCode, paramValues, serror.GetSError(200500, paramValues))
}
func Test200510Error(t *testing.T) {
	var errCode = 200510
	var paramValues = []string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(200510, paramValues))
}
func Test200511Error(t *testing.T) {
	var errCode = 200511
	var paramValues = []string{"PARAMETER_NAME", "PARAMETER_NAME"}
	validateReply(t, errCode, paramValues, serror.GetSError(200511, paramValues))
}
func Test201000Error(t *testing.T) {
	var errCode = 201000
	var paramValues = []string{"DETAILED_MESSAGE"}
	validateReply(t, errCode, paramValues, serror.GetSError(201000, paramValues))
}
func Test202000Error(t *testing.T) {
	var errCode = 202000
	var paramValues = []string{"ENVIRONMENT"}
	validateReply(t, errCode, paramValues, serror.GetSError(202000, paramValues))
}
func Test310000Error(t *testing.T) {
	var errCode = 310000
	var paramValues = []string{"KEY_NAME"}
	validateReply(t, errCode, paramValues, serror.GetSError(310000, paramValues))
}
func Test310005Error(t *testing.T) {
	var errCode = 310005
	var paramValues = []string{"KEY_NAME"}
	validateReply(t, errCode, paramValues, serror.GetSError(310005, paramValues))
}
func Test320000Error(t *testing.T) {
	var errCode = 320000
	var paramValues = []string{"PARAMETER_LIST"}
	validateReply(t, errCode, paramValues, serror.GetSError(320000, paramValues))
}
func Test400000Error(t *testing.T) {
	var errCode = 400000
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400000, paramValues))
}
func Test400010Error(t *testing.T) {
	var errCode = 400010
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400010, paramValues))
}
func Test400020Error(t *testing.T) {
	var errCode = 400020
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400020, paramValues))
}
func Test400030Error(t *testing.T) {
	var errCode = 400030
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400030, paramValues))
}
func Test400040Error(t *testing.T) {
	var errCode = 400040
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400040, paramValues))
}
func Test400050Error(t *testing.T) {
	var errCode = 400050
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400050, paramValues))
}
func Test400060Error(t *testing.T) {
	var errCode = 400060
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400060, paramValues))
}
func Test400070Error(t *testing.T) {
	var errCode = 400070
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400070, paramValues))
}
func Test400080Error(t *testing.T) {
	var errCode = 400080
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400080, paramValues))
}
func Test400090Error(t *testing.T) {
	var errCode = 400090
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"}
	validateReply(t, errCode, paramValues, serror.GetSError(400090, paramValues))
}
func Test400100Error(t *testing.T) {
	var errCode = 400100
	var paramValues = []string{"FIELD_VALUE"}
	validateReply(t, errCode, paramValues, serror.GetSError(400100, paramValues))
}
func Test405110Error(t *testing.T) {
	var errCode = 405110
	var paramValues = []string{"OBJECT_TYPE", "SYSTEM_ID"}
	validateReply(t, errCode, paramValues, serror.GetSError(405110, paramValues))
}
func Test405120Error(t *testing.T) {
	var errCode = 405120
	var paramValues = []string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"}
	validateReply(t, errCode, paramValues, serror.GetSError(405120, paramValues))
}
func Test601000Error(t *testing.T) {
	var errCode = 601000
	var paramValues = []string{"ENVIRONMENT"}
	validateReply(t, errCode, paramValues, serror.GetSError(601000, paramValues))
}
func Test602000Error(t *testing.T) {
	var errCode = 602000
	var paramValues = []string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"}
	validateReply(t, errCode, paramValues, serror.GetSError(602000, paramValues))
}
func Test609999Error(t *testing.T) {
	var errCode = 609999
	var paramValues = []string{"START_VARIABLE_MISSING"}
	validateReply(t, errCode, paramValues, serror.GetSError(609999, paramValues))
}
func Test700000Error(t *testing.T) {
	var errCode = 700000
	var paramValues = []string{"LIST_PARAMETERS"}
	validateReply(t, errCode, paramValues, serror.GetSError(700000, paramValues))
}
func TestGenMarkDown(t *testing.T) {
	if x := serror.GenMarkDown(); !strings.Contains(x, serror.MakeDownTitleBar) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	}
}
func validateReply(t *testing.T, errCode int, paramValues []string, x serror.SoteError) {
	if x.ErrCode != 109999 && x.ErrCode != 230000 {
		for i, _ := range paramValues {
			if !strings.Contains(x.FmtErrMsg, paramValues[i]) {
				t.Errorf("Error Code Tested: %v - The %v was not found in the FmtErrMsg property returned", errCode, paramValues[i])
			}
		}
	}
}
