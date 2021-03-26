/* This will prepare uploaded Sote documents for Optical Character Recognition (OCR)  */
package sDocument

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gocv.io/x/gocv"
	"os"
)

type PreprocessManager struct {
}

func NewPreprocessor() (preprocessManagerPtr *PreprocessManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	preprocessManagerPtr = &PreprocessManager{}

	return
}

/* CorrectSkew fixes direction and angle of skewed images */
func (pm *PreprocessManager) CorrectSkew(sFilePath string) (sGrayScaleImage gocv.Mat, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tGrayScaleImage gocv.Mat = gocv.NewMat()

	if _, soteErr = pm.checkIfPathExists(sFilePath); soteErr.ErrCode == nil {
		sOriginalImage := gocv.IMRead(sFilePath, -1) // Load image and return it in original format
		sGrayScaleImage = gocv.NewMat() // Create a new empty mat
		//convert the image to grayscale
		gocv.CvtColor(sOriginalImage, &tGrayScaleImage, gocv.ColorBGRToGray)
		//Flip the foreground and background to ensure foreground is now "white" and background is "black"
		gocv.BitwiseNot(tGrayScaleImage, &sGrayScaleImage)

		window := gocv.NewWindow("Hello")
		window.IMShow(sGrayScaleImage)
		window.WaitKey(0)
	}

	return
}

/* checkIfPathExists determines whether path to uploaded file exists */
func (pm *PreprocessManager) checkIfPathExists(sFilePath string) (sExistingFilePath os.FileInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var sErr error

	if sExistingFilePath, sErr = os.Stat(sFilePath); os.IsNotExist(sErr) {
		//TODO Add sError support for incorrect file path
		soteErr = sError.GetSError(210400, nil, sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
		//panic("sDocument.checkIfPathExists failed")
	}

	return
}
