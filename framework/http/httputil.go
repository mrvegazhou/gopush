package httphelper

import (
	"encoding/json"
	"gopush/const"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func NewWithError(code int, msg string) *Response {
	return &Response{
		Code:   code,
		Msg: 	msg,
		Data:   nil,
	}
}

func NewSuccess(data interface{}) *Response {
	return &Response{
		Code:    constdefine.SUCCESS,
		Msg: "success",
		Data:    data,
	}
}