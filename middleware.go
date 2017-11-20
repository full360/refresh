package refresh

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
}

func NewLoggingMiddleware(logger log.Logger) *loggingMiddleware {
	return &loggingMiddleware{logger}
}

func (l *loggingMiddleware) LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			l.logger.Log(
				"took", time.Since(begin).Seconds(),
			)
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}
