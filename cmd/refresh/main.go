package main

import (
	"flag"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-kit/kit/log"
	"gitlab.full360.com/full360/refresh"
)

func main() {
	s3Bucket := flag.String("s3-bucket", "dev.playground.prometheus", "Name of the AWS S3 Bucket")
	s3BucketPrefix := flag.String("s3-bucket-prefix", "", "Name of the AWS S3 Bucket Prefix")
	awsRegion := flag.String("aws-region", "us-east-1", "AWS Region")
	downloadDir := flag.String("download-dir", "something", "Download directory")

	// log setup
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*awsRegion),
	}))

	downloader := &refresh.Downloader{
		Svc:     s3.New(sess),
		Manager: s3manager.NewDownloader(sess),
		Logger:  log.With(logger, "service", "s3Service"),
		Config: struct {
			Bucket string
			Prefix string
		}{
			Bucket: *s3Bucket,
			Prefix: *s3BucketPrefix,
		},
	}

	if err := downloader.Download(*downloadDir); err != nil {
		logger.Log("msg", err.Error())
	}

}
