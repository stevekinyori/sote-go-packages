package sDocument

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	leptonica "gopkg.in/GeertJohan/go.leptonica.v1"
	tesseract "gopkg.in/GeertJohan/go.tesseract.v1"
	"path/filepath"
	"strings"
	"sync"
)

/* Limits characters tesseract is looking for */
const (
	WHITELIST = ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~` + "`"
)

type DocumentManager struct {
	SManager        *tesseract.Tess
	stessdataPrefix string
	pPix            *leptonica.Pix

	sync.Mutex
}

func New(tessdataLocation string) (pdocumentManager *DocumentManager, soteError sError.SoteError) {
	sLogger.DebugMethod()
	var sinstance *tesseract.Tess

	pdocumentManager = &DocumentManager{stessdataPrefix: tessdataLocation}

	if sinstance, soteError = pdocumentManager.connect(); soteError.ErrCode == nil {
		pdocumentManager = &DocumentManager{SManager: sinstance}
	}

	return
}

/*
	This will connect to a new tesseract instance and point it to tessdata location.
*/
func (dm *DocumentManager) connect() (pti *tesseract.Tess, soteError sError.SoteError) {
	sLogger.DebugMethod()
	var serr error

	// Create a new tesseract instance
	pti, serr = tesseract.NewTess(filepath.Join(dm.stessdataPrefix, "tessdata"), "eng")
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

func (dm *DocumentManager) GetTextFromDocument(sfilename string) (stext string, soteError sError.SoteError) {
	sLogger.DebugMethod()

	// open a new Pix from file with leptonica
	if ppix, serr := leptonica.NewPixFromFile(sfilename); serr == nil {
		// set the page seg mode to autodetect
		dm.SManager.SetPageSegMode(tesseract.PSM_AUTO_OSD)

		// setup a whitelist of all basic ascii
		if soteError = dm.setWhitelist(); soteError.ErrCode == nil {
			// set the image to the tesseract instance
			dm.SManager.SetImagePix(ppix)
			stext = dm.SManager.Text()
		}

	} else {
		soteError = sError.GetSError(209110, sError.BuildParams([]string{serr.Error()}), sError.EmptyMap)
		sLogger.Info(soteError.FmtErrMsg)

	}

	return
}

/* setWhitelist restricts the TesseractOCR function to a set of pre-defined (white-listed) characters */
func (dm *DocumentManager) setWhitelist() (soteError sError.SoteError) {
	sLogger.DebugMethod()

	if serr := dm.SManager.SetVariable("tessedit_char_whitelist", WHITELIST); serr != nil {
		soteError = sError.GetSError(209110, sError.BuildParams([]string{"tessedit_char_whitelist"}), sError.EmptyMap)
		sLogger.Info(soteError.FmtErrMsg)
	}

	return
}

/*
	This will close connection to tesseract instance
*/
func (dm *DocumentManager) close() {
	sLogger.DebugMethod()

	dm.SManager.Close()

	dm = nil

	return

}
