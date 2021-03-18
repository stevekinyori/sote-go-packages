package sDocument

import (
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"os"
	"runtime"
)

/*  This will return prefix location of tesseracts training data */
func GetTessdataPrefix()  (tessdataPrefix string){
	sLogger.DebugMethod()

	if tessdataPrefix = os.Getenv("TESSDATA_PREFIX"); tessdataPrefix == ""{
		sgoOs := runtime.GOOS
		//TODO automate this for different operating systems
		switch sgoOs {
		case "darwin":
			tessdataPrefix = "/opt/homebrew/Cellar/tesseract/4.1.1/share"
		default:
			sLogger.Info("TESSDATA_PREFIX not set:Couldn't determine OS type")
		}

	}

	return
}
