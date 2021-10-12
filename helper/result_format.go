package helper

import (
	"encoding/json"
)

type ResultEasyUItem struct {
	Pages   *Pagination
	Models  interface{}
	_result map[string]interface{}
}

func (this *ResultEasyUItem) process() {
	this._result = map[string]interface{}{
		"rows":  this.Models,
		"total": this.Pages.GetTotal(),
	}
}

func (this *ResultEasyUItem) ToJson() string {
	this.process()
	bt, _ := json.Marshal(this._result)
	return string(bt)
}

func (this *ResultEasyUItem) ToMap() map[string]interface{} {
	this.process()
	return this._result
}

type resultFormat struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func ResultError(msg string) string {
	o_o := &resultFormat{
		Status: 0,
		Msg:    msg,
		Data:   nil,
	}
	bt, _ := json.Marshal(o_o)
	return string(bt)
}

func ResultSuccess(msg string, data interface{}) string {
	o_o := resultFormat{
		Status: 1,
		Msg:    msg,
		Data:   data,
	}
	bt, _ := json.Marshal(o_o)
	return string(bt)
}
