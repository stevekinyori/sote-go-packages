package sDocument

import (
	"fmt"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	tesseract "gopkg.in/GeertJohan/go.tesseract.v1"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type DocumentManager struct {
	Manager *tesseract.Tess

	sync.Mutex

}

func New()  (documentManager *DocumentManager, soteError sError.SoteError){
	sLogger.DebugMethod()
	var instance *tesseract.Tess

	if instance, soteError = documentManager.connect(); soteError.ErrCode == nil{
		documentManager = &DocumentManager{Manager: instance}
	}

	return
}

/*  This will return location of tesseracts training data */
func (dm *DocumentManager) getTessdataLocation()  (tessdataPrefix string){
	if tessdataPrefix = os.Getenv("TESSDATA_PREFIX"); tessdataPrefix == ""{
		goos := runtime.GOOS
		//TODO automate this for different operating systems
		switch goos {
		case "darwin":
			fmt.Println("Mac os")
			tessdataPrefix = "/opt/homebrew/Cellar/tesseract/4.1.1/share"
		}

	}

	return
}

/*
	This will connect to a new tesseract instance and point it to tessdata location.
*/
func (dm *DocumentManager) connect() (ti *tesseract.Tess, soteErr sError.SoteError)  {
	sLogger.DebugMethod()
	var err error

	// Create a new tesseract instance
	if ti, err = tesseract.NewTess(filepath.Join(dm.getTessdataLocation(), "tessdata"), "eng"); err != nil {
		sLogger.Info(err.Error())
		dm.Manager.Close()
	}


	return
}

/*
	This will close connection to tesseract instance
*/
func (dm *DocumentManager) close()  {
	sLogger.DebugMethod()

	dm.Manager.Close()

	dm = nil

	return

}