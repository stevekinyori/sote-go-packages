/*
	This is a wrapper for validation messages used by Sote GO software developers.

	The GO helper package should not be used directly.  This package sets the format of 
	the validation message so they are uniform.
*/

package shelper

import (
	"regexp"
	"fmt"
	"runtime"
	"gitlab.com/soteapps/packages/slogger"
	"gitlab.com/soteapps/packages/serror"
)

type soteRet struct {
	ErrCode          interface{}
	ErrType          string
	ParamCount       int
	ParamDescription string
	FmtErrMsg        string
	LUErrorDetails   map[string]string 
	Loc              string
}

func ErrLoc() string {
    pc := make([]uintptr, 10) 
    runtime.Callers(2, pc)
    f := runtime.FuncForPC(pc[0])
    return f.Name()
}
/**
 * contains Special Characters will test if the following characters are in the string
 * !"#\$%&''()*+,-./:;<=>?@[\]^_`{|}~
 *
 * @param string $myFieldName
 * @param string $myFieldValue
 *
 * @return array Keys: ErrCode, ErrType, ParamCount, ParamDescription, FmtErrMsg, LUErrorDetails, Loc
 *               FmtErrMsg on an error will return special_characters (key) with the list of special characters (value)
 */


func ContainsSpecialCharacters(myFieldName string, myFieldValue string) soteRet {

	slogger.DebugMethod()

	var result soteRet

	var re = regexp.MustCompile(`(?m)[!"#\$%&\'\'\(\)\*\+,-\.\/:;<=>\?@\[\\\\\]\^_\x60\{\|}~]`)

	if(re.MatchString(myFieldValue)) {
		result = soteRet{ 0, "Success", 0, "", "", make(map[string]string), ErrLoc()}
		
		
	} else { 
		result = soteRet(serror.GetSError(400060, []string{myFieldName , myFieldValue }, serror.EmptyMap))
	}
	
	slogger.Debug("ErrCode: " + fmt.Sprint(result.ErrCode)+ " ErrType: " + result.ErrType + " FmtErrMsg: " + result.FmtErrMsg + " Loc: " + result.Loc)
	return result
}