package sHelper

import (
	"fmt"

	"bou.ke/monkey"
	"gitlab.com/soteapps/packages/v2021/sError"
)

var (
	Patch       = monkey.Patch
	PatchMethod = monkey.PatchInstanceMethod
)

type PatchGuard = monkey.PatchGuard

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

func MockRunHelper(t iTesting, verifyConsumerName string, verifySubject ...string) Environment {
	var (
		helperPatch *PatchGuard
	)
	helperPatch = Patch(NewHelper, func(env Environment) *Helper {
		helperPatch.Unpatch()
		helper := Helper{
			Env: env,
		}
		helper.Run = func(isGoroutine bool) {}
		helper.AddSubscriber = func(consumerName, subject string, _ MessageListener, _ *Schema, _ ...string) sError.SoteError {
			AssertEqual(t, consumerName, verifyConsumerName)
			found := false
			for _, s := range verifySubject {
				if s == subject {
					found = true
					break
				}
			}
			if !found {
				AssertEqual(t, subject, verifySubject)
			}
			return sError.SoteError{}
		}
		return &helper
	})
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	return env
}
