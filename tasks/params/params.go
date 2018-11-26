package params

import (
	"context"
	"encoding/xml"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	//"pinylin.top/excutor/wechat/domain"
)

var (
	//PathVarNotExistError router 变量不存在
	PathVarNotExistError = "path variable not exist error"
)

type WxCheckReq struct {
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Echostr   string `json:"echostr"`
}

type WxCheckResp struct {
	Echostr string `json:"echostr"`
}

type WxEventReq struct {
	XMLName      xml.Name      `xml:"xml"`
	ToUserName   string        `xml:"ToUserName,CDATA"`
	FromUserName string        `xml:"FromUserName,CDATA,omitempty"`
	CreateTime   time.Duration `xml:"CreateTime"`
	MsgType      string        `xml:"MsgType,CDATA"`
	// 关注/取关
	Event string `xml:"Event,CDATA,omitempty"`
	// 自定义菜单事件
	EventKey string `xml:"EventKey,CDATA,omitempty"`
	// 扫描带参数二维码事件
	Ticket string `xml:"Ticket,CDATA,omitempty"`

	// 上报地理位置事件
	//Latitude string `xml:"Latitude,CDATA,omitempty"` // 纬度
	//Longitude string `xml:"Longitude,CDATA,omitempty"`	 // 经度
	//Precision string `xml:"Precision,CDATA,omitempty"`  // 精度
}

type DefaultResp struct {
	Err error `json:"err"`
}
type GetAllCmtsReq struct {
	CurrPage int `bson:"currPage" json:"currPage"`
	PageSize int `bson:"pageSize" json:"pageSize"`
}

type GetCommentsReq struct {
	MessID   string `bson:"messID" json:"messID"`
	UserName string `bson:"userName" json:"userName"`
	CurrPage int    `bson:"currPage" json:"currPage"`
	PageSize int    `bson:"pageSize" json:"pageSize"`
}

func DecodeWxCheckReq(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if vars == nil {
		return nil, errors.New(PathVarNotExistError)
	}
	request := WxCheckReq{Timestamp: vars["timestamp"], Nonce: vars["nonce"], Signature: vars["signature"]}
	return request, nil
}

func DecodeWxEventReq(_ context.Context, r *http.Request) (interface{}, error) {
	var request WxEventReq
	err := xml.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, err
}

//func DecodeGetCmtsReq(_ context.Context, r *http.Request) (interface{}, error) {
//	var request GetCommentsReq
//	err := json.NewDecoder(r.Body).Decode(&request)
//	if err != nil {
//		return nil, err
//	}
//	return request, err
//}
