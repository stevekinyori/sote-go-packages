package sDocument

import (
	"os"
	"runtime"

	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*  This will return prefix location of tesseracts training data */
func SGetTessdataPrefix() (tessdataPrefix string) {
	sLogger.DebugMethod()

	if tessdataPrefix = os.Getenv("TESSDATA_PREFIX"); tessdataPrefix == "" {
		sgoOs := runtime.GOOS
		// TODO automate this for different operating systems, download if not available and set it manually
		switch sgoOs {
		case "darwin":
			tessdataPrefix = "/opt/homebrew/Cellar/tesseract/4.1.1/share"
		default:
			sLogger.Info("TESSDATA_PREFIX not set:Couldn't determine OS type")
		}

	}

	return
}
