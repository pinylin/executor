package wechat

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"pinylin.top/executor/coder"
	"pinylin.top/executor/wechat/params"

	//"pinylin.top/executor/wechat/authorization"
)

//MakeHTTPHandler for wechat
func MakeHTTPHandler(ctx context.Context, e Endpoints, logger log.Logger, tracer stdopentracing.Tracer) *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(coder.EncodeError),
	}
	//tokenBefore := httptransport.ServerBefore(authorization.JWTToContext)
	//tokenOptions := append(options, tokenBefore)
	r.Methods("GET").Path("/wx").Handler(httptransport.NewServer(
		e.WxCheckEndpoint,
		params.DecodeWxCheckReq,
		coder.EncodeJSONResp,
		options...,
	))
	r.Methods("POST").Path("/wx").Handler(httptransport.NewServer(
		e.WxEventEndpoint,
		params.DecodeWxEventReq,
		coder.EncodeJSONResp,
		options...,
	))
	r.Path("/metrics").Handler(promhttp.Handler())
	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

//GatewayEndpoints  endpoints for apigateway
type GatewayEndpoints struct {
	GetAllCmtsEndpoint  endpoint.Endpoint
	GetCommentsEndpoint endpoint.Endpoint
}