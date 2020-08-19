package sDatabase

import (
	"fmt"
	"os/exec"

	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func ExecutePSQLWithFile(fileName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	out, err := exec.Command("ls").Output()
	if err != nil {
		sLogger.Debug("ls command failed.")
	}
	fmt.Printf("ls output: %v",out)
	return 
}
