/*
This will create an tesseract instance that is used for Optical Character Recognitions(OCR) of Sote Documents.
*/
package sDocument

import (
	"path/filepath"
	"strings"
	"sync"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	leptonica "gopkg.in/GeertJohan/go.leptonica.v1"
	tesseract "gopkg.in/GeertJohan/go.tesseract.v1"

)

const (
	/* Limits characters tesseract is looking for */
	WHITELIST = ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~` + "`"
)

type TesseractServerManager struct {
	Manager        *tesseract.Tess
	tessdataPrefix string

	sync.Mutex
}

/* Creates a new tesseract instance for OCR. */
func NewTesseractServer(tessdataPrefix string) (tesseractServerManagerPtr *TesseractServerManager, soteError sError.SoteError) {
	sLogger.DebugMethod()
	var sinstance *tesseract.Tess

	tesseractServerManagerPtr = &TesseractServerManager{tessdataPrefix: tessdataPrefix}

	if sinstance, soteError = tesseractServerManagerPtr.connect(); soteError.ErrCode == nil {
		tesseractServerManagerPtr = &TesseractServerManager{Manager: sinstance}
	}

	return
}

/* GetTextFromFile performs Optical Character Recognition on a file/image. It returns resulting text and a SoteError */
func (tsm *TesseractServerManager) GetTextFromFile(sfilename string) (stext string, soteError sError.SoteError) {
	sLogger.DebugMethod()

	// open a new Pix from file with leptonica
	if ppix, serr := leptonica.NewPixFromFile(sfilename); serr == nil {
		// set the page seg mode to autodetect
		tsm.Manager.SetPageSegMode(tesseract.PSM_AUTO_OSD)

		// setup a whitelist of all basic ascii
		if soteError = tsm.setWhitelist(); soteError.ErrCode == nil {
			// set the image to the tesseract instance
			tsm.Manager.SetImagePix(ppix)
			stext = tsm.Manager.Text()
		}

	} else {
		soteError = sError.GetSError(209110, sError.BuildParams([]string{serr.Error()}), sError.EmptyMap)
		sLogger.Info(soteError.FmtErrMsg)

	}

	return
}

/*
	This will connect to a new tesseract instance and point it to tessdata location.
*/
func (tsm *TesseractServerManager) connect() (pti *tesseract.Tess, soteError sError.SoteError) {
	sLogger.DebugMethod()
	var serr error

	// Create a new tesseract instance
	pti, serr = tesseract.NewTess(filepath.Join(tsm.tessdataPrefix, "tessdata"), "eng")
	if serr != nil {
		if strings.Contains(serr.Error(), "could not initiate new Tess instance") {
			soteError = sError.GetSError(209100, sError.BuildParams([]string{"TESSDATA_PREFIX"}), sError.EmptyMap)
			sLogger.Info(soteError.FmtErrMsg)
		} else {
			var errDetails = make(map[string]string)
			errDetails, soteError = sError.ConvertErr(serr)
			if soteError.ErrCode != nil {
				sLogger.Info(soteError.FmtErrMsg)
				panic("sError.ConvertErr Failed")
			}
			sLogger.Info(sError.GetSError(210400, nil, errDetails).FmtErrMsg)
			panic("sDocument.connect Failed")
		}

	}

	return
}

/* setWhitelist restricts the TesseractOCR function to a set of pre-defined (white-listed) characters */
func (tsm *TesseractServerManager) setWhitelist() (soteError sError.SoteError) {
	sLogger.DebugMethod()

	if serr := tsm.Manager.SetVariable("tessedit_char_whitelist", WHITELIST); serr != nil {
		soteError = sError.GetSError(209110, sError.BuildParams([]string{"tessedit_char_whitelist"}), sError.EmptyMap)
		sLogger.Info(soteError.FmtErrMsg)
	}

	return
}

/*
	This will close connection to tesseract instance
*/
func (tsm *TesseractServerManager) close() {
	sLogger.DebugMethod()

	tsm.Manager.Close()

	tsm = nil

	return

}
