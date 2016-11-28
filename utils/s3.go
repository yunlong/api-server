package utils

import (

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"fmt"
	"os"
	"bytes"
	"net/http"
)

func UploadToS3(user ,filename, filepath string) (string, error) {
	aws_access_key_id := "AKIAIP6YQLDQ426QQ3DA"
	aws_secret_access_key := "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"
	token := ""

	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err := creds.Get()
	if err != nil {
    	fmt.Println(err)
	}
	config := &aws.Config{Region: aws.String("ap-southeast-1"), Credentials: creds}

	svc := s3.New(session.New(config))
	bucket := "chanhlvbucket"

	bucketName := bucket
	//fileToUpload := pathConfig + "tmp_files" + osSplit + hfilename

	file, err := os.Open(filepath)

	if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	var fileInS3 string = user + "/" + filename
	linkReturn := "https://s3-ap-southeast-1.amazonaws.com/" + bucket + "/" + fileInS3

	params := &s3.PutObjectInput{
    	Bucket:        aws.String(bucketName), // required
    	Key:           aws.String(fileInS3),       // required
    	ACL:           aws.String("public-read"),
    	Body:          fileBytes,
    	ContentLength: aws.Int64(size),
    	ContentType:   aws.String(fileType),
    	Metadata: map[string]*string{
    		"Key": aws.String("MetadataValue"), //required
    	},
	}
	_, err = svc.PutObject(params)
	if err != nil {
    	log.Println(err)
		return "", err
	}

	return linkReturn, nil
}