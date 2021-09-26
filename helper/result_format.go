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
