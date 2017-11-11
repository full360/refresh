package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.full360.com/full360/refresh"
)

func main() {
	s3Bucket := flag.String("s3-bucket", "dev.playground.prometheus", "Name of the AWS S3 Bucket")
	awsRegion := flag.String("aws-region", "us-east-1", "AWS Region")
	downloadDir := flag.String("download-dir", "something", "Download directory")

	// log setup
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*awsRegion),
	}))

	downloader := &refresh.Downloader{
		Svc:     s3.New(sess),
		Manager: s3manager.NewDownloader(sess),
		Logger:  logger,
		Config: struct {
			Bucket string
			Prefix string
		}{
			Bucket: *s3Bucket,
			Prefix: "",
		},
	}

	downloader.Download(*downloadDir)
}
