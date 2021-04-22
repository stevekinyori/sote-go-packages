/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other catagories of restrictions}
    * {List of restrictions for the catagory

NOTES:
    {Enter any additional notes that you believe will help the next developer.}
*/
package packages

import (
	// Add imports here

	"os"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sDocument"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	LOGMESSAGEPREFIX = "sDocument_external_test"
	DOCUMENTS        = "documents"
)

// List type's here

var (
// Add Variables here for the file (Remember, they are global)
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

func TestDocumentNew(tPtr *testing.T) {
	if _, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		tPtr.Fail()
	}

	// if _, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode != 109999 {
	// 	tPtr.Errorf("TestNew failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	// }

	if _, soteErr := sDocument.New("", sConfigParams.STAGING, true); soteErr.ErrCode != 200513 {
		tPtr.Errorf("TestNew failed: Expected error code of 200513 got %v", soteErr.FmtErrMsg)
		tPtr.Fail()
	}
}
func TestDownloadDocument(tPtr *testing.T) {
	if dmPtr, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./sDocument/inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestDocumentManager_DownloadDocument failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}
	}

	if dmPtr, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpg", "./sDocument/inbound/", true); soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestDocumentManager_DownloadDocument failed: Expected error code of 109999 got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}
	}

	os.Remove("./sDocument/inbound/upload_test.jpeg")

}
func TestSConvertImageFormat(tPtr *testing.T) {
	// Download test file
	if dmPtr, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./sDocument/inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSConvertImageFormat failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}

		if soteErr.ErrCode == nil {

			if _, soteErr = sDocument.SConvertImageFormat("./sDocument/inbound/upload_test.jpeg",
				"./sDocument/processed/out.pdf"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestSConvertImageFormat failed: Expected error code to be nil got %v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}

			if _, soteErr = sDocument.SConvertImageFormat("./inbound/upload_test.jpeg", "./sDocument/processed/out.pdf"); soteErr.ErrCode != 109999 {
				tPtr.Errorf("TestSConvertImageFormat failed: Expected error code of 109999 got %v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}
	}

	os.Remove("./sDocument/inbound/upload_test.jpeg")
	os.Remove("./sDocument/processed/out.pdf")
}
func TestUploadDocument(tPtr *testing.T) {
	if dmPtr, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./sDocument/inbound/", true); soteErr.ErrCode == nil {

			// Convert downloaded document to PDF
			if _, soteErr = sDocument.SConvertImageFormat("./sDocument/inbound/upload_test.jpeg",
				"./sDocument/processed/out.pdf"); soteErr.ErrCode == nil {

				// Upload converted document
				if _, soteErr = dmPtr.UploadDocument("./sDocument/processed/out.pdf"); soteErr.ErrCode != nil {
					tPtr.Errorf("TestDocumentManager_UploadDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
					tPtr.Fail()
				}

				// Delete uploaded document
				if soteErr.ErrCode == nil {
					if _, soteErr = dmPtr.DeleteDocument("processed/out.pdf", true); soteErr.ErrCode != nil {
					tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
					tPtr.Fail()
				}
				}
			}
		}
	}

	os.Remove("./sDocument/inbound/upload_test.jpeg")
	os.Remove("./sDocument/processed/out.pdf")
}
func TestDeleteDocument(tPtr *testing.T) {
	if dmPtr, soteErr := sDocument.New(DOCUMENTS, sConfigParams.STAGING, true); soteErr.ErrCode == nil {
		if _, soteErr = dmPtr.DownloadDocument("inbound/upload_test.jpeg", "./sDocument/inbound/", true); soteErr.ErrCode != nil {
			tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			tPtr.Fail()
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = sDocument.SConvertImageFormat("./sDocument/inbound/upload_test.jpeg",
				"./sDocument/processed/out.pdf"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestDocumentManager_DeleteDocument failed: Expected error code to be nil got %v ", soteErr.FmtErrMsg)
				tPtr.Fail()
			}
		}

		if soteErr.ErrCode == nil {
			if _, soteErr = dmPtr.UploadDocument("./sDocument/processed/out.pdf"); soteErr.ErrCode != nil {
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

	os.Remove("./sDocument/inbound/upload_test.jpeg")
	os.Remove("./sDocument/processed/out.pdf")
}
