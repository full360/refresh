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
	"github.com/full360/refresh"
	"github.com/full360/refresh/app"
	"github.com/full360/refresh/health"
	"github.com/full360/refresh/storage"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	addr := flag.String("address", "127.0.0.1", "Listen address")
	port := flag.Int("port", 3000, "Listen port")
	appUrl := flag.String("app-url", "", "Application URL")
	appUrlMethod := flag.String("app-url-method", "POST", "Application URL Method")
	s3Bucket := flag.String("s3-bucket", "", "Name of the AWS S3 Bucket")
	s3BucketPrefix := flag.String("s3-bucket-prefix", "", "Name of the AWS S3 Bucket Prefix")
	awsRegion := flag.String("aws-region", "us-east-1", "AWS Region")
	downloadDir := flag.String("download-dir", "", "Download directory")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.Parse()

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

	// health service setup
	var hs health.Service
	hs = health.NewService()

	// storage service setup
	var ss storage.Storage
	ss = storage.NewS3Storage(
		s3.New(sess),
		s3manager.NewDownloader(sess),
		log.With(logger, "component", "storage"),
		*s3Bucket,
		*s3BucketPrefix,
		*downloadDir,
	)

	// app service setup
	var as app.Service
	as = app.NewService(
		ss,
		httpClient,
		*appUrl,
		*appUrlMethod,
	)
	as = app.NewLoggingService(log.With(logger, "component", "app"), as)

	// middleware setup
	lm := refresh.NewLoggingMiddleware(log.With(logger, "component", "server"))
	li := refresh.NewInstrumentingMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "requests_total",
			Help:      "Number of requests received",
		}, []string{"path", "method"}),
		kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "request_duration_seconds",
			Help:      "Total duration of requests in seconds",
		}, []string{"path", "method"}),
	)

	r := mux.NewRouter()
	r.Handle(
		"/health",
		li.InstrumentingHandler(health.HealthHandler(hs)),
	).Methods("GET")
	r.Handle(
		"/app/refresh",
		li.InstrumentingHandler(lm.LoggingHandler(app.RefreshHandler(as))),
	).Methods("POST")
	r.Handle(
		"/metrics",
		promhttp.Handler(),
	).Methods("GET")

	srv := http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", *addr, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Log("msg", "HTTP server", "addr", *addr, "port", *port)
	logger.Log("err", srv.ListenAndServe())
}
