package tests

/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS}

    {Other categories of restrictions}
    * AWS S3 Bucket must be mounted.
    * Mount point environment variable must be defined (Check constant.go for the name).

NOTES:
    {Enter any additional notes that you believe will help the next developer.}
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sDocument"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	TESTMODE           = true
	TESTAPPENVIRONMENT = "staging"
	// TESTCLIENTCOMPANYID            = 1
	TESTCLIENTCOMPANYNAME          = "Test Client Company One"
	TESTDOCUMENTSMOUNTPOINTENVNAME = sDocument.DOCUMENTSMOUNTPOINTENVIRONMENTVARNAME
	TESTFILESFOLDER                = "test-files"
	TESTLOCALFILENAME              = "invoice.jpeg"
	TESTLOCALPDFFILENAME           = "invoice.jpeg"
	TESTINVALIDFILEPATH            = "INVALID PATH"
)

type formDataTest struct {
	fieldName          string
	filepath           string
	testFilename       string
	microserviceIdName string
	microserviceId     int
}

type UploadRes interface {
	*sDocument.UploadResponse | map[string][]*sDocument.UploadResponse
}

type DeleteTestFile struct {
	s3ClientServerPtr *sDocument.S3ClientServer
	testForm          *multipart.Form
	filename          string
	filenames         []string
}

var (
	sDocumentCtx = context.Background()
)

func init() {
	sLogger.SetLogMessagePrefix("packages")
}

func TestNewS3ClientServer(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Initialize S3 Client Server", func(tPtr *testing.T) {
		if _, soteErr = sDocument.NewS3ClientServer(sDocumentCtx, &sDocument.DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			AppEnvironment:    TESTAPPENVIRONMENT,
			ClientCompanyName: TESTCLIENTCOMPANYNAME,
			TestMode:          TESTMODE,
		}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("Initialize S3 Client Server Using Invalid Application Name", func(tPtr *testing.T) {
		if _, soteErr = sDocument.NewS3ClientServer(sDocumentCtx, &sDocument.DocumentParams{
			AppConfigName:     "INVALIDAPPNAME",
			AppEnvironment:    TESTAPPENVIRONMENT,
			ClientCompanyName: TESTCLIENTCOMPANYNAME,
			TestMode:          TESTMODE,
		}); soteErr.ErrCode != sError.ErrItemNotFound {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrItemNotFound,
				soteErr.FmtErrMsg)
		}
	})
}
func TestDocumentUpload(tPtr *testing.T) {
	var (
		function, _, _, _       = runtime.Caller(0)
		testName                = runtime.FuncForPC(function).Name()
		soteErr                 sError.SoteError
		err                     error
		s3ClientServerPtr       *sDocument.S3ClientServer
		formFiles               map[string][]*multipart.FileHeader
		fileFieldNameOne        = sDocument.ParentDocumentKey.String()
		testSingleFilenameOne   = "single-upload.jpeg"
		testSingleFilenameTwo   = "1"
		testMultipleFilenameOne = "multiple-upload-one.jpeg"
		testMultipleFilenameTwo = "multiple-upload-two.jpeg"
		uploadResponseOne       *sDocument.UploadResponse
		uploadResponseTwo       *sDocument.UploadResponse
		uploadsResponse         map[string][]*sDocument.UploadResponse
		openedFile              multipart.File
		singleDocForm           *multipart.Form
		multipleDocForm         *multipart.Form
		tMultipleDocForm        *multipart.Form
	)

	tPtr.Cleanup(func() {
		if uploadResponseOne != nil {
			deleteTestFile(tPtr, uploadResponseOne, &DeleteTestFile{
				s3ClientServerPtr: s3ClientServerPtr,
				testForm:          singleDocForm,
				filename:          testSingleFilenameOne,
			})
		}

		if uploadResponseTwo != nil {
			deleteTestFile(tPtr, uploadResponseTwo, &DeleteTestFile{
				s3ClientServerPtr: s3ClientServerPtr,
				testForm:          singleDocForm,
				filename:          testSingleFilenameTwo,
			})
		}

		if len(uploadsResponse) > 0 {
			deleteTestFile(tPtr, uploadsResponse, &DeleteTestFile{
				s3ClientServerPtr: s3ClientServerPtr,
				testForm:          multipleDocForm,
				filenames:         []string{testMultipleFilenameOne, testMultipleFilenameTwo},
			})
		}
	})

	tPtr.Run("Upload Single Document", func(tPtr *testing.T) {
		_, singleDocForm = createTestFormData(tPtr, testSingleFilenameOne, TESTLOCALFILENAME, false)
		formFiles = singleDocForm.File

		if s3ClientServerPtr, soteErr = sDocument.NewS3ClientServer(sDocumentCtx, &sDocument.DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			AppEnvironment:    TESTAPPENVIRONMENT,
			ClientCompanyName: TESTCLIENTCOMPANYNAME,
			TestMode:          TESTMODE,
		}); soteErr.ErrCode == nil {
			tPtr.Run("Using Valid files", func(tPtr *testing.T) {
				for _, file := range formFiles[fileFieldNameOne] {
					if openedFile, err = file.Open(); err != nil {
						break
					} else {
						if uploadResponseOne, soteErr = s3ClientServerPtr.SingleDocumentUpload(sDocumentCtx, openedFile,
							file.Filename); soteErr.ErrCode != nil {
							tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
						}
					}
				}
			})

			tPtr.Run("Using Invalid Filename", func(tPtr *testing.T) {
				for _, file := range formFiles[fileFieldNameOne] {
					if openedFile, err = file.Open(); err != nil {
						break
					} else {
						if uploadResponseTwo, soteErr = s3ClientServerPtr.SingleDocumentUpload(sDocumentCtx, openedFile,
							"1"); soteErr.ErrCode != nil && uploadResponseTwo != nil {
							tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
						}
					}
				}
			})
		}
	})

	tPtr.Run("Upload Multiple Documents", func(tPtr *testing.T) {
		_, multipleDocForm = createTestFormData(tPtr, testMultipleFilenameOne, TESTLOCALFILENAME, false)
		_, tMultipleDocForm = createTestFormData(tPtr, testMultipleFilenameTwo, TESTLOCALFILENAME, true)
		formFiles = multipleDocForm.File
		formFiles[sDocument.SupportingDocumentsKey.String()] = tMultipleDocForm.File[sDocument.SupportingDocumentsKey.String()]

		if s3ClientServerPtr, soteErr = sDocument.NewS3ClientServer(sDocumentCtx, &sDocument.DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			AppEnvironment:    TESTAPPENVIRONMENT,
			ClientCompanyName: TESTCLIENTCOMPANYNAME,
			TestMode:          TESTMODE,
			FormFiles:         formFiles,
		}); soteErr.ErrCode == nil {
			if uploadsResponse, _, soteErr = s3ClientServerPtr.MultipleDocumentsUpload(sDocumentCtx); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		}
	})
}
func TestDirectDocumentUpload(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		s3ClientServerPtr *sDocument.S3ClientServer
		contents          = make([]byte, 0)
		filenameOne       = "single-direct-image-upload.jpeg"
		filenameTwo       = "single-direct-file-upload.pdf"
		keysOne           = new(sDocument.ObjectKeys)
		keysTwo           = new(sDocument.ObjectKeys)
	)

	tPtr.Cleanup(func() {
		if keysOne != nil {
			s3ClientServerPtr.DocumentDelete(sDocumentCtx, keysOne.ProcessedObjectKey)
		}

		if keysTwo != nil {
			s3ClientServerPtr.DocumentDelete(sDocumentCtx, keysTwo.ProcessedObjectKey)
		}
	})

	if s3ClientServerPtr, soteErr = sDocument.NewS3ClientServer(sDocumentCtx, &sDocument.DocumentParams{
		AppConfigName:     sConfigParams.DOCUMENTS,
		AppEnvironment:    TESTAPPENVIRONMENT,
		ClientCompanyName: TESTCLIENTCOMPANYNAME,
		TestMode:          TESTMODE,
	}); soteErr.ErrCode == nil {
		tPtr.Run("Using a Valid Image Ignoring Mount Point", func(tPtr *testing.T) {
			if contents, soteErr = sDocument.ReadFile(sDocumentCtx,
				strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME},
					"/")); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			keysOne = sDocument.GetObjectKeys(filenameOne, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			if soteErr = s3ClientServerPtr.DocumentUpload(sDocumentCtx, keysOne.ProcessedObjectKey, contents,
				sDocument.GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		})

		tPtr.Run("Using a Valid File Ignoring Mount Point", func(tPtr *testing.T) {
			if contents, soteErr = sDocument.ReadFile(sDocumentCtx, strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER,
				TESTLOCALPDFFILENAME},
				"/")); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			keysTwo = sDocument.GetObjectKeys(filenameTwo, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			if soteErr = s3ClientServerPtr.DocumentUpload(sDocumentCtx, keysTwo.ProcessedObjectKey, contents,
				sDocument.GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		})
	}
}
func TestDocumentPreSignedURL(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		s3ClientServerPtr *sDocument.S3ClientServer
		keys              *sDocument.ObjectKeys
		documentLinks     *sDocument.DocumentLinks
		contents          []byte
		sourceFilepath    string
		filename          = "test-presigned-url.jpeg"
	)

	tPtr.Cleanup(func() {
		if keys != nil {
			s3ClientServerPtr.DocumentDelete(sDocumentCtx, keys.ProcessedObjectKey)
			s3ClientServerPtr.DocumentDelete(sDocumentCtx, keys.InboundObjectKey)
		}
	})

	if s3ClientServerPtr, soteErr = sDocument.NewS3ClientServer(sDocumentCtx, &sDocument.DocumentParams{
		AppConfigName:     sConfigParams.DOCUMENTS,
		AppEnvironment:    TESTAPPENVIRONMENT,
		ClientCompanyName: TESTCLIENTCOMPANYNAME,
		TestMode:          TESTMODE,
	}); soteErr.ErrCode == nil {
		tPtr.Run("Generate presigned URL using valid document-link", func(tPtr *testing.T) {
			sourceFilepath = strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME}, "/")
			keys = sDocument.GetObjectKeys(filename, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			// Read contents of file to be uploaded
			if contents, soteErr = sDocument.ReadFile(parentCtx, sourceFilepath); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
			// Upload source document
			if s3ClientServerPtr.DocumentUpload(parentCtx, keys.ProcessedObjectKey, contents,
				sDocument.GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			if documentLinks, soteErr = sDocument.GetDocumentLinks(sDocumentCtx, s3ClientServerPtr.BucketName, keys); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			s3ClientServerPtr.DocumentParamsPtr.DocumentsLink = documentLinks.ProcessedDocumentLink
			if _, soteErr = s3ClientServerPtr.DocumentPreSignedURL(sDocumentCtx, 6); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		})

		tPtr.Run("Check Invalid Pre-Signed Document URL", func(tPtr *testing.T) {
			if _, soteErr = sDocument.ValidatePreSignedDocumentURL(TESTINVALIDFILEPATH); soteErr.ErrCode != sError.ErrItemNotFound {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrItemNotFound, soteErr.FmtErrMsg)
			}
		})
	}
}

/**
TEST HELPER FUNCTIONS
*/
// todo: Move these functions to a utility file to allow importing
func createTestFormData(tPtr *testing.T, testName, localFilepath string, isSupportingDoc bool) (writer *multipart.Writer,
	form *multipart.Form) {
	sLogger.DebugMethod()

	tPtr.Helper()

	var (
		fileFieldNameOne = sDocument.ParentDocumentKey.String()
		filepathOne      = strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER, localFilepath}, "/")
	)
	if isSupportingDoc {
		fileFieldNameOne = sDocument.SupportingDocumentsKey.String()
	}

	form, _ = createTestForm(tPtr, &formDataTest{
		fieldName:    fileFieldNameOne,
		filepath:     filepathOne,
		testFilename: testName,
	})

	return
}

// createTestFormData Will create form-data test
func createTestForm(tPtr *testing.T, fdPtr *formDataTest) (form *multipart.Form, err error) {
	sLogger.DebugMethod()

	tPtr.Helper()

	var (
		writer         *multipart.Writer
		body           *bytes.Buffer
		contents       []byte
		soteErr        sError.SoteError
		targetFilepath = strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER, fdPtr.testFilename}, "/")
		part           io.Writer
	)

	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)

	_, _ = os.Create(targetFilepath)
	if part, err = writer.CreateFormFile(fdPtr.fieldName, targetFilepath); err == nil {
		if contents, soteErr = sDocument.ReadFile(sDocumentCtx, fdPtr.filepath); soteErr.ErrCode == nil {
			if _, err = part.Write(contents); err == nil {
				// close the writer before making the request
				writer.Close()

				r := multipart.NewReader(body, writer.Boundary())
				form, err = r.ReadForm(25)
			}
		}
	}

	return
}

func deleteTestFile[V UploadRes](tPtr *testing.T, uploadRes V, delTestFile *DeleteTestFile) {
	tPtr.Helper()

	var (
		tUploadRes  = sDocument.UploadResponse{}
		tUploadsRes = make(map[string][]*sDocument.UploadResponse, 2)
		data        []byte
		err         error
	)

	switch reflect.TypeOf(uploadRes) {
	case reflect.TypeOf(&tUploadRes):
		if data, err = json.Marshal(uploadRes); err == nil {
			if err = json.Unmarshal(data, &tUploadRes); err == nil {
				if objKeys := tUploadRes.ObjKeys; objKeys != nil {
					delTestFile.s3ClientServerPtr.DocumentDelete(sDocumentCtx, objKeys.InboundObjectKey)
					delTestFile.s3ClientServerPtr.DocumentDelete(sDocumentCtx, objKeys.ProcessedObjectKey)
				}
				_ = delTestFile.testForm.RemoveAll()
				sDocument.RemoveFile(sDocumentCtx,
					strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER, delTestFile.filename}, "/"))
			}
		}
	default:
		if data, err = json.Marshal(uploadRes); err == nil {
			if err = json.Unmarshal(data, &tUploadsRes); err == nil {
				for _, v := range tUploadsRes {
					for _, uploadResPtr := range v {
						if objKeys := uploadResPtr.ObjKeys; objKeys != nil {
							delTestFile.s3ClientServerPtr.DocumentDelete(sDocumentCtx, objKeys.InboundObjectKey)
							delTestFile.s3ClientServerPtr.DocumentDelete(sDocumentCtx, objKeys.ProcessedObjectKey)
						}
					}
				}
			}
		}
		_ = delTestFile.testForm.RemoveAll()
		for _, filename := range delTestFile.filenames {
			sDocument.RemoveFile(sDocumentCtx, strings.Join([]string{sDocument.GetFullDirectoryPath(), TESTFILESFOLDER, filename}, "/"))
		}
	}
}
