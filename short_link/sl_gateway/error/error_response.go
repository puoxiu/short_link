package gateway_error

import (
	"encoding/json"
	"net/http"
)

// Data 结构体用于标准化返回的 JSON 数据
type BaseResponse struct {
	Code int    `json:"code"` 
	Data any    `json:"data"` 
	Msg  string `json:"msg"`  // 错误或提示信息
}

func FailResponse(msg string, res http.ResponseWriter) {
	resp := BaseResponse{Code: 7, Msg: msg}
	byteData, _ := json.Marshal(resp)
	res.Write(byteData)
}