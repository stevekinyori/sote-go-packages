/* This will prepare uploaded Sote documents for Optical Character Recognition (OCR)  */
package sDocument

import (
	"os"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gocv.io/x/gocv"
)

type PreprocessManager struct {
	sOriginalImage gocv.Mat
}

func NewPreprocessor(sFilePath string) (preprocessManagerPtr *PreprocessManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, soteErr = CheckIfPathExists(sFilePath); soteErr.ErrCode == nil {
		tOriginalImage := gocv.IMRead(sFilePath, -1) // Load image and return it in original format
		preprocessManagerPtr = &PreprocessManager{sOriginalImage: tOriginalImage}
	}

	return
}

/* CheckIfPathExists determines whether path to uploaded file exists */
func CheckIfPathExists(sFilePath string) (sExistingFilePath os.FileInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var sErr error

	if sExistingFilePath, sErr = os.Stat(sFilePath); os.IsNotExist(sErr) {
		//TODO Add sError support for incorrect file path
		soteErr = sError.GetSError(210400, nil, sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
		//TODO Determine whether program should panic if upload path is invalid
		//panic("sDocument.checkIfPathExists failed")
	}

	return
}

/* CorrectSkew fixes direction and angle of skewed images */
func (pm *PreprocessManager) CorrectSkew() (sGrayScaleImage gocv.Mat, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
	//sBitwiseImage   = gocv.NewMat()
	//tThresholdedImage = gocv.NewMat()
	//sThreshold        float32
	)

	sGrayScaleImage = pm.convertImageToGrayScale()
	////Flip the foreground and background to ensure foreground is now "white" and background is "black"
	//gocv.BitwiseNot(tGrayScaleImage, &sGrayScaleImage)
	//// Threshold image, set all foreground pixels to 255 and all background pixels to 0
	//sThreshold = gocv.Threshold(sGrayScaleImage, &tThresholdedImage, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
	//fmt.Println(sThreshold)

	window := gocv.NewWindow("Hello")
	window.IMShow(sGrayScaleImage)
	window.WaitKey(0)

	return
}

/* convertImageToGrayScale Converts image to grayscale */
func (pm PreprocessManager) convertImageToGrayScale() (sGrayScaleImage gocv.Mat) {
	sLogger.DebugMethod()

	sGrayScaleImage = gocv.NewMat()

	gocv.CvtColor(pm.sOriginalImage, &sGrayScaleImage, gocv.ColorBGRToGray)

	return
}
