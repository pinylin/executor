package tasks

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	"pinylin.top/executor/tasks/params"
)

//Endpoints collects the endpoints that comprise the Service
type Endpoints struct {
	WxCheckEndpoint    endpoint.Endpoint
	WxEventEndpoint  endpoint.Endpoint
}

//MakeEndpoints for service
func MakeEndpoints(s Service, logger log.Logger, tracer stdopentracing.Tracer) Endpoints {
	//auth := authorization.AuthorizationMiddleware(logger)
	//adminAuth := authorization.AdminAuthMiddleware(logger)
	return Endpoints{
		WxCheckEndpoint:  opentracing.TraceServer(tracer, "GET /wx")(MakeWxCheckEndpoint(s)),
		WxEventEndpoint:  opentracing.TraceServer(tracer, "POST /wx")(MakeWxEventEndpoint(s)),
		//GetAllCmtsEndpoint:  opentracing.TraceServer(tracer, "POST /getallcmts")(adminAuth(MakeGetAllCmtsEndpoint(s))),
	}
}


//
func MakeWxCheckEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		var span stdopentracing.Span
		span, ctx = stdopentracing.StartSpanFromContext(ctx, "wx check")
		span.SetTag("service", "wechat")
		defer span.Finish()
		req := request.(params.WxCheckReq)
		if s.WxCheck(req.Timestamp, req.Nonce, req.Signature) {

		}else{

		}
		return req.Echostr, nil
		//var mwe msgdef.MsgWithError
		//mwe.Load(params.WxCheckResp{Echostr: ""})
		//return mwe, err
	}
}

func MakeWxEventEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		var span stdopentracing.Span
		span, ctx = stdopentracing.StartSpanFromContext(ctx, "wx check")
		span.SetTag("service", "wechat")
		defer span.Finish()
		// TODO 事件处理
		//req := request.(params.WxCheckReq)
		//if s.WxCheck(req.Timestamp, req.Nonce, req.Signature) {
		//
		//}else{
		//
		//}
		return nil, nil
		//var mwe msgdef.MsgWithError
		//mwe.Load(params.WxCheckResp{Echostr: ""})
		//return mwe, err
	}
}
