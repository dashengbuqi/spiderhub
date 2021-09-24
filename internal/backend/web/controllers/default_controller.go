package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type DefaultController struct {
	Ctx iris.Context
}

//加载框架
func (this *DefaultController) Get() mvc.Result {
	return &mvc.View{
		Name: "default/index.html",
		Data: iris.Map{"Title": "SpiderHub 管理系统"},
	}
}
