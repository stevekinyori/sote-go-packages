package sHelper

import (
	"fmt"
	"os"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sLogger"

	"github.com/integrii/flaggy"
)

type Parameter struct {
	AppName     string
	Description string
	Version     string
	Notes       string
}

func (p Parameter) Init() Environment {
	sLogger.DebugMethod()
	applicationName := ENVDEFAULTAPPNAME
	targetEnvironment := ENVDEFAULTTARGET
	isVerbose := false
	appDescription := `%s.
		Version:
			- %s

		Constraints:
			- At start up, you must pass the application name and the environment for the business service.
			- There is no CLI for this business service.

		Notes:
			%s`

	// Set your program's name and description.  These appear in help outputPtr.
	flaggy.SetName(p.AppName)
	flaggy.SetDescription(fmt.Sprintf(appDescription, p.Description, p.Version, p.Notes))

	// You can disable various things by changing bools on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://gitlab.com/getsote/business-service-layer/back-end"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&applicationName, "a", "appName",
		"Used to pull the stored parameters in System Manager for the message broker backbone provider (Ex: synadia).  "+
			"This requires that you have aws credentials/config setup on the system at ~/.aws. (default: '"+ENVDEFAULTAPPNAME+"')")
	flaggy.String(&targetEnvironment, "t", "targetEnv",
		"Pulls configuration information from aws based on the environment supplied (development|staging|demo|production).  "+
			"This requires that you have aws credentials/config setup on the system at ~/.aws. (default: '"+ENVDEFAULTTARGET+"')")
	flaggy.Bool(&isVerbose, "v", "verbose",
		"Verbose output: log all tests as they are run. Also print all text from Log and Logf calls even if the test succeeds.")

	// Set the version and parse all inputs into variables.
	flaggy.SetVersion(p.Version)
	flaggy.Parse()

	if isVerbose {
		sLogger.SetLogLevelDebug()
	}
	sLogger.SetLogMessagePrefix(applicationName)

	appEnvironment, soteErr := sConfigParams.GetEnvironmentAppEnvironment()
	if soteErr.ErrCode != nil && appEnvironment == "" {
		appEnvironment = targetEnvironment
		os.Setenv("APP_ENVIRONMENT", appEnvironment)
	}

	env, soteErr := NewEnvironment(applicationName, targetEnvironment, appEnvironment)

	if soteErr.ErrCode != nil {
		flaggy.ShowHelpAndExit(soteErr.FmtErrMsg)
	}
	return env
}
