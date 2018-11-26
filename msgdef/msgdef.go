package msgdef

import "encoding/json"

type MsgWithError struct {
	Ret   int             `json:"ret"`
	Error string          `json:"error"`
	Msg   json.RawMessage `json:"msg"`
}

func (mwe *MsgWithError) Load(msg interface{}) {
	mbytes, _ := json.Marshal(msg)
	mwe.Msg = json.RawMessage(mbytes)
}
