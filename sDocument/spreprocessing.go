/*
	This will prepare uploaded Sote documents for Optical Character Recognition (OCR).

	The algorithm used to automatically detect and correct text on a skewed image is found here:
		https://www.pyimagesearch.com/2017/02/20/text-skew-correction-opencv-python/

	RESTRICTIONS:
		Preprocessor Manager functions:
		* Must has OpenCV 4 computer vision library installed in the system to be able to use GoCV.
*/
package sDocument

import (
	"os"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gocv.io/x/gocv"
)

type PreprocessManager struct {
	sOriginalImage    gocv.Mat
	sThresholdedImage gocv.Mat
	sThreshold        float32
}

func NewPreprocessor(sFilePath string) (preprocessManagerPtr *PreprocessManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, soteErr = CheckIfPathExists(sFilePath); soteErr.ErrCode == nil {
		tOriginalImage := gocv.IMRead(sFilePath, gocv.IMReadColor) // Load image and return it in original format
		preprocessManagerPtr = &PreprocessManager{sOriginalImage: tOriginalImage}
	}

	return
}

/* CheckIfPathExists determines whether path to uploaded file exists */
func CheckIfPathExists(sFilePath string) (sExistingFilePath os.FileInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var sErr error

	if sExistingFilePath, sErr = os.Stat(sFilePath); os.IsNotExist(sErr) {
		// TODO Add sError support for incorrect file path
		soteErr = sError.GetSError(199999, sError.BuildParams([]string{"File path (" + sFilePath + ") doesn't exist"}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
		// TODO Determine whether program should panic if upload path is invalid
		// panic("sDocument.checkIfPathExists failed")
	}

	return
}

/* CorrectSkew fixes direction and angle of skewed images */
func (pm *PreprocessManager) CorrectSkew() (sGrayScaleImage gocv.Mat, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		sBitwiseImage     = gocv.NewMat()
		tThresholdedImage = gocv.NewMat()
		tThreshold        float32
	)

	sGrayScaleImage = pm.convertImageToGrayScale()
	sBitwiseImage = pm.convertGrayscaleImageToBitwise(sGrayScaleImage)
	tThreshold, tThresholdedImage = pm.thresholdImage(sBitwiseImage)
	pm.sThreshold = tThreshold
	pm.sThresholdedImage = tThresholdedImage
	// gocv.IMWrite("gray.png", sBitwiseImage)
	// window := gocv.NewWindow("Hello")
	// window.IMShow(sGrayScaleImage)
	// window.WaitKey(0)

	return
}

/* convertImageToGrayScale Converts image to grayscale */
func (pm *PreprocessManager) convertImageToGrayScale() (sGrayScaleImage gocv.Mat) {
	sLogger.DebugMethod()

	sGrayScaleImage = gocv.NewMat()

	gocv.CvtColor(pm.sOriginalImage, &sGrayScaleImage, gocv.ColorBGRToGray)

	return
}

/*
	convertGrayscaleImageToBitwise flips the foreground and background to ensure foreground is "white" and background is "black".
	This makes all text dark and background light.
*/
func (pm *PreprocessManager) convertGrayscaleImageToBitwise(sGrayScaleImage gocv.Mat) (sBitwiseImage gocv.Mat) {
	sLogger.DebugMethod()

	sBitwiseImage = gocv.NewMat()

	gocv.BitwiseNot(sGrayScaleImage, &sBitwiseImage)

	return
}

/*
	thresholdImage thresholds bitwise image, set all foreground pixels to 255(black) and all background pixels to 0(white).
 	This makes all text light and background dark.
*/
func (pm *PreprocessManager) thresholdImage(sBitwiseImage gocv.Mat) (sThreshold float32, tThresholdedImage gocv.Mat) {
	sLogger.DebugMethod()

	tThresholdedImage = gocv.NewMat()

	sThreshold = gocv.Threshold(sBitwiseImage, &tThresholdedImage, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	return
}
