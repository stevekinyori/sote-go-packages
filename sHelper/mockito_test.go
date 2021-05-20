package sHelper

import (
	"testing"
)

type fakeT struct {
	iTesting
}

func (fakeT) Helper() {
}

func (fakeT) Fatal(args ...interface{}) {
}

func TestAssertEqual(t *testing.T) {
	AssertEqual(t, "HELLO", "HELLO")
	AssertEqual(&fakeT{}, "HELLO", "WORLD")
}

func TestMockRunHelper(t *testing.T) {
	env := MockRunHelper(t, "bsl-organization-wildcard", "bsl.organization.list")
	AssertEqual(t, env.ApplicationName, ENVDEFAULTAPPNAME)
	AssertEqual(t, env.TargetEnvironment, ENVDEFAULTTARGET)
	AssertEqual(t, env.AppEnvironment, ENVDEFAULTTARGET)

	helper := NewHelper(env)
	helper.AddSubscriber("bsl-organization-wildcard", "bsl.organization.list", nil, nil)
}

func TestMockRunHelperInvalidSubject(t *testing.T) {
	env := MockRunHelper(&fakeT{}, "bsl-organization-wildcard", "bsl.organization.add", "bsl.organization.change", "bsl.organization.delete")
	AssertEqual(t, env.ApplicationName, ENVDEFAULTAPPNAME)
	AssertEqual(t, env.TargetEnvironment, ENVDEFAULTTARGET)
	AssertEqual(t, env.AppEnvironment, ENVDEFAULTTARGET)

	helper := NewHelper(env)
	helper.AddSubscriber("bsl-organization-wildcard", "bsl.organization.list", nil, nil)
}
