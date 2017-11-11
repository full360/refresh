package refresh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-kit/kit/log"
)

type Downloader struct {
	Svc     *s3.S3
	Manager *s3manager.Downloader
	Logger  log.Logger
	Config  struct {
		Bucket string
		Prefix string
	}
}

func (s *Downloader) Download(downloadDir string) error {
	err := s.Svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &s.Config.Bucket,
	}, func(page *s3.ListObjectsOutput, last bool) bool {
		for _, obj := range page.Contents {
			// skip directories
			if strings.HasSuffix(*obj.Key, "/") {
				continue
			}
			// Create the directories in the path
			file := filepath.Join(downloadDir, *obj.Key)
			if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
				s.Logger.Log("msg", err.Error())
			}

			// Setup the local file
			fd, err := os.Create(file)
			if err != nil {
				s.Logger.Log("msg", err.Error())
			}

			defer fd.Close()

			s.Logger.Log("msg", fmt.Sprintf("downloading s3://%s/%s to %s", s.Config.Bucket, *obj.Key, file))
			s.Manager.Download(fd, &s3.GetObjectInput{Bucket: &s.Config.Bucket, Key: obj.Key})
		}
		return true
	})
	if err != nil {
		panic(err)
	}

	return err
}
