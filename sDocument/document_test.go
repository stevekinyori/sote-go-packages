package sDocument

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

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
// todo: Add negative tests

const (
// Add Constants here
)

// List type's here

var (
	parentCtx = context.Background()
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
	testMode = true
	appEnvironment = TESTAPPENVIRONMENT
}

func TestNewS3ClientServer(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Init S3 Client Server", func(tPtr *testing.T) {
		if _, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			AppEnvironment:    TESTAPPENVIRONMENT,
			ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
			TestMode:          testMode,
		}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("Init S3 Client Server Using Invalid Application Name", func(tPtr *testing.T) {
		if _, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
			AppConfigName:     "INVALIDAPPNAME",
			ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
			AppEnvironment:    TESTAPPENVIRONMENT,
			TestMode:          testMode,
		}); soteErr.ErrCode != sError.ErrItemNotFound {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrItemNotFound, soteErr.FmtErrMsg)
		}
	})
}
func TestGetObjectKeys(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		inboundObjectKey  = strings.Join([]string{BUCKETDOCSUBFOLDER, INBOUNDFOLDER, TESTAPPENVIRONMENT, TESTCLIENTCOMPANYNAMEONE, "testOne.jpeg"},
			"/")
	)

	tPtr.Run("Get Object Keys From Parent Document", func(tPtr *testing.T) {
		if objectKeys := GetObjectKeys("testOne.jpeg", fmt.Sprint(TESTCLIENTCOMPANYNAMEONE)); objectKeys.InboundObjectKey != inboundObjectKey {
			tPtr.Errorf("%v Failed: Expected object key to be %v but got %v", testName, inboundObjectKey, objectKeys.InboundObjectKey)
		}
	})
}
func TestConvertInterfaceToByteSlice(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Valid interface", func(tPtr *testing.T) {
		if _, soteErr = ConvertInterfaceToByteSlice("valid"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
func TestReadFile(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Read invalid file", func(tPtr *testing.T) {
		if _, soteErr = ReadFile(parentCtx, TESTINVALIDFILEPATH); soteErr.ErrCode != sError.ErrItemNotFound {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrItemNotFound, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("Read valid file", func(tPtr *testing.T) {
		if _, soteErr = ReadFile(parentCtx, testLocalFilepath); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
func TestWriteFile(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		contents          []byte
		destFilepath      string
	)

	tPtr.Cleanup(func() {
		RemoveFile(parentCtx, destFilepath)
	})

	// tPtr.Run("Write invalid file", func(tPtr *testing.T) {
	// 	if _, soteErr = WriteFile(parentCtx, TESTINVALIDFILEPATH, []byte("INVALID TEST")); soteErr.ErrCode != sError.ErrItemNotFound {
	// 		tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrItemNotFound, soteErr.FmtErrMsg)
	// 	}
	// })

	tPtr.Run("Write valid file", func(tPtr *testing.T) {
		// Read contents to be written
		if contents, soteErr = ReadFile(parentCtx, testLocalFilepath); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
		// 	Write the file
		targetFilepath := strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, "create-test-file.jpeg"}, "/")
		if destFilepath, soteErr = WriteFile(parentCtx, targetFilepath, contents); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
func TestGetMIMEType(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		contents          = make([]byte, 0)
	)

	tPtr.Run("Get File MIME Type Using a Valid File", func(tPtr *testing.T) {
		if contents, soteErr = ReadFile(parentCtx, strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME},
			"/")); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if contentType := GetMIMEType(contents); contentType != "image/jpeg" {
			tPtr.Errorf("%v Failed: Expected Content-Type to be %v but got %v", testName, "image/jpeg", contentType)
		}
	})
}
func TestDocumentUpload(tPtr *testing.T) {
	var (
		function, _, _, _       = runtime.Caller(0)
		testName                = runtime.FuncForPC(function).Name()
		soteErr                 sError.SoteError
		err                     error
		s3ClientServerPtr       *S3ClientServer
		formFiles               map[string][]*multipart.FileHeader
		fileFieldNameOne        = ParentDocumentKey.String()
		testSingleFilenameOne   = "single-upload.jpeg"
		testSingleFilenameTwo   = "1"
		testMultipleFilenameOne = "multiple-upload-one.jpeg"
		testMultipleFilenameTwo = "multiple-upload-two.jpeg"
		uploadResponseOne       *UploadResponse
		uploadResponseTwo       *UploadResponse
		uploadsResponse         map[string][]*UploadResponse
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

		if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
			AppEnvironment:    TESTAPPENVIRONMENT,
			TestMode:          testMode,
		}); soteErr.ErrCode == nil {
			tPtr.Run("Using Valid files", func(tPtr *testing.T) {
				for _, file := range formFiles[fileFieldNameOne] {
					if openedFile, err = file.Open(); err != nil {
						break
					} else {
						if uploadResponseOne, soteErr = s3ClientServerPtr.SingleDocumentUpload(parentCtx, openedFile,
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
						if uploadResponseTwo, soteErr = s3ClientServerPtr.SingleDocumentUpload(parentCtx, openedFile,
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
		formFiles[SupportingDocumentsKey.String()] = tMultipleDocForm.File[SupportingDocumentsKey.String()]

		if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			AppEnvironment:    TESTAPPENVIRONMENT,
			ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
			FormFiles:         formFiles,
			TestMode:          testMode,
		}); soteErr.ErrCode == nil {
			if uploadsResponse, _, soteErr = s3ClientServerPtr.MultipleDocumentsUpload(parentCtx); soteErr.ErrCode != nil {
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
		s3ClientServerPtr *S3ClientServer
		contents          = make([]byte, 0)
		filenameOne       = "single-direct-image-upload.jpeg"
		filenameTwo       = "single-direct-file-upload.pdf"
		keysOne           = new(ObjectKeys)
		keysTwo           = new(ObjectKeys)
	)

	tPtr.Cleanup(func() {
		if keysOne != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keysOne.ProcessedObjectKey)
		}

		if keysTwo != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keysTwo.ProcessedObjectKey)
		}
	})

	if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
		AppConfigName:     sConfigParams.DOCUMENTS,
		ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
		AppEnvironment:    TESTAPPENVIRONMENT,
		TestMode:          testMode,
	}); soteErr.ErrCode == nil {
		tPtr.Run("Using a Valid Image", func(tPtr *testing.T) {
			if contents, soteErr = ReadFile(parentCtx, strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME},
				"/")); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			keysOne = GetObjectKeys(filenameOne, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			if soteErr = s3ClientServerPtr.DocumentUpload(parentCtx, keysOne.ProcessedObjectKey, contents,
				GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		})

		tPtr.Run("Using a Valid File", func(tPtr *testing.T) {
			if contents, soteErr = ReadFile(parentCtx, strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALPDFFILENAME},
				"/")); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			keysTwo = GetObjectKeys(filenameTwo, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			if soteErr = s3ClientServerPtr.DocumentUpload(parentCtx, keysTwo.ProcessedObjectKey, contents,
				GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		})
	}
}
func TestDocumentDelete(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		s3ClientServerPtr *S3ClientServer
		keys              *ObjectKeys
		filename          = "test-delete-document.jpeg"
	)

	if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
		AppConfigName:     sConfigParams.DOCUMENTS,
		ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
		AppEnvironment:    TESTAPPENVIRONMENT,
		TestMode:          testMode,
	}); soteErr.ErrCode == nil {
		tPtr.Run("Using Invalid Bucket Name", func(tPtr *testing.T) {
			keys = GetObjectKeys(filename, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			s3ClientServerPtr.BucketName = ""
			if soteErr = s3ClientServerPtr.DocumentDelete(parentCtx, keys.InboundObjectKey); soteErr.ErrCode != sError.ErrMissingParameters {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrMissingParameters, soteErr.FmtErrMsg)
			}
		})
	}
}
func TestDocumentPreSignedURL(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		s3ClientServerPtr *S3ClientServer
		keys              *ObjectKeys
		documentLinks     *DocumentLinks
		contents          []byte
		sourceFilepath    string
		filename          = "test-presigned-url.jpeg"
	)

	tPtr.Cleanup(func() {
		if keys != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keys.ProcessedObjectKey)
			s3ClientServerPtr.DocumentDelete(parentCtx, keys.InboundObjectKey)
		}
	})

	if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
		AppConfigName:     sConfigParams.DOCUMENTS,
		AppEnvironment:    TESTAPPENVIRONMENT,
		ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
		TestMode:          testMode,
	}); soteErr.ErrCode == nil {
		tPtr.Run("Generate presigned URL using valid document-link", func(tPtr *testing.T) {
			sourceFilepath = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME}, "/")
			keys = GetObjectKeys(filename, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			// Read contents of file to be uploaded
			if contents, soteErr = ReadFile(parentCtx, sourceFilepath); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
			// Upload source document
			if s3ClientServerPtr.DocumentUpload(parentCtx, keys.ProcessedObjectKey, contents,
				GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			if documentLinks, soteErr = GetDocumentLinks(parentCtx, s3ClientServerPtr.BucketName, keys); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}

			s3ClientServerPtr.DocumentParamsPtr.DocumentsLink = documentLinks.ProcessedDocumentLink
			if _, soteErr = s3ClientServerPtr.DocumentPreSignedURL(parentCtx, 6); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		})

		tPtr.Run("Check Invalid Pre-Signed Document URL", func(tPtr *testing.T) {
			if _, soteErr = ValidatePreSignedDocumentURL(TESTINVALIDFILEPATH); soteErr.ErrCode != sError.ErrItemNotFound {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrItemNotFound, soteErr.FmtErrMsg)
			}
		})
	}
}
func TestDocumentCopy(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		s3ClientServerPtr *S3ClientServer
		keys              *ObjectKeys
		sourceFilepath    string
		contents          []byte
		filename          = "test-copy-document.jpeg"
	)

	tPtr.Cleanup(func() {
		if keys != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keys.InboundObjectKey)
			s3ClientServerPtr.DocumentDelete(parentCtx, keys.ProcessedObjectKey)
		}
	})

	tPtr.Run("Copy document", func(tPtr *testing.T) {
		if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
			AppEnvironment:    TESTAPPENVIRONMENT,
			TestMode:          testMode,
		}); soteErr.ErrCode == nil {
			sourceFilepath = strings.Join([]string{GetFullDirectoryPath(), TESTFILESFOLDER, TESTLOCALFILENAME}, "/")
			keys = GetObjectKeys(filename, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
			// Read contents of file to be copied
			if contents, soteErr = ReadFile(parentCtx, sourceFilepath); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
			// Upload source document
			if s3ClientServerPtr.DocumentUpload(parentCtx, keys.InboundObjectKey, contents, GetMIMEType(contents)); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
			// Copy uploaded document
			if soteErr = s3ClientServerPtr.DocumentCopy(parentCtx, keys.InboundObjectKey, keys.ProcessedObjectKey); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
			}
		}
	})
}
func TestEmbedMetadata(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		s3ClientServerPtr *S3ClientServer
		keys              *ObjectKeys
		metadata          = map[string]interface{}{
			"client-company":   TESTCLIENTCOMPANYNAMEONE,
			"company-supplier": TESTCOMPANYSUPPLIERNAMEONE,
			"invoice-number":   "0223BE726355",
			"convert-test":     map[string]string{"convert-interface": "test"},
		}
		keysOne     *ObjectKeys
		keysTwo     *ObjectKeys
		filenameOne = "test-get-embed.jpeg"
		filenameTwo = "test-invalid-get-embed.jpeg"
		metadataTwo = map[string]interface{}{
			"inbound-object-key": fmt.Sprintf(strings.Join([]string{filepath.Dir(TESTOBJECTKEYONE), filenameOne}, "/")),
		}
	)

	tPtr.Cleanup(func() {
		s3ClientServerPtr.BucketName = TESTS3BUCKETNAME

		if keys != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keys.InboundObjectKey)
		}

		if keysOne != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keysOne.InboundObjectKey)
		}

		if keysTwo != nil {
			s3ClientServerPtr.DocumentDelete(parentCtx, keysTwo.InboundObjectKey)
		}
	})

	tPtr.Run("Embed metadata", func(tPtr *testing.T) {
		if s3ClientServerPtr, soteErr = NewS3ClientServer(parentCtx, &DocumentParams{
			AppConfigName:     sConfigParams.DOCUMENTS,
			ClientCompanyName: TESTCLIENTCOMPANYNAMEONE,
			AppEnvironment:    TESTAPPENVIRONMENT,
			TestMode:          testMode,
		}); soteErr.ErrCode == nil {
			tPtr.Run("Using invalid document-link", func(tPtr *testing.T) {
				tObjectKeys := GetObjectKeys(TESTINVALIDFILEPATH, fmt.Sprint(s3ClientServerPtr.DocumentParamsPtr.ClientCompanyName))
				if soteErr = s3ClientServerPtr.EmbedMetadata(parentCtx, tObjectKeys.InboundObjectKey,
					metadata); soteErr.ErrCode != sError.ErrItemNotFound {
					tPtr.Errorf("%v Failed: Expected error code to be %v or %v but got %v", testName, sError.ErrBusinessServiceError,
						sError.ErrItemNotFound,
						soteErr.FmtErrMsg)
				}
			})

			tPtr.Run("Using valid document-link", func(tPtr *testing.T) {
				if _, _, keys, soteErr = copyTestDocument(tPtr, "test-embed.jpeg", false, false); soteErr.ErrCode == nil {
					fmt.Println("HERE")
					if soteErr = s3ClientServerPtr.EmbedMetadata(parentCtx, keys.InboundObjectKey,
						metadata); soteErr.ErrCode != nil {
						tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
					}
				}
			})

			tPtr.Run("Get Embedded Metadata Using Valid Bucket Name", func(tPtr *testing.T) {
				if _, _, keysOne, soteErr = copyTestDocument(tPtr, filenameOne, false, false); soteErr.ErrCode == nil {
					if soteErr = s3ClientServerPtr.EmbedMetadata(parentCtx, keysOne.InboundObjectKey, metadataTwo); soteErr.ErrCode == nil {
						if _, soteErr = s3ClientServerPtr.GetEmbeddedDocumentMetadata(parentCtx, keysOne); soteErr.ErrCode != nil {
							tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
						}
					}
				}
			})

			tPtr.Run("Get Embedded Metadata Using Invalid Bucket Name", func(tPtr *testing.T) {
				if _, _, keysTwo, soteErr = copyTestDocument(tPtr, filenameTwo, false, false); soteErr.ErrCode == nil {
					if soteErr = s3ClientServerPtr.EmbedMetadata(parentCtx, keysTwo.InboundObjectKey, metadataTwo); soteErr.ErrCode == nil {
						s3ClientServerPtr.BucketName = ""
						if _, soteErr = s3ClientServerPtr.GetEmbeddedDocumentMetadata(parentCtx,
							keysTwo); soteErr.ErrCode != sError.ErrMissingParameters {
							tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrMissingParameters,
								soteErr.FmtErrMsg)
						}
					}
				}
			})
		}
	})
}
