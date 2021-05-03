package sHelper

import (
	"fmt"
)

type iTesting interface {
	Helper()
	Fatal(args ...interface{})
}

func AssertEqual(t iTesting, actual, expected interface{}) {
	t.Helper() // get caller function a line number
	if expected != actual {
		t.Fatal(fmt.Sprintf("Not equal:\nexpected: %v\nactual:   %v", expected, actual))
	}
}
