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
	"context"
	"encoding/json"
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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sCustom"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

const (
	LOGMESSAGEPREFIX = "sDocument"
)

// DocumentParams Holds params used for extracting metadata from a document
type DocumentParams struct {
	DocumentsLink        string
	ClientCompanyId      int
	ClientCompanyName    string
	MountPointEnvVarName string
	AppConfigName        string
	AppEnvironment       string
	TestMode             bool
	FormFiles            map[string][]*multipart.FileHeader
}

// S3ClientServer  Holds params used by AWS S3 Service
type S3ClientServer struct {
	DocumentParamsPtr  *DocumentParams
	BucketName         string
	S3BucketMountPoint string
	S3ClientPtr        *s3.Client
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

	// Get document mount point
	documentsMountPoint, soteErr = GetDocumentsMountPoint(ctx, documentParamsPtr.MountPointEnvVarName)
	if soteErr.ErrCode != nil {
		return
	}

	s3ClientServerPtr = new(S3ClientServer)
	s3ClientServerPtr.DocumentParamsPtr = documentParamsPtr
	s3ClientServerPtr.S3BucketMountPoint = documentsMountPoint
	cfg, err = config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		sLogger.Info(err.Error())
		soteErr = sError.GetSError(210399, nil, sError.EmptyMap)
		return
	}

	// Initialize S3 Client
	s3ClientServerPtr.S3ClientPtr = s3.NewFromConfig(cfg)
	// Get S3 Bucket Name
	bucketName, soteErr = sConfigParams.GetAWSS3Bucket(ctx, documentParamsPtr.AppConfigName)
	if soteErr.ErrCode != nil {
		return
	}
	s3ClientServerPtr.BucketName = bucketName

	return
}

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
	)
	uploadResponse = make(map[string][]*UploadResponse)
	wg.Add(2)

	go func() {
		wg.Wait()
		close(parentDocChan)
		close(supportingDocsChan)
		close(supportingDocsAttachmentsChan)
		close(parentDocAttachmentsChan)
	}()

	go func() {
		var (
			tUploadRes *UploadResponse
			uploadsRes []*UploadResponse
		)

		sLogger.Info("starting parent-document upload...")
		for _, file := range formFiles[ParentDocumentKey.String()] {
			if openedFile, err := file.Open(); err != nil {
				break
			} else {
				if tUploadRes, soteErr = s3sPtr.SingleDocumentUpload(ctx, openedFile, file.Filename); soteErr.ErrCode == nil {
					uploadsRes = append(uploadsRes, tUploadRes)
					parentDocChan <- uploadsRes
					parentDocAttachmentsChan <- []string{uploadsRes[0].ProcessedDocumentLink}
				} else {
					parentDocChan <- []*UploadResponse{}
					parentDocAttachmentsChan <- []string{}
				}
			}
		}

		wg.Done()
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

				if openedFile, err := file.Open(); err != nil {
					break
				} else {
					if tUploadRes, soteErr = s3sPtr.SingleDocumentUpload(ctx, openedFile, file.Filename); soteErr.ErrCode == nil {
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

		wg.Done()
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

	return
}

func (s3sPtr *S3ClientServer) SingleDocumentUpload(ctx context.Context, file multipart.File, filename string) (uploadResponse *UploadResponse,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		err      error
		contents []byte

		inboundFilepath   string
		processedFilepath string
		documentLinks     = new(DocumentLinks)
		tDocumentLinks    = new(DocumentLinks)
		metadata          = make(map[string]interface{})
		tObjectKeys       = new(ObjectKeys)
		wg                sync.WaitGroup
	)
	wg.Add(3)
	soteErrChan := make(chan map[string]sError.SoteError)
	docLinkChan := make(chan *DocumentLinks)

	go func() {
		wg.Wait()
		close(soteErrChan)
		close(docLinkChan)
	}()

	if contents, err = ioutil.ReadAll(file); err == nil {
		tObjectKeys = GetObjectKeys(filename, fmt.Sprint(s3sPtr.DocumentParamsPtr.ClientCompanyId))
		inboundFilepath, processedFilepath, _ = s3sPtr.GetMountPointFilepath(tObjectKeys)
		go func() {
			if _, soteErr = WriteFile(ctx, inboundFilepath, contents); soteErr.ErrCode == nil {
				soteErrChan <- map[string]sError.SoteError{}
			} else {
				soteErrChan <- map[string]sError.SoteError{INBOUNDLINKKEY: soteErr}
			}

			wg.Done()
		}()

		go func() {
			if _, soteErr = WriteFile(ctx, processedFilepath, contents); soteErr.ErrCode == nil {
				// Add inbound object key to metadata to enable inbound document to be deleted
				metadata[InboundObjectKeyFieldName.String()] = tObjectKeys.InboundObjectKey
				if soteErr = s3sPtr.EmbedMetadata(ctx, tObjectKeys.ProcessedObjectKey, metadata); soteErr.ErrCode != nil {
					soteErrChan <- map[string]sError.SoteError{PROCESSEDLINKKEY: soteErr}
				}
				soteErrChan <- map[string]sError.SoteError{}
			} else {
				soteErrChan <- map[string]sError.SoteError{PROCESSEDLINKKEY: soteErr}
			}

			wg.Done()
		}()

		go func() {
			if tDocumentLinks, soteErr = GetDocumentLinks(ctx, s3sPtr.BucketName, tObjectKeys); soteErr.ErrCode == nil {
				docLinkChan <- tDocumentLinks
			} else {
				docLinkChan <- &DocumentLinks{}
				soteErrChan <- map[string]sError.SoteError{"document-links": soteErr}
			}

			wg.Done()
		}()

		documentLinks = <-docLinkChan
		soteError := <-soteErrChan

		for _, v := range soteError {
			if v.ErrCode != nil {
				soteErr = v
				break
			}
		}

		if soteErr.ErrCode == nil {
			uploadResponse = &UploadResponse{
				FileName:              filename,
				DocumentLinks:         documentLinks,
				ObjKeys:               tObjectKeys,
				ProcessedDocumentLink: documentLinks.ProcessedDocumentLink,
			}

			sLogger.Info(fmt.Sprintf("successfully uploaded document %v", filename))
		}
	}

	if err != nil {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return
}

/*GetMountPointFilepath Will return display name,source document filepath and target document filepath from the mount point*/
func (s3sPtr *S3ClientServer) GetMountPointFilepath(objectKeys *ObjectKeys) (sourceFilepath, targetFilepath, displayName string) {
	sLogger.DebugMethod()

	displayName = GetDocumentName(objectKeys.InboundObjectKey)
	sourceFilepath = fmt.Sprintf("%v", strings.Join([]string{s3sPtr.S3BucketMountPoint, objectKeys.MountPointInboundObjectKey}, "/"))
	targetFilepath = fmt.Sprintf("%v", strings.Join([]string{s3sPtr.S3BucketMountPoint, PROCESSEDFOLDER, appEnvironment,
		fmt.Sprint(s3sPtr.DocumentParamsPtr.ClientCompanyId), displayName}, "/"))

	return
}

/*
DocumentCopy Copies a document to specified path in AWS S3 Bucket.
	Input Parameters:
		* sourceFilepath - location of document to be copied
		* targetFilepath - location where document will be copied
*/
func (s3sPtr *S3ClientServer) DocumentCopy(ctx context.Context, sourceFilepath, targetFilepath string) (sFilepath string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		contents  []byte
		tFilepath string
	)

	if contents, soteErr = ReadFile(ctx, sourceFilepath); soteErr.ErrCode == nil {
		if soteErr = CreateSubdirectories(targetFilepath); soteErr.ErrCode == nil {
			_, _ = os.Create(targetFilepath) // Ignore this error
			if tFilepath, soteErr = WriteFile(ctx, targetFilepath, contents); soteErr.ErrCode == nil {
				sFilepath = tFilepath
			}
		}
	}

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
					soteErr = sError.GetSError(210599, nil, sError.EmptyMap)
				}
			} else {
				sLogger.Info(err.Error())
				soteErr = sError.GetSError(207110, sError.BuildParams([]string{"objHeaders"}), sError.EmptyMap)
			}
		} else {
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(207105, []interface{}{"sourceObject", sourceObject}, sError.EmptyMap)
		}
	} else {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

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

	keys = GetObjectKeys(GetDocumentName(s3sPtr.DocumentParamsPtr.DocumentsLink), fmt.Sprint(s3sPtr.DocumentParamsPtr.ClientCompanyId))
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
		found         bool
		tFilepath     = strings.Join([]string{s3sPtr.S3BucketMountPoint, keys.MountPointInboundObjectKey}, "/")
	)

	if found, soteErr = ValidateFilepath(tFilepath); found && soteErr.ErrCode == nil {
		if headObjOutput, err = s3sPtr.S3ClientPtr.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(s3sPtr.BucketName),
			Key:    aws.String(keys.InboundObjectKey),
		}); err == nil {
			metadata = headObjOutput.Metadata
		} else {
			soteErr = AmazonTextractErrorHandler(ctx, err)
		}
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
					soteErr = sError.GetSError(207110, sError.BuildParams([]string{"expires-in"}), sError.EmptyMap)
				}
			} else {
				soteErr = sError.GetSError(207110, sError.BuildParams([]string{"created-date string"}), sError.EmptyMap)
			}
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{"document"}), sError.EmptyMap)
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
func GetObjectKeys(filename, clientCompanyId string) (objectKeys *ObjectKeys) {
	sLogger.DebugMethod()

	objectKeys = &ObjectKeys{
		MountPointInboundObjectKey: strings.Join([]string{INBOUNDFOLDER, appEnvironment, clientCompanyId, filename}, "/"),
		InboundObjectKey:           strings.Join([]string{BUCKETDOCSUBFOLDER, INBOUNDFOLDER, appEnvironment, clientCompanyId, filename}, "/"),
		ProcessedObjectKey:         strings.Join([]string{BUCKETDOCSUBFOLDER, PROCESSEDFOLDER, appEnvironment, clientCompanyId, filename}, "/"),
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
		soteErr = sError.GetSError(210599, nil, sError.EmptyMap)
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