package sDocument

import (
	"os"
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	DOCUMENTS = "documents"
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

func TestNew(tPtr *testing.T) {
	if _, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		tPtr.Fail()
	}

	// if _, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode != 109999 {
	// 	tPtr.Errorf("TestNew failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	// }
}
func TestDocumentManager_DownloadDocument(tPtr *testing.T) {
	if dmPtr, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestDocumentManager_DownloadDocument failed: Expected error code  to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}
	}

	if dmPtr, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpg", "./inbound/", true); soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestDocumentManager_DownloadDocument failed: Expected error code of 109999 got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}
	}

	os.Remove("./inbound/upload_test.jpeg")
}
func TestSConvertImageFormat(tPtr *testing.T) {
	if dmPtr, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		// Download test file
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSConvertImageFormat failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = SConvertImageFormat(strings.Join([]string{GetFullDirectoryPath(), "inbound/upload_test.jpeg"}, "/"),
				strings.Join([]string{GetFullDirectoryPath(), "processed/out.pdf"}, "/")); soteErr.ErrCode != nil {
				tPtr.Errorf("TestSConvertImageFormat failed: Expected error code to be nil:%v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}

			if _, soteErr = SConvertImageFormat("inbound/upload_test.jpg",
				"../processed/out.pdf"); soteErr.ErrCode != 109999 {
				tPtr.Errorf("TestSConvertImageFormat failed: Expected error code of 109999 got %v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}
	}

	os.Remove("./inbound/upload_test.jpeg")
	os.Remove("./processed/out.pdf")
}
func TestDocumentManager_UploadDocument(tPtr *testing.T) {
	if dmPtr, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		// Download test file
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestDocumentManager_UploadDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = SConvertImageFormat(strings.Join([]string{GetFullDirectoryPath(), "inbound/upload_test.jpeg"}, "/"),
				strings.Join([]string{GetFullDirectoryPath(), "processed/out.pdf"}, "/")); soteErr.ErrCode != nil {
				tPtr.Errorf("TestDocumentManager_UploadDocument failed: Expected error code to be nil:%v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = dmPtr.UploadDocument(strings.Join([]string{GetFullDirectoryPath(), "processed/out.pdf"}, "/")); soteErr.ErrCode != nil {
				tPtr.Errorf("TestDocumentManager_UploadDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}
	}

	os.Remove("./inbound/upload_test.jpeg")
	os.Remove("./processed/out.pdf")
}
func TestDocumentManager_DeleteDocument(tPtr *testing.T) {
	if dmPtr, soteErr := New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = SConvertImageFormat("./inbound/upload_test.jpeg",
				"./processed/out.pdf"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil:%v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = dmPtr.UploadDocument("./processed/out.pdf"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = dmPtr.DeleteDocument("processed/out.pdf", true); soteErr.ErrCode != nil {
				tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}
	}

	os.Remove("./inbound/upload_test.jpeg")
	os.Remove("./processed/out.pdf")
}
