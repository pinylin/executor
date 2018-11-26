package coder

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	httptransprot "github.com/go-kit/kit/transport/http"
	"pinylin.top/executor/msgdef"
)

type errorer interface {
	error() error
}

//EncodeError encode httptransport error
func EncodeError(context context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	serviceName, _ := context.Value("ServiceName").(string)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ret":   http.StatusInternalServerError,
		"error": err.Error(),
		"from":  serviceName,
	})
}

//EncodeJSONResp  json encode to w.write()
func EncodeJSONResp(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func EncodeJSONRespWithMWE(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var mwe msgdef.MsgWithError
	mwe.Load(response)
	return json.NewEncoder(w).Encode(&mwe)
}

//DecodeJSONRequest 不关心数据， 仅仅转发
func DecodeJSONRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request interface{}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

//EncodeJSONRequest  json ecode to req.body
func EncodeJSONRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

//DecodeJSONResponse  不关心具体数据 仅仅转发
func DecodeJSONResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func DecodeJSONRequestToObj(objFactory func() interface{}) httptransprot.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		request := objFactory()
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			return nil, err
		}
		return request, nil
	}
}

// func DecodePathRequest(cxt context.Context, req *http.Request)
