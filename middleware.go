package refresh

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func NewInstrumentingMiddleware(counter metrics.Counter, histogram metrics.Histogram) *instrumentingMiddleware {
	return &instrumentingMiddleware{counter, histogram}
}

func (i *instrumentingMiddleware) InstrumentingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			i.requestCount.With("path", r.URL.Path, "method", r.Method).Add(1)
			i.requestLatency.With("path", r.URL.Path, "method", r.Method).Observe(time.Since(begin).Seconds())
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}
