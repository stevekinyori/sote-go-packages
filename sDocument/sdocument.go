/*
	This is a wrapper for Sote Document storage which uses AWS S3
*/
package sDocument

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)


// Downloads an item from an S3 Bucket in the region configured in the shared config
// or AWS_REGION environment variable.
//
// Usage:
//    go run s3_download_object.go BUCKET ITEM
func New() {
	bucket := "sote-test-bucket-scott-yacko"
	item := "/inbound/927d056b-6f90-4048-a462-5ecb5ae2d5a5.PDF"

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create("/tmp" + item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", item, err)
	}

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}