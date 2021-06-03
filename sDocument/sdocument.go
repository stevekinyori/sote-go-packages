/*
	This is a wrapper for Sote Golang developers to access services from AWS S3 Bucket.
*/
package sDocument

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const (
	LOGMESSAGEPREFIX = "sDocument"
)

type DocumentManager struct {
	sInboundS3BucketURL    string
	sProcessedS3BucketURL  string
	tesseractServerManager *TesseractServerManager

	sync.Mutex
}

var (
	sSession *session.Session
	tErr     error
)

/*
This will establish a session using the default .aws location
*/
func init() {
	sLogger.DebugMethod()

	sSession, tErr = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if tErr != nil {
		log.Fatalln(tErr)
	}
}

/*
	New will create a Sote Document Manager and a connection to AWS S3 Bucket URL. The connection is established to
	processed/unprocessed URL.
*/
func New(application, environment string, testMode bool) (SDocumentManagerPtr *DocumentManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tS3BucketURL string
	)

	// Initialize the values for Document Manager
	SDocumentManagerPtr = &DocumentManager{sInboundS3BucketURL: "", sProcessedS3BucketURL: ""}

	// Get AWS S3 Bucket URL for inbound files
	if tS3BucketURL, soteErr = sConfigParams.SGetS3BucketURL(application, environment,
		sConfigParams.UNPROCESSEDDOCUMENTSKEY); soteErr.ErrCode == nil {
		SDocumentManagerPtr.sInboundS3BucketURL = tS3BucketURL
	}

	// Get AWS S3 Bucket URL for processed files
	if tS3BucketURL, soteErr = sConfigParams.SGetS3BucketURL(application, environment, sConfigParams.PROCESSEDDOCUMENTSKEY); soteErr.ErrCode == nil {
		SDocumentManagerPtr.sProcessedS3BucketURL = tS3BucketURL
	}

	return
}

/*
	SConvertImageFormat writes out the image in the same/different format

	EXAMPLE:
		if dsm, soteErr := New(); soteErr.ErrCode == nil {
			if _, soteErr = dsm.ConvertImageFormat("logo.jpg","out.png"); soteErr.ErrCode != nil {
				sLogger.Info(fmt.Sprintf("ConvertImageFormat failed:Expected soteErr to be nil:%v ", soteErr.FmtErrMsg))
			}
		}
*/
func SConvertImageFormat(sourcePath string, targetPath string) (pdfFilePtr *imagick.ImageCommandResult,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var sErr error
	imagick.Initialize()
	defer imagick.Terminate() // Memory leak cleanup

	if pdfFilePtr, sErr = imagick.ConvertImageCommand([]string{"convert", sourcePath, targetPath}); sErr != nil {
		fmt.Println(sErr, pdfFilePtr, GetFullDirectoryPath())
		if strings.Contains(sErr.Error(), "No such file or directory") {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{"upload path or filename"}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/* GetFullDirectoryPath will get current working directory */
func GetFullDirectoryPath() (path string) {
	sLogger.DebugMethod()

	path, _ = os.Getwd()

	return
}

/* UploadDocument will upload a document to AWS S3 Bucket. */
func (dmPtr *DocumentManager) UploadDocument(processedFilePath string) (sPutObjectOutput *s3.PutObjectOutput, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		sErr                   error
		sProcessedDocument     *os.File
		sProcessedDocumentInfo os.FileInfo
		sFileSize              int64
		sFileBuffer            []byte
		sFileName              string
	)

	if sProcessedDocument, sErr = os.Open(processedFilePath); sErr == nil {
		// Get processed file information
		if sProcessedDocumentInfo, sErr = sProcessedDocument.Stat(); sErr == nil {
			sFileSize = sProcessedDocumentInfo.Size()
			sFileName = sProcessedDocumentInfo.Name()
			sFileBuffer = make([]byte, sFileSize)
			_, _ = sProcessedDocument.Read(sFileBuffer)
		}

		// Add object to Amazon S3 Bucket
		sPutObjectOutput, sErr = s3.New(sSession).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(strings.Split(dmPtr.sProcessedS3BucketURL, "/")[0]),
			Key:                  aws.String(strings.Join([]string{"processed", sFileName}, "/")),
			ACL:                  aws.String("bucket-owner-full-control"),
			Body:                 bytes.NewReader(sFileBuffer),
			ContentLength:        aws.Int64(sFileSize),
			ContentType:          aws.String(http.DetectContentType(sFileBuffer)),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
		})
	}

	if sErr != nil {
		sLogger.Info(sErr.Error())
		panic(fmt.Sprintf("%v", sErr.Error()))
	}

	return
}

/* DownloadDocument wil get a document from AWS S3 Bucket. */
func (dmPtr *DocumentManager) DownloadDocument(objectKey, localInboundPath string, testMode bool) (sFileBuffer *bytes.Buffer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		sErr         error
		sRawObject   *s3.GetObjectOutput
		sInboundFile *os.File
	)

	params := make(map[string]string)
	params["objectKey"] = objectKey
	params["testMode"] = strconv.FormatBool(testMode)

	// Retrieve object from Amazon S3 Bucket using key
	if sRawObject, sErr = s3.New(sSession).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(strings.Split(dmPtr.sInboundS3BucketURL, "/")[0]),
		Key:    aws.String(objectKey),
	}); sErr == nil {
		sFileBuffer = new(bytes.Buffer)
		if _, sErr = sFileBuffer.ReadFrom(sRawObject.Body); sErr == nil {
			// Create	file
			if sInboundFile, sErr = os.Create(localInboundPath + strings.Split(objectKey,
				"/")[1]); sErr == nil {
				_, sErr = sFileBuffer.WriteTo(sInboundFile)
			}
		}
	}

	if sErr != nil {
		sLogger.Info(fmt.Sprintf("%v", sErr.Error()))
		if strings.Contains(sErr.Error(), "NoSuchKey: The specified key does not exist") {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{fmt.Sprintf("key (%v)", objectKey)}), sError.EmptyMap)
		} else {
			panic(fmt.Sprintf("%v", sErr.Error()))
		}
	}

	return
}

// DeleteDocument  will delete a document from AWS S3 Bucket. Arguments objectKey and testMode are required .
func (dmPtr *DocumentManager) DeleteDocument(objectKey string, testMode bool) (sDeleteObjectOutput *s3.DeleteObjectOutput, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		sErr error
	)

	params := make(map[string]string)
	params["objectKey"] = objectKey
	params["testMode"] = strconv.FormatBool(testMode)

	// Delete object from Amazon S3 Bucket using key
	if sDeleteObjectOutput, sErr = s3.New(sSession).DeleteObject(&s3.DeleteObjectInput{
		Bucket:    aws.String(strings.Split(dmPtr.sInboundS3BucketURL, "/")[0]),
		Key:       aws.String(objectKey),
		VersionId: nil,
	}); sErr != nil {
		soteErr = dmPtr.s3BucketErrorHandle(sErr, params)
	}

	return
}

func (dmPtr *DocumentManager) s3BucketErrorHandle(err error, params map[string]string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		panicError = true
		testMode   = false
	)

	if strings.Contains(err.Error(), "NoSuchKey: The specified key does not exist.") {
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{params["objectKey"]}), sError.EmptyMap)
		panicError = false
	}

	if testMode, err = strconv.ParseBool(params["testMode"]); err != nil {
		testMode = false
	}

	sLogger.Info(soteErr.FmtErrMsg)
	if panicError && !testMode {
		panic(soteErr.FmtErrMsg)
	}

	return
}
