package helper

import (
	"encoding/json"
)

type ResultEUINode struct {
	Rows  interface{} `json:"rows"`
	Total int64       `json:"total"`
}

type ResultEasyUItem struct {
	Pages   *Pagination
	Models  interface{}
	_result *ResultEUINode
}

func (this *ResultEasyUItem) process() {
	this._result = &ResultEUINode{
		Rows:  this.Models,
		Total: this.Pages.GetTotal(),
	}
}

func (this *ResultEasyUItem) ToJson() string {
	this.process()
	if this._result.Total == 0 {
		this._result.Rows = []string{}
	}
	bt, _ := json.Marshal(this._result)
	return string(bt)
}

func (this *ResultEasyUItem) ToMap() *ResultEUINode {
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
