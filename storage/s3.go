package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-kit/kit/log"
)

type s3Storage struct {
	svc     *s3.S3
	manager *s3manager.Downloader
	logger  log.Logger
	config  struct {
		bucket string
		prefix string
		dir    string
	}
}

func NewS3Storage(s3 *s3.S3, manager *s3manager.Downloader, logger log.Logger, bucket, prefix, dir string) *s3Storage {
	return &s3Storage{
		svc:     s3,
		manager: manager,
		logger:  logger,
		config: struct {
			bucket string
			prefix string
			dir    string
		}{
			bucket: bucket,
			prefix: prefix,
			dir:    dir,
		},
	}
}

func (s *s3Storage) Download() error {
	err := s.svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &s.config.bucket,
	}, func(page *s3.ListObjectsOutput, last bool) bool {
		for _, obj := range page.Contents {
			// skip directories
			if strings.HasSuffix(*obj.Key, "/") {
				continue
			}
			// Create the directories in the path
			file := filepath.Join(s.config.dir, *obj.Key)
			if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
				s.logger.Log("err", err.Error())
			}

			// Setup the local file
			fd, err := os.Create(file)
			if err != nil {
				s.logger.Log("err", err.Error())
			}

			defer fd.Close()

			s.logger.Log("msg", fmt.Sprintf("downloading s3://%s/%s to %s", s.config.bucket, *obj.Key, file))
			s.manager.Download(fd, &s3.GetObjectInput{Bucket: &s.config.bucket, Key: obj.Key})
		}
		return true
	})
	if err != nil {
		return err
	}

	return nil
}
