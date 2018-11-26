package wechat

import (
	"time"

	"github.com/go-kit/kit/log"
)

//LoggingMiddleware for
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) WxCheck(timestamp, nonce, signature string) bool {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "WxCheck",
			"timestamp", timestamp,
			"nonce", nonce,
			"signature", signature,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.WxCheck(timestamp, nonce, signature)
}

func (mw loggingMiddleware) WxEvent(xmlText string) bool {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "WxCallBack",
			"xmlText", xmlText,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.WxEvent(xmlText)
}

