package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handleRequest(event events.S3Event) {
	for _, record := range event.Records {
		// Retrieve the S3 bucket and key from the event
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		// Create an S3 session
		sess, err := session.NewSession()
		if err != nil {
			log.Fatalf("Failed to create S3 session: %v", err)
		}

		// Create an S3 service client
		svc := s3.New(sess)

		// Get the contents of the zipped file from S3
		input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
		result, err := svc.GetObject(input)
		if err != nil {
			log.Fatalf("Failed to get object from S3: %v", err)
		}

		// Read the contents of the zipped file into a byte slice
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, result.Body); err != nil {
			log.Fatalf("Failed to read object from S3: %v", err)
		}
		zipData := buf.Bytes()
		fmt.Println("Read", len(zipData), "bytes from S3") // Worked up till here

		//Out the files from zipData to S3
		zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
		if err != nil {
			log.Fatal(err)
		}
		keyTrimmed := strings.TrimSuffix(key, ".zip") // remove .zip from the key
		// Read all the files from zip archive
		for _, zipFile := range zipReader.File {
			// get file minetype
			contentType := mime.TypeByExtension(filepath.Ext(zipFile.Name))
			unzippedFileBytes, err := readZipFile(zipFile)
			if err != nil {
				log.Println(err)
				continue
			}

			//_ = unzippedFileBytes // this is unzipped file bytes
			// Create an S3 object input
			input := &s3.PutObjectInput{
				Bucket:      aws.String(bucket),
				Key:         aws.String(fmt.Sprintf("%v/%v", keyTrimmed, zipFile.Name)),
				Body:        bytes.NewReader(unzippedFileBytes),
				ContentType: aws.String(contentType),
			}

			// Upload the []byte data as an S3 object
			_, err = svc.PutObject(input)
			if err != nil {
				fmt.Println("Failed to upload object to S3:", err)
				return
			}
		}

	}
}

func main() {
	lambda.Start(handleRequest)
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
