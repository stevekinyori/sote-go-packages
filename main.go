package main

import (
	"context"
	"log"

	"github.com/integrii/flaggy"
	"gitlab.com/soteapps/packages/v2022/sLogger"
	"gitlab.com/soteapps/packages/v2022/sMigration"
)

const (
	// VERSION Most Recent Version of the packages - THIS MUST MATCH THE RELEASE TAG!!!!!
	VERSION = "v2022.21.0"
	// FOURSPACES Formatting for the application display when run with help
	FOURSPACES       string = "    "
	LOGMESSAGEPREFIX        = "packages"
)

var (
	targetEnvironment string
	service           string
	action            string
	setupDir          string
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)

	appDescription := "This contains customized packages specific to Sote" +
		"\nVersion: \n" +
		FOURSPACES + "- " + VERSION + "\n" +
		"\nConstraints: \n" +
		FOURSPACES + "- When you want to run commands related to these packages, " +
		"you must provide the required parameters specific to each package \n" +
		FOURSPACES + "- There is no CLI for these packages.\n" +
		"\nNotes:\n" +
		FOURSPACES + "None\n"

	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("sote-packages")
	flaggy.SetDescription(appDescription)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://gitlab.com/soteapps/packages"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&targetEnvironment, "e", "targetEnv",
		"Pulls configuration information from aws based on the environment supplied (development|staging|demo|production).  "+
			"This requires that you have aws credentials/config setup on the system at ~/.aws. (default: '')")
	flaggy.String(&service, "s", "service",
		"This defines the service you are calling .Can be one of these [migration|seed]")
	flaggy.String(&action, "a", "action",
		"This defines the action to perform on the set service .Can be one of these [setup|run]")
	flaggy.String(&setupDir, "d", "dir",
		"This defines the directory where the installation files are to be set. Defaults to current directory path of this file")
	// Set the version
	flaggy.SetVersion(VERSION)
	// parse all inputs into variables.
	flaggy.Parse()
}

func main() {
	sLogger.DebugMethod()

	switch service {
	case sMigration.MigrationType, sMigration.SeedingAction:
		if targetEnvironment == "" || action == "" {
			flaggy.ShowHelpAndExit("")
		}

		if soteErr := sMigration.Run(context.Background(), targetEnvironment, service, action, setupDir); soteErr.ErrCode != nil {
			log.Fatal(soteErr.FmtErrMsg)
		}
	default:
		flaggy.ShowHelpAndExit("")
	}
}
