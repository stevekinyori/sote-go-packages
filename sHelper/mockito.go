package sHelper

import (
	"fmt"
	"testing"

	"bou.ke/monkey"
	"gitlab.com/soteapps/packages/v2021/sError"
)

var Patch = monkey.Patch

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

func MockRunHelper(t *testing.T, verifyConsumerName, verifySubject string) Environment {
	Patch(NewHelper, func(env Environment) *Helper {
		helper := Helper{
			Env: env,
		}
		helper.Run = func(isGoroutine bool) {}
		helper.AddSubscriber = func(consumerName, subject string, _ MessageListener, _ *Schema) sError.SoteError {
			AssertEqual(t, consumerName, verifyConsumerName)
			AssertEqual(t, subject, verifySubject)
			return sError.SoteError{}
		}
		return &helper
	})
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	return env
}
