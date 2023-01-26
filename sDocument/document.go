package sDocument

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
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sCustom"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	LOGMESSAGEPREFIX = "sDocument"
)

// DocumentParams Holds params used for extracting metadata from a document
type DocumentParams struct {
	DocumentsLink     string
	ClientCompanyName string
	AppConfigName     string
	AppEnvironment    string
	TestMode          bool
	FormFiles         map[string][]*multipart.FileHeader
	// ClientCompanyId   int
	// MountPointEnvVarName string
	// IgnoreMountPoint     bool
}

// FileParams Holds params used in file operations i.e. write, delete, read
type FileParams struct {
	Wg                 *sync.WaitGroup
	ErrChan            chan sError.SoteError
	DocumentLinksChan  chan *DocumentLinks
	UploadResponseChan chan *UploadResponse
	Contents           []byte
	BucketName         string
	ObjectKey          string
	ObjectKeys         *ObjectKeys
	File               multipart.File
	Filename           string
	// TargetFilepath     string
}

// S3ClientServer  Holds params used by AWS S3 Service
type S3ClientServer struct {
	DocumentParamsPtr  *DocumentParams
	BucketName         string
	S3BucketMountPoint string
	S3ClientPtr        *s3.Client
	UploaderPtr        *manager.Uploader
}

type ObjectKeys struct {
	MountPointInboundObjectKey string
	InboundObjectKey           string
	ProcessedObjectKey         string
}

// DocumentLinks Holds generated processed and inbound document links
type DocumentLinks struct {
	InboundDocumentLink   string
	ProcessedDocumentLink string
}

type UploadResponse struct {
	FileName              string         `json:"file-name"`
	DocumentLinks         *DocumentLinks `json:"document-links"`
	ObjKeys               *ObjectKeys    `json:"object-keys"`
	ProcessedDocumentLink string         `json:"processed-document-link"`
}

// S3PresignGetObjectAPI defines the interface for the PreSignGetObject function.
// We use this interface to test the function using a mocked service.
type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

var (
	testMode       bool
	appEnvironment string
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

/*
NewS3ClientServer Creates a new S3 Client Instance Server
*/
func NewS3ClientServer(ctx context.Context, documentParamsPtr *DocumentParams,
	optFns ...func(*config.LoadOptions) error) (s3ClientServerPtr *S3ClientServer, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err                 error
		cfg                 aws.Config
		bucketName          string
		documentsMountPoint string
	)

	testMode = documentParamsPtr.TestMode
	appEnvironment = documentParamsPtr.AppEnvironment

	// Initialize AWS S3 Server params
	s3ClientServerPtr = new(S3ClientServer)
	s3ClientServerPtr.DocumentParamsPtr = documentParamsPtr
	s3ClientServerPtr.S3BucketMountPoint = documentsMountPoint
	if cfg, err = config.LoadDefaultConfig(ctx, optFns...); err != nil {
		sLogger.Info(err.Error())
		soteErr = sError.GetSError(sError.ErrAWSSessionError, nil, sError.EmptyMap)
		return
	}
	// Create S3 Client
	s3ClientServerPtr.S3ClientPtr = s3.NewFromConfig(cfg)
	// Create an uploader with the session and default options
	partMiBs := int64(10)
	s3ClientServerPtr.UploaderPtr = manager.NewUploader(s3ClientServerPtr.S3ClientPtr, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	// Get S3 Bucket Name
	if bucketName, soteErr = sConfigParams.GetAWSS3Bucket(ctx, documentParamsPtr.AppConfigName); soteErr.ErrCode != nil {
		return
	}
	// Initialize S3 Bucket Name
	s3ClientServerPtr.BucketName = bucketName

	return
}

// MultipleDocumentsUpload Uploads multiple documents
func (s3sPtr *S3ClientServer) MultipleDocumentsUpload(ctx context.Context) (uploadResponse map[string][]*UploadResponse,
	attachment []string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		formFiles                     = s3sPtr.DocumentParamsPtr.FormFiles
		tParentDocResponse            []*UploadResponse
		tSupportingDocsResponse       []*UploadResponse
		supportingDocsAttachments     []string
		parentDocAttachment           []string
		wg                            sync.WaitGroup
		parentDocChan                 = make(chan []*UploadResponse)
		supportingDocsChan            = make(chan []*UploadResponse)
		supportingDocsAttachmentsChan = make(chan []string)
		parentDocAttachmentsChan      = make(chan []string)
		soteErrChan                   = make(chan sError.SoteError)
		uploadResChan                 = make(chan *UploadResponse)
		parentSoteErr                 sError.SoteError

		// supportingSoteErr             sError.SoteError
	)
	uploadResponse = make(map[string][]*UploadResponse)
	wgCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Make sure cancel is called to release resources even if no errors
	defer close(parentDocChan)
	defer close(supportingDocsChan)
	defer close(supportingDocsAttachmentsChan)
	defer close(parentDocAttachmentsChan)
	go func() {
		var (
			tUploadRes *UploadResponse
			uploadsRes []*UploadResponse
		)
		wg.Add(1)

		sLogger.Info("starting parent-document upload...")
		for _, file := range formFiles[ParentDocumentKey.String()] {
			if openedFile, err := file.Open(); err != nil {
				break
			} else {
				go s3sPtr.DocumentsUpload(wgCtx, &FileParams{
					Wg:                 &wg,
					ErrChan:            soteErrChan,
					UploadResponseChan: uploadResChan,
					File:               openedFile,
					Filename:           file.Filename,
				})
				if parentSoteErr = <-soteErrChan; parentSoteErr.ErrCode == nil {
					tUploadRes = <-uploadResChan
					uploadsRes = append(uploadsRes, tUploadRes)
					parentDocChan <- uploadsRes
					parentDocAttachmentsChan <- []string{uploadsRes[0].ProcessedDocumentLink}
				} else {
					parentDocChan <- []*UploadResponse{}
					parentDocAttachmentsChan <- []string{}
				}
			}
		}
	}()

	go func() {
		var (
			tUploadRes  *UploadResponse
			uploadsRes  []*UploadResponse
			attachments []string
		)

		if files, found := formFiles[SupportingDocumentsKey.String()]; found {
			sLogger.Info("starting supporting-documents upload...")
			for _, file := range files {
				wg.Add(1)

				if openedFile, err := file.Open(); err != nil {
					break
				} else {
					go s3sPtr.DocumentsUpload(wgCtx, &FileParams{
						Wg:                 &wg,
						ErrChan:            soteErrChan,
						UploadResponseChan: uploadResChan,
						File:               openedFile,
						Filename:           file.Filename,
					})
					// todo: determine if we care about this error
					if tSoteErr := <-soteErrChan; tSoteErr.ErrCode == nil {
						tUploadRes = <-uploadResChan
						uploadsRes = append(uploadsRes, tUploadRes)
						attachments = append(attachments, tUploadRes.ProcessedDocumentLink)
					}
				}
			}
			supportingDocsChan <- uploadsRes
			supportingDocsAttachmentsChan <- attachments
		} else {
			supportingDocsChan <- []*UploadResponse{}
			supportingDocsAttachmentsChan <- []string{}
		}
	}()

	tParentDocResponse = <-parentDocChan
	tSupportingDocsResponse = <-supportingDocsChan
	supportingDocsAttachments = <-supportingDocsAttachmentsChan
	parentDocAttachment = <-parentDocAttachmentsChan
	attachment = append(parentDocAttachment, supportingDocsAttachments...)
	// Populate Upload Response
	uploadResponse[ParentDocumentKey.String()] = tParentDocResponse
	if len(tSupportingDocsResponse) > 0 {
		uploadResponse[SupportingDocumentsKey.String()] = tSupportingDocsResponse
	}

	if parentSoteErr.ErrCode != nil {
		soteErr = parentSoteErr
	}

	wg.Wait()

	return
}

/*DocumentsUpload  Will upload a document to S3 Bucket without use of a mount-point*/
func (s3sPtr *S3ClientServer) DocumentsUpload(ctx context.Context, params *FileParams) {
	sLogger.DebugMethod()

	defer params.Wg.Done()
	wgCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Make sure cancel is called to release resources even if no errors

	var (
		soteErr        sError.SoteError
		uploadResponse = new(UploadResponse)
	)

	checkContextCancel(wgCtx)
	if uploadResponse, soteErr = s3sPtr.SingleDocumentUpload(wgCtx, params.File, params.Filename); soteErr.ErrCode != nil {
		params.ErrChan <- soteErr
		params.UploadResponseChan <- &UploadResponse{}
		cancel()

		return
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		sLogger.Info(fmt.Sprintf("%v document couldn't be uploaded to S3 bucket", params.Filename))
		params.ErrChan <- sError.GetSError(sError.ErrContextCancelled, sError.BuildParams([]string{"document upload"}), sError.EmptyMap)
		params.UploadResponseChan <- &UploadResponse{}
		cancel()

		return
	}

	params.ErrChan <- sError.SoteError{}
	params.UploadResponseChan <- uploadResponse

	return
}

func (s3sPtr *S3ClientServer) SingleDocumentUpload(ctx context.Context, file multipart.File, filename string) (uploadResponse *UploadResponse,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		err            error
		contents       []byte
		tDocumentLinks = new(DocumentLinks)
		documentLinks  = new(DocumentLinks)
		metadata       = make(map[string]interface{})
		objectKeys     = new(ObjectKeys)
		wg             sync.WaitGroup
		errChan        = make(chan sError.SoteError, 2)
		docLinkChan    = make(chan *DocumentLinks, 1)
	)
	wg.Add(3)
	wgCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Make sure cancel is called to release resources even if no errors
	defer close(errChan)
	defer close(docLinkChan)

	checkContextCancel(wgCtx)
	objectKeys = GetObjectKeys(filename, fmt.Sprint(s3sPtr.DocumentParamsPtr.ClientCompanyName))
	if contents, err = ioutil.ReadAll(file); err != nil {
		soteErr = AmazonTextractErrorHandler(ctx, err)
		return
	}
	// Upload document to inbound folder
	go s3sPtr.UploadDocument(wgCtx, &FileParams{
		Wg:        &wg,
		ErrChan:   errChan,
		Contents:  contents,
		ObjectKey: objectKeys.InboundObjectKey,
	})
	if soteErr = <-errChan; soteErr.ErrCode != nil {
		cancel()
		return
	}
	// Upload document to processed folder
	go s3sPtr.UploadDocument(wgCtx, &FileParams{
		Wg:        &wg,
		ErrChan:   errChan,
		Contents:  contents,
		ObjectKey: objectKeys.ProcessedObjectKey,
	})
	if soteErr = <-errChan; soteErr.ErrCode != nil {
		cancel()
		return
	}
	// Get the document links
	go GetDocumentsLinks(wgCtx, &FileParams{
		Wg:                &wg,
		ErrChan:           errChan,
		DocumentLinksChan: docLinkChan,
		BucketName:        s3sPtr.BucketName,
		ObjectKeys:        objectKeys,
	})
	if soteErr = <-errChan; soteErr.ErrCode != nil {
		cancel()
		return
	}

	// Add inbound object key to metadata to enable inbound document to be deleted
	metadata[InboundObjectKeyFieldName.String()] = objectKeys.InboundObjectKey
	soteErr = s3sPtr.EmbedMetadata(ctx, objectKeys.ProcessedObjectKey, metadata)
	if soteErr.ErrCode != nil {
		cancel()
		return
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		sLogger.Info(fmt.Sprintf("%v document couldn't be uploaded to S3 bucket", filename))
		errChan <- sError.GetSError(sError.ErrContextCancelled, sError.BuildParams([]string{"document upload"}), sError.EmptyMap)
		cancel()

		return
	}

	if tDocumentLinks = <-docLinkChan; tDocumentLinks != nil {
		documentLinks = tDocumentLinks
	}

	if soteErr.ErrCode == nil {
		uploadResponse = &UploadResponse{
			FileName:              filename,
			DocumentLinks:         documentLinks,
			ObjKeys:               objectKeys,
			ProcessedDocumentLink: documentLinks.ProcessedDocumentLink,
		}

		sLogger.Info(fmt.Sprintf("successfully uploaded document %v", filename))
	}

	wg.Wait()

	return
}

// UploadDocument Will upload document to AWS S3 Bucket directly
func (s3sPtr *S3ClientServer) UploadDocument(ctx context.Context, params *FileParams) {

	defer params.Wg.Done()
	wgCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Make sure cancel is called to release resources even if no errors
	sLogger.Info(fmt.Sprintf("started uploading document with object key: %v", params.ObjectKey))

	var (
		soteErr sError.SoteError
	)

	checkContextCancel(wgCtx)
	if soteErr = s3sPtr.DocumentUpload(ctx, params.ObjectKey, params.Contents, GetMIMEType(params.Contents)); soteErr.ErrCode != nil {
		params.ErrChan <- soteErr
		cancel()

		return
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		sLogger.Info(fmt.Sprintf("Document couldn't be uploaded to S3"))
		params.ErrChan <- sError.GetSError(sError.ErrContextCancelled, sError.BuildParams([]string{"document upload"}), sError.EmptyMap)
		cancel()

		return
	}

	params.ErrChan <- sError.SoteError{}

	return
}

func GetDocumentsLinks(ctx context.Context, params *FileParams) {
	sLogger.DebugMethod()

	defer params.Wg.Done()
	wgCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Make sure cancel is called to release resources even if no errors
	sLogger.Info(fmt.Sprintf("started generating document links for %v", filepath.Base(params.ObjectKeys.InboundObjectKey)))

	var (
		documentLinks = new(DocumentLinks)
		soteErr       sError.SoteError
	)

	checkContextCancel(wgCtx)
	if documentLinks, soteErr = GetDocumentLinks(wgCtx, params.BucketName, params.ObjectKeys); soteErr.ErrCode != nil {
		params.ErrChan <- soteErr
		params.DocumentLinksChan <- &DocumentLinks{}
		cancel()
		return
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		sLogger.Info(fmt.Sprintf("document links couldn't be created"))
		params.ErrChan <- sError.GetSError(sError.ErrContextCancelled, sError.BuildParams([]string{"get document link"}), sError.EmptyMap)
		params.DocumentLinksChan <- &DocumentLinks{}
		cancel()
		return
	}

	params.ErrChan <- sError.SoteError{}
	params.DocumentLinksChan <- documentLinks
	sLogger.Info(fmt.Sprintf("finished generating document links for %v", filepath.Base(params.ObjectKeys.InboundObjectKey)))

	return
}

/*
DocumentCopy Creates a copy of an object that is already stored in Amazon S3
	Input Parameters:
		* sourceObjectKey - Specifies the source object for the copy operation
		* targetObjectKey - The key of the destination object
*/
func (s3sPtr *S3ClientServer) DocumentCopy(ctx context.Context, sourceObjectKey, targetObjectKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		copyObjectOutput *s3.CopyObjectOutput
		err              error
		metadata         = make(map[string]interface{}, 0)
	)

	if copyObjectOutput, err = s3sPtr.S3ClientPtr.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s3sPtr.BucketName),
		CopySource: aws.String(fmt.Sprintf("%v/%v", s3sPtr.BucketName, sourceObjectKey)),
		Key:        aws.String(targetObjectKey),
	}); err != nil {
		soteErr = AmazonTextractErrorHandler(ctx, err)
		return
	}

	// Add inbound object key to metadata to enable inbound document to be deleted
	metadata[InboundObjectKeyFieldName.String()] = sourceObjectKey
	soteErr = s3sPtr.EmbedMetadata(ctx, targetObjectKey, metadata)

	sLogger.Info(fmt.Sprintf("successfully copied document with object key %v version-id %v to document with object key %v version-id %v",
		sourceObjectKey, *copyObjectOutput.CopySourceVersionId, targetObjectKey, *copyObjectOutput.VersionId))

	return
}

// EmbedMetadata Will embed document information a document
func (s3sPtr *S3ClientServer) EmbedMetadata(ctx context.Context, processedObjectKey string,
	metadata map[string]interface{}) (soteErr sError.SoteError) {
	var (
		inInterface  map[string]interface{}
		err          error
		sourceObject *s3.HeadObjectOutput
		newMetadata  = make(map[string]string)
		objHeaders   []byte
		value        string
		buffer       []byte
	)

	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s3sPtr.BucketName),
		Key:    aws.String(processedObjectKey),
	}

	if sourceObject, err = s3sPtr.S3ClientPtr.HeadObject(ctx, headInput); err == nil {
		// A map of metadata to store with the object in S3
		if sourceObject.Metadata != nil {
			newMetadata = sourceObject.Metadata
		}

		for k, v := range metadata {
			if reflect.ValueOf(v).Kind() == reflect.String {
				value = fmt.Sprintf("%v", v)
			} else {
				if buffer, soteErr = ConvertInterfaceToByteSlice(v); soteErr.ErrCode == nil {
					value = fmt.Sprintf("%v", strings.Trim(string(buffer), "\n"))
				}
			}
			newMetadata[k] = fmt.Sprintf("%v", value)
		}

		if objHeaders, err = json.Marshal(sourceObject); err == nil {

			if err = json.Unmarshal(objHeaders, &inInterface); err == nil {
				params := &s3.CopyObjectInput{
					Bucket:            aws.String(s3sPtr.BucketName),
					CopySource:        aws.String(fmt.Sprintf("%v/%v", s3sPtr.BucketName, processedObjectKey)),
					Key:               aws.String(processedObjectKey),
					MetadataDirective: "REPLACE",
					Metadata:          newMetadata,
				}

				for field, val := range inInterface {
					if val != nil {
						switch field {
						// Not a complete list of available headers to copy
						case "CacheControl":
							params.CacheControl = aws.String(fmt.Sprintf("%v", val))
						case "ContentDisposition":
							params.ContentDisposition = aws.String(fmt.Sprintf("%v", val))
						case "ContentEncoding":
							params.ContentEncoding = aws.String(fmt.Sprintf("%v", val))
						case "ContentLanguage":
							params.ContentLanguage = aws.String(fmt.Sprintf("%v", val))
						case "ContentType":
							params.ContentType = aws.String(fmt.Sprintf("%v", val))
						case "ServerSideEncryption":
							params.ServerSideEncryption = s3Types.ServerSideEncryption(fmt.Sprintf("%v", val))
						default:
							break
						}
					}
				}

				if _, err = s3sPtr.S3ClientPtr.CopyObject(ctx, params); err != nil {
					sLogger.Info(err.Error())
					soteErr = sError.GetSError(sError.ErrBusinessServiceError, nil, sError.EmptyMap)
				}
			} else {
				sLogger.Info(err.Error())
				soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{"objHeaders"}), sError.EmptyMap)
			}
		} else {
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrJSONConversionError, []interface{}{"sourceObject", sourceObject}, sError.EmptyMap)
		}
	} else {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return

}

/*DocumentUpload Uses an upload manager to upload data to an object in a bucket.
The upload manager breaks large data into parts and uploads the parts concurrently.
  Params:
    objectKey    - Name of object in AWS S3 Bucket
    contents     - Data to be uploaded
    contentType  - Object media type
*/
func (s3sPtr *S3ClientServer) DocumentUpload(ctx context.Context, objectKey string, contents []byte,
	contentType string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	largeBuffer := bytes.NewReader(contents)
	if _, err := s3sPtr.UploaderPtr.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s3sPtr.BucketName),
		Key:         aws.String(objectKey),
		Body:        largeBuffer,
		ContentType: aws.String(contentType),
	}); err != nil {
		sLogger.Info(fmt.Sprintf("couldn't upload large object to %v:%v. An error occured: %v\n",
			s3sPtr.BucketName, objectKey, err.Error()))
		soteErr = AmazonTextractErrorHandler(ctx, err)
		return
	}

	sLogger.Info("successfully uploaded document with object key " + objectKey)

	return
}

// DocumentDelete  will delete a document from AWS S3 Bucket.
func (s3sPtr *S3ClientServer) DocumentDelete(ctx context.Context, objectKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		deleteObjectOutput *s3.DeleteObjectOutput
		err                error
	)

	// Delete object from Amazon S3 Bucket using key
	if deleteObjectOutput, err = s3sPtr.S3ClientPtr.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3sPtr.BucketName),
		Key:    aws.String(objectKey),
	}); err == nil {
		if deleteObjectOutput.VersionId != nil {
			sLogger.Info(fmt.Sprintf("Deleted object with key %v version %v", objectKey, *deleteObjectOutput.VersionId))
		}
	} else {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return
}

// DocumentPreSignedURL Will return a pre-signed document URL
func (s3sPtr *S3ClientServer) DocumentPreSignedURL(ctx context.Context, expiryDuration time.Duration) (preSignedDocumentURL string,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		keys *ObjectKeys
	)

	keys = GetObjectKeys(GetDocumentName(s3sPtr.DocumentParamsPtr.DocumentsLink), fmt.Sprint(s3sPtr.DocumentParamsPtr.ClientCompanyName))
	preSignedDocumentURL, soteErr = s3sPtr.getDocumentPreSignedURL(ctx, keys.ProcessedObjectKey, expiryDuration)

	return
}

// getDocumentPreSignedURL will return a pre-signed document URL
func (s3sPtr *S3ClientServer) getDocumentPreSignedURL(ctx context.Context, objectKey string,
	expiryDuration time.Duration) (preSignedDocumentURL string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		presignedClient *s3.PresignClient
		err             error
		preSignedReq    *v4.PresignedHTTPRequest
	)

	presignedClient = s3.NewPresignClient(s3sPtr.S3ClientPtr, s3.WithPresignExpires(expiryDuration*time.Second))

	if preSignedReq, err = GetPresignedURL(ctx, presignedClient, &s3.GetObjectInput{
		Bucket: aws.String(s3sPtr.BucketName),
		Key:    aws.String(objectKey),
	}); err == nil {
		preSignedDocumentURL = preSignedReq.URL
	} else {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return
}

// GetEmbeddedDocumentMetadata Will return embedded document metadata
func (s3sPtr *S3ClientServer) GetEmbeddedDocumentMetadata(ctx context.Context, keys *ObjectKeys) (metadata map[string]string,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		headObjOutput *s3.HeadObjectOutput
		err           error
	)

	if headObjOutput, err = s3sPtr.S3ClientPtr.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s3sPtr.BucketName),
		Key:    aws.String(keys.InboundObjectKey),
	}); err == nil {
		metadata = headObjOutput.Metadata
	} else {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return
}

// GetPresignedURL retrieves a presigned URL for an Amazon S3 bucket object.
// Inputs:
//     ctx is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, the pre-signed URL for the object and nil.
//     Otherwise, nil and an error from the call to PreSignGetObject.
func GetPresignedURL(ctx context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(ctx, input)
}

// ValidatePreSignedDocumentURL Will return a pre-signed document URL
func ValidatePreSignedDocumentURL(preSignedDocURL string) (isPreSignedDocumentURLExpired bool,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		req             *http.Request
		err             error
		urlVals         url.Values
		createdDateStr  string
		expiresIn       string
		createdDatetime time.Time
		expiryDuration  time.Duration
		currentDatetime = time.Now().UTC()
		expiryDatetime  time.Time
	)

	if req, err = http.NewRequest("GET", preSignedDocURL, nil); err == nil {
		if urlVals = req.URL.Query(); len(urlVals) > 0 {
			createdDateStr = urlVals.Get("X-Amz-Date")
			expiresIn = urlVals.Get("X-Amz-Expires")
			// 	Get document creation date in Unix Date format
			if createdDatetime, err = time.Parse(strings.ReplaceAll(strings.ReplaceAll(time.RFC3339, "-", ""), ":", ""), createdDateStr); err == nil {
				// 	Parse expiry str to time.Duration format
				if expiryDuration, err = time.ParseDuration(expiresIn + "s"); err == nil {
					// Get expiry date
					expiryDatetime = createdDatetime.Add(expiryDuration)
					// Check if pre-signed document URL is expired
					isPreSignedDocumentURLExpired = currentDatetime.After(expiryDatetime)
				} else {
					soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{"expires-in"}), sError.EmptyMap)
				}
			} else {
				soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{"created-date string"}), sError.EmptyMap)
			}
		} else {
			soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{"document"}), sError.EmptyMap)
		}
	}

	return
}

/*
GetDocumentLinks Will return a document's S3 Bucket URLs
	Input Parameters:
		* bucketName - The name of the bucket containing the object
		* objectKeys - Object keys of input document
*/
func GetDocumentLinks(ctx context.Context, bucketName string, objectKeys *ObjectKeys) (documentLinks *DocumentLinks, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	runtime.GOMAXPROCS(2)
	var (
		region        string
		wg            sync.WaitGroup
		s3URIProtocol = "s3"
	)
	documentLinks = new(DocumentLinks)

	if region, soteErr = sConfigParams.GetRegion(ctx); soteErr.ErrCode == nil {
		wg.Add(2)
		pChan := make(chan string)
		iChan := make(chan string)

		go func() {
			wg.Wait()
			close(pChan)
			close(iChan)
		}()

		go func() {
			iChan <- fmt.Sprintf("%v://%v.s3-%v.amazonaws.com/%v", s3URIProtocol, bucketName, region,
				objectKeys.InboundObjectKey)

			wg.Done()
		}()

		go func() {
			pChan <- fmt.Sprintf("%v://%v.s3-%v.amazonaws.com/%v", s3URIProtocol, bucketName, region,
				objectKeys.ProcessedObjectKey)

			wg.Done()
		}()

		documentLinks.InboundDocumentLink = <-iChan
		documentLinks.ProcessedDocumentLink = <-pChan
	}

	return
}

/*GetObjectKeys Will return object keys associated with a document
  Params:
      filename        - Name of file
      clientCompanyId - This is the system id assigned to the organization client company
*/
func GetObjectKeys(filename, clientCompanyName string) (objectKeys *ObjectKeys) {
	sLogger.DebugMethod()

	objectKeys = &ObjectKeys{
		MountPointInboundObjectKey: strings.Join([]string{INBOUNDFOLDER, appEnvironment, clientCompanyName, filename}, "/"),
		InboundObjectKey:           strings.Join([]string{BUCKETDOCSUBFOLDER, INBOUNDFOLDER, appEnvironment, clientCompanyName, filename}, "/"),
		ProcessedObjectKey:         strings.Join([]string{BUCKETDOCSUBFOLDER, PROCESSEDFOLDER, appEnvironment, clientCompanyName, filename}, "/"),
	}

	return
}

/*CreateSubdirectories Will create subdirectories in a document's path
  Params:
     * documentFilepath - Path containing subdirectories to be created
*/
func CreateSubdirectories(documentFilepath string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if err = os.MkdirAll(filepath.Dir(documentFilepath), os.ModePerm); err != nil {
		soteErr = sError.GetSError(sError.ErrBusinessServiceError, nil, sError.EmptyMap)
	}

	return
}

/*
ConvertInterfaceToByteSlice Converts interface to slice of bytes
*/
func ConvertInterfaceToByteSlice(value interface{}) (buffer []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tBuffer []byte
	)

	if tBuffer, soteErr = sCustom.JSONMarshal(value); soteErr.ErrCode == nil {
		buffer = tBuffer
	}

	return
}

// GetDocumentName will return name of document
func GetDocumentName(objectKey string) (fileName string) {
	sLogger.DebugMethod()

	fileName = filepath.Base(objectKey)

	return
}

/*
GetMIMEType Will return a file's MimeType
	Input Parameters:
		* Content - File data
*/
func GetMIMEType(content []byte) (mimeType string) {
	sLogger.DebugMethod()

	mimeType = http.DetectContentType(content)

	return
}

/*
ReadFile Will return contents of file
	Input Parameters:
		* filepath - Path of file to be read
*/
func ReadFile(ctx context.Context, filepath string) (byteSlice []byte, soteErr sError.SoteError) {
	var (
		err        error
		tByteSlice []byte
	)

	if tByteSlice, err = ioutil.ReadFile(filepath); err == nil {
		byteSlice = tByteSlice
	} else {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return
}

/*
WriteFile Will create a file with specified contents
	Input Parameters:
		* targetFilepath - Path of file to be created
		* contents       - Details be written
*/
func WriteFile(ctx context.Context, targetFilepath string, contents []byte) (destinationFilepath string, soteErr sError.SoteError) {
	var (
		err error
	)

	if soteErr = CreateSubdirectories(targetFilepath); soteErr.ErrCode == nil {
		_, _ = os.Create(targetFilepath) // Ignore this error
		if err = ioutil.WriteFile(targetFilepath, contents, fs.ModePerm); err == nil {
			destinationFilepath = targetFilepath
		} else {
			soteErr = AmazonTextractErrorHandler(ctx, err)
		}
	}

	return
}

// checkContextCancel checks if an error has happened during the context then terminates all processes using the context(
// only if context is started with a cancel signal)
func checkContextCancel(ctx context.Context) {
	select {
	case <-ctx.Done():
		return // Error somewhere? then terminate
	default: // avoid blocking
	}
}
