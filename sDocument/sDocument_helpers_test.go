package sDocument

/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other categories of restrictions}
    * AWS S3 Bucket must be mounted.
    * Mount point environment variable must be defined (Check constant.go for the name).

NOTES:
    {Enter any additional notes that you believe will help the next developer.}
*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

const (
	// Add Constants here
	TESTCLIENTCOMPANYNAMEONE   = "SILAFRICA KENYA LTD"
	TESTCOMPANYSUPPLIERNAMEONE = "ExxonMobil Petroleum & Chemical BV"
	TESTCLIENTCOMPANYID        = 1
	TESTCLIENTCOMPANYIDSTR     = "1"
	// TESTFILENAMEONE            = "test-invoice.jpeg"
	// TESTS3BUCKETNAME           = "sote-internal-technology-data"
	TESTAPPENVIRONMENT    = "staging"
	TESTMOUNTPOINTENVNAME = MOUNTPOINTENVIRONMENTVARNAME
	TESTFILESFOLDER       = "test-files"
	TESTLOCALFILENAME     = "invoice.jpeg"
	TESTLOCALPDFFILENAME  = "invoice-two.pdf"
	TESTINVALIDFILEPATH   = "INVALID PATH"
)

// List type's here
type UploadRes interface {
	*UploadResponse | map[string][]*UploadResponse
}

type FormDataTest struct {
	fieldName          string
	filepath           string
	testFilename       string
	microserviceIdName string
	microserviceId     int
}

type TestError struct {
	ErrStr string
	Err    error
}

type DeleteTestFile struct {
	s3ClientServerPtr *S3ClientServer
	testForm          *multipart.Form
	filename          string
	filenames         []string
}

var (
	testLocalFilepath = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME}, "/")
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
	testMode = true
	appEnvironment = TESTAPPENVIRONMENT
}

func TestGetDocumentsMountPoint(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Get S3 Bucket mount point location", func(tPtr *testing.T) {
		if _, soteErr = GetDocumentsMountPoint(parentCtx, TESTMOUNTPOINTENVNAME); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("Validate missing mount point environment variable", func(tPtr *testing.T) {
		if _, soteErr = GetDocumentsMountPoint(parentCtx, ""); soteErr.ErrCode != 209100 {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "209100", soteErr.FmtErrMsg)
		}
	})
}
func TestRemoveFile(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Remove file with invalid path", func(tPtr *testing.T) {
		if soteErr = RemoveFile(parentCtx, TESTINVALIDFILEPATH); soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "109999", soteErr.FmtErrMsg)
		}
	})
}
func TestValidateFilepath(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Check invalid path", func(tPtr *testing.T) {
		if _, soteErr = ValidateFilepath(TESTINVALIDFILEPATH); soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "109999", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("Check valid path", func(tPtr *testing.T) {
		if _, soteErr = ValidateFilepath(testLocalFilepath); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "109999", soteErr.FmtErrMsg)
		}
	})
}
func TestPanicService(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("nil error", func(tPtr *testing.T) {
		testMode = true
		if soteErr = panicService(parentCtx, sError.SoteError{}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("non nil error", func(tPtr *testing.T) {
		testMode = true
		if soteErr = panicService(parentCtx, sError.SoteError{ErrCode: 207500}); soteErr.ErrCode != 207500 {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "207500", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("non-testMode error", func(tPtr *testing.T) {
		testMode = false
		if soteErr = panicService(parentCtx, sError.SoteError{ErrCode: 207500}); soteErr.ErrCode != 199999 {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "199999", soteErr.FmtErrMsg)
		}
		testMode = true
	})
}
func TestAmazonTextractErrorHandler(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		err               error
	)

	tPtr.Run("InvalidParameterException Error", func(tPtr *testing.T) {
		err = createCustomErr("InvalidParameterException")
		if soteErr = AmazonTextractErrorHandler(parentCtx, err); soteErr.ErrCode != 200514 {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "200514", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("General Business Service Error", func(tPtr *testing.T) {
		err = createCustomErr("General Error")
		if soteErr = AmazonTextractErrorHandler(parentCtx, err); soteErr.ErrCode != 210599 {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "200514", soteErr.FmtErrMsg)
		}
	})
}

/**
TEST HELPER FUNCTIONS
*/

func createTestFormData(tPtr *testing.T, testName, localFilepath string, isSupportingDoc bool) (writer *multipart.Writer,
	form *multipart.Form) {
	sLogger.DebugMethod()

	tPtr.Helper()

	var (
		fileFieldNameOne = ParentDocumentKey.String()
		filepathOne      = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, localFilepath}, "/")
	)
	if isSupportingDoc {
		fileFieldNameOne = SupportingDocumentsKey.String()
	}

	form, _ = createTestForm(tPtr, &FormDataTest{
		fieldName:    fileFieldNameOne,
		filepath:     filepathOne,
		testFilename: testName,
	})

	return
}

// createTestFormData Will create form-data test
func createTestForm(tPtr *testing.T, fdPtr *FormDataTest) (form *multipart.Form, err error) {
	sLogger.DebugMethod()

	tPtr.Helper()

	var (
		writer         *multipart.Writer
		body           *bytes.Buffer
		contents       []byte
		soteErr        sError.SoteError
		targetFilepath = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, fdPtr.testFilename}, "/")
		part           io.Writer
	)

	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)

	_, _ = os.Create(targetFilepath)
	if part, err = writer.CreateFormFile(fdPtr.fieldName, targetFilepath); err == nil {
		if contents, soteErr = ReadFile(parentCtx, fdPtr.filepath); soteErr.ErrCode == nil {
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

func (te *TestError) Error() string {
	return fmt.Sprintf("Error is %q: err %v", te.ErrStr, te.Err)
}

func createCustomErr(errStr string) error {
	return &TestError{
		ErrStr: errStr,
		Err:    errors.New("unavailable"),
	}
}

func deleteTestFile[V UploadRes](tPtr *testing.T, uploadRes V, delTestFile *DeleteTestFile) {
	tPtr.Helper()

	var (
		tUploadRes  = UploadResponse{}
		tUploadsRes = make(map[string][]*UploadResponse, 2)
		data        []byte
		err         error
	)

	switch reflect.TypeOf(uploadRes) {
	case reflect.TypeOf(&tUploadRes):
		if data, err = json.Marshal(uploadRes); err == nil {
			if err = json.Unmarshal(data, &tUploadRes); err == nil {
				if objKeys := tUploadRes.ObjKeys; objKeys != nil {
					delTestFile.s3ClientServerPtr.DocumentDelete(parentCtx, objKeys.InboundObjectKey)
					delTestFile.s3ClientServerPtr.DocumentDelete(parentCtx, objKeys.ProcessedObjectKey)
				}
				_ = delTestFile.testForm.RemoveAll()
				RemoveFile(parentCtx, strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, delTestFile.filename}, "/"))
			}
		}
	default:
		if data, err = json.Marshal(uploadRes); err == nil {
			if err = json.Unmarshal(data, &tUploadsRes); err == nil {
				for _, v := range tUploadsRes {
					for _, uploadResPtr := range v {
						if objKeys := uploadResPtr.ObjKeys; objKeys != nil {
							delTestFile.s3ClientServerPtr.DocumentDelete(parentCtx, objKeys.InboundObjectKey)
							delTestFile.s3ClientServerPtr.DocumentDelete(parentCtx, objKeys.ProcessedObjectKey)
						}
					}
				}
			}
		}
		_ = delTestFile.testForm.RemoveAll()
		for _, filename := range delTestFile.filenames {
			RemoveFile(parentCtx, strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, filename}, "/"))
		}
	}
}

func copyTestDocument(tPtr *testing.T, filename string, useProcessedFolder, usePDF bool) (sourceFilepath, targetFilepath string, keys *ObjectKeys,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	tPtr.Helper()
	var (
		s3ClientServerPtr *S3ClientServer
		inboundFilepath   string
		processedFilepath string
	)
	if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
		AppConfigName:        sConfigParams.DOCUMENTS,
		MountPointEnvVarName: TESTMOUNTPOINTENVNAME,
		ClientCompanyId:      TESTCLIENTCOMPANYID,
		AppEnvironment:       TESTAPPENVIRONMENT,
		TestMode:             testMode,
	}); soteErr.ErrCode == nil {
		keys = GetObjectKeys(filename, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyId))
		inboundFilepath, processedFilepath, _ = s3ClientServerPtr.getMountPointFilepath(keys)

		switch usePDF {
		case true:
			sourceFilepath = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALPDFFILENAME}, "/")
		case false:
			sourceFilepath = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME}, "/")

		}

		switch useProcessedFolder {
		case true:
			targetFilepath = processedFilepath
		case false:
			targetFilepath = inboundFilepath
		}
		_, soteErr = s3ClientServerPtr.DocumentCopy(parentCtx, sourceFilepath, targetFilepath)
	}

	return
}
