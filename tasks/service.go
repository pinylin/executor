package tasks

import (
	"pinylin.top/executor/tasks/tools"
)

/*Service wechat服务
 */
type Service interface {
	WxCheck(timestamp, nonce, signature string) bool
	WxEvent(xmlText string) bool
}

//Middleware for service
type Middleware func(Service) Service

//NewFixedService is a simple implementation of service
func NewFixedService() Service {
	return &fixedService{}
}

//fixedService for Users
type fixedService struct{}

func (s *fixedService) WxCheck(timestamp, nonce, signature string) bool {
	if timestamp == "" || nonce == "" || signature == "" {
		return false
	}
	signatureGen := tools.MakeSignature(timestamp, nonce)
	return signatureGen == signature
}

func (s *fixedService) WxEvent(xmlText string) bool {
	if xmlText == "" {
		return false
	}
	// TODO 事件处理
	return true
}



