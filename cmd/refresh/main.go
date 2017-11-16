package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"gitlab.full360.com/full360/refresh/prom"
	"gitlab.full360.com/full360/refresh/storage"
)

func main() {
	addr := flag.String("address", "127.0.0.1", "Listen address")
	port := flag.Int("port", 3000, "Listen port")
	promUrl := flag.String("prom-url", "", "Prometheus URL")
	s3Bucket := flag.String("s3-bucket", "", "Name of the AWS S3 Bucket")
	s3BucketPrefix := flag.String("s3-bucket-prefix", "", "Name of the AWS S3 Bucket Prefix")
	awsRegion := flag.String("aws-region", "us-east-1", "AWS Region")
	downloadDir := flag.String("download-dir", "", "Download directory")

	// log setup
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// http client setup
	httpClient := cleanhttp.DefaultClient()

	// aws session setup
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*awsRegion),
	}))

	// storage setup
	storage := storage.NewS3Storage(
		s3.New(sess),
		s3manager.NewDownloader(sess),
		log.With(logger, "storage", "s3"),
		*s3Bucket,
		*s3BucketPrefix,
		*downloadDir,
	)

	// prometheus service setup
	var ps prom.Service
	ps = prom.NewService(
		storage,
		httpClient,
		*promUrl,
		"POST",
	)
	ps = prom.NewLoggingService(log.With(logger, "service", "prom"), ps)

	r := mux.NewRouter()
	r.Handle(
		"/prom/refresh",
		prom.RefreshHandler(ps),
	).Methods("POST")

	srv := http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", "localhost", 3000),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Log("msg", "HTTP server", "addr", *addr, "port", *port)
	logger.Log("err", srv.ListenAndServe())
}
