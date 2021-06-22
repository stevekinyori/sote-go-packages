package sHelper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/integrii/flaggy"

	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func newParam() Parameter {
	return Parameter{
		Version:     "v2021.1.0",
		AppName:     "My App",
		Description: "My Description",
	}
}

func TestParameterInit(t *testing.T) {
	params := newParam()
	AssertEqual(t, params.Version, "v2021.1.0")
	AssertEqual(t, params.AppName, "My App")
	AssertEqual(t, params.Description, "My Description")
}

func TestParameterDefaultEnv(t *testing.T) {
	os.Setenv("APP_ENVIRONMENT", "")
	os.Setenv("XDG_CONFIG_HOME", "/")
	params := newParam()
	env := params.Init()
	AssertEqual(t, flaggy.DefaultParser.Version, params.Version)
	AssertEqual(t, flaggy.DefaultParser.Name, params.AppName)

	AssertEqual(t, env.ApplicationName, ENVDEFAULTAPPNAME)
	AssertEqual(t, env.TargetEnvironment, ENVDEFAULTTARGET)
	AssertEqual(t, env.AppEnvironment, ENVDEFAULTTARGET)
	AssertEqual(t, env.TestMode, true)
}

func TestParameterShort(t *testing.T) {
	flaggy.ResetParser()
	os.Args = []string{"main", "-a", "SoteApp", "-t", "production"}
	os.Setenv("APP_ENVIRONMENT", "development")
	params := newParam()
	env := params.Init()
	AssertEqual(t, env.ApplicationName, "SoteApp")
	AssertEqual(t, env.TargetEnvironment, "production")
	AssertEqual(t, env.AppEnvironment, "development")
	AssertEqual(t, env.TestMode, false)
}

func TestParameterLong(t *testing.T) {
	flaggy.ResetParser()
	os.Args = []string{"main", "--appName", "SoteApp", "--targetEnv", "production"}
	params := newParam()
	env := params.Init()
	AssertEqual(t, env.ApplicationName, "SoteApp")
	AssertEqual(t, env.TargetEnvironment, "production")
	AssertEqual(t, env.TestMode, false)
}

func TestParameterVerbose(t *testing.T) {
	flaggy.ResetParser()
	os.Args = []string{"main", "-v"}
	AssertEqual(t, sLogger.GetLogLevel(), sLogger.InfoLogLevel)
	params := newParam()
	params.Init()
	AssertEqual(t, sLogger.GetLogLevel(), sLogger.DebugLogLevel)
	sLogger.SetLogLevelInfo()
}

func TestParameterVersion(t *testing.T) {
	flaggy.PanicInsteadOfExit = true
	flaggy.ResetParser()
	os.Args = []string{"main", "--version"}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic on show version and exit call")
		}
	}()
	newParam().Init()
}

func TestParameterInvalidEnv(t *testing.T) {
	flaggy.PanicInsteadOfExit = true
	flaggy.ResetParser()
	os.Args = []string{"main", "--targetEnv", "UNKNOWN"}
	defer func() {
		r := recover()
		AssertEqual(t, r, "Panic instead of exit with code: 2")
	}()
	newParam().Init()
}

func TestParameterTargetEnv(t *testing.T) {
	flaggy.ResetParser()
	parent, _ := filepath.Abs(".config")
	tempNatsDir := filepath.Join(parent, "nats")
	err := os.MkdirAll(tempNatsDir, os.ModePerm)
	AssertEqual(t, err, nil)
	ioutil.WriteFile(filepath.Join(tempNatsDir, "context.txt"), []byte("sote-production"), 0644)
	os.Args = []string{"main", "--config", parent}
	params := newParam()
	env := params.Init()
	os.RemoveAll(tempNatsDir)
	AssertEqual(t, env.TargetEnvironment, "production")
	AssertEqual(t, os.Getenv("XDG_CONFIG_HOME"), parent)
	os.Setenv("XDG_CONFIG_HOME", "")
}
