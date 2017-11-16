package prom

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, p Service) Service {
	return &loggingService{logger, p}
}

func (l *loggingService) Refresh() (client *http.Response, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "refresh",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Refresh()
}
