package sDocument

import (
	"fmt"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type DocumentServerManager struct {
	tesseractServer *TesseractServerManager
}

/* New initializes Document Management Server Manager*/
func New() (documentServerPtr *DocumentServerManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tTesseractServer *TesseractServerManager
	)

	if tTesseractServer, soteErr = NewTesseractServer(SGetTessdataPrefix()); soteErr.ErrCode == nil {
		documentServerPtr = &DocumentServerManager{tesseractServer: tTesseractServer}
	}

	return
}

/* ConvertFileFormat writes out the image in the same/different format */
func (dsm *DocumentServerManager) ConvertImageFormat(sourcePath string, targetPath string) (pdfFilePtr *imagick.ImageCommandResult,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var serr error
	imagick.Initialize()
	defer imagick.Terminate() // Memory leak cleanup

	if pdfFilePtr, serr = imagick.ConvertImageCommand([]string{"convert", sourcePath, targetPath}); serr != nil {
		fmt.Printf("Metadata:%v", pdfFilePtr)
	}

	return
}
