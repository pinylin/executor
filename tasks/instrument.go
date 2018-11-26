package tasks

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

//MetricsMiddleware metrics service middleware
func MetricsMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return metricsMiddleware{
			next:           next,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}

type metricsMiddleware struct {
	next           Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (mw metricsMiddleware) WxCheck(timestamp, nonce, signature string) bool {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetAllCmts"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.WxCheck(timestamp, nonce, signature)
}
func (mw metricsMiddleware) WxEvent(xmlText string) bool {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetComments"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.WxEvent(xmlText)
}
