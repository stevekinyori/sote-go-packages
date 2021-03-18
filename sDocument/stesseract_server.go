package sDocument

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	tesseract "gopkg.in/GeertJohan/go.tesseract.v1"
	"path/filepath"
	"strings"
	"sync"
)

/* Limits characters tesseract is looking for */
const (
	WHITELIST = `!"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~`+"`"
	)

type DocumentManager struct {
	SManager *tesseract.Tess
	stessdataPrefix string

	sync.Mutex

}

func New(tessdataLocation string)  (pdocumentManager *DocumentManager, soteError sError.SoteError){
	sLogger.DebugMethod()
	var sinstance *tesseract.Tess

	pdocumentManager = &DocumentManager{stessdataPrefix: tessdataLocation}

	if sinstance, soteError = pdocumentManager.connect(); soteError.ErrCode == nil{
		pdocumentManager = &DocumentManager{SManager: sinstance}
	}

	return
}



/*
	This will connect to a new tesseract instance and point it to tessdata location.
*/
func (dm *DocumentManager) connect() (pti *tesseract.Tess, soteErr sError.SoteError)  {
	sLogger.DebugMethod()
	var err error

	// Create a new tesseract instance
	pti, err = tesseract.NewTess(filepath.Join(dm.stessdataPrefix, "tessdata"), "eng")
	if err != nil {
		if strings.Contains(err.Error(), "could not initiate new Tess instance"){
			soteErr = sError.GetSError(209100, sError.BuildParams([]string{"TESSDATA_PREFIX"}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}else {
			var errDetails = make(map[string]string)
			errDetails, soteErr = sError.ConvertErr(err)
			if soteErr.ErrCode != nil {
				sLogger.Info(soteErr.FmtErrMsg)
				panic("sError.ConvertErr Failed")
			}
			sLogger.Info(sError.GetSError(210400, nil, errDetails).FmtErrMsg)
			panic("sDocument.connect Failed")
		}


	}

	return
}



/*
	This will close connection to tesseract instance
*/
func (dm *DocumentManager) close()  {
	sLogger.DebugMethod()

	dm.SManager.Close()

	dm = nil

	return

}