package sHelper

import (
	"os"
	"testing"
)

func TestEnvCreate(t *testing.T) {
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, os.Getenv("APP_ENVIRONMENT"))
	AssertEqual(t, env.ApplicationName, ENVDEFAULTAPPNAME)
	AssertEqual(t, env.TargetEnvironment, ENVDEFAULTTARGET)
	AssertEqual(t, env.TestMode, env.TargetEnvironment != "production")
}

func TestEnvCreateProduction(t *testing.T) {
	env, _ := NewEnvironment("ProductionApp", "production", "staging")
	AssertEqual(t, env.ApplicationName, "ProductionApp")
	AssertEqual(t, env.TargetEnvironment, "production")
	AssertEqual(t, env.TestMode, false)
}

func TestEnvCreateInvalidEnv(t *testing.T) {
	_, soteErr := NewEnvironment(ENVDEFAULTAPPNAME, "UNKNOWN", "staging")
	if soteErr.ErrCode == nil {
		t.Fatal("Shouldn't be nil")
	}
}
