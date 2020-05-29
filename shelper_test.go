package shelper

import (
	"strings"
	"testing"
)

func TestContainsSpecialCharacters(t *testing.T) {
	myFieldName := "TestField"
	myFieldValue := " Test.Field"
	if x := shelper.ContainsSpecialCharacters(myFieldName, myFieldValue); x.ErrCode = 400060 {
		t.Errorf(x.FmtErrMsg)
	}
}

func TestContainsSpecialCharactersNone(t *testing.T) {
	myFieldName := "TestField"
	myFieldValue := " Test Field"
	if x := shelper.ContainsSpecialCharacters(myFieldName, myFieldValue); x.ErrCode = 0 {
		t.Errorf("%v (%v) does not contain special characters", myFieldName, myFieldValue)
	}
}