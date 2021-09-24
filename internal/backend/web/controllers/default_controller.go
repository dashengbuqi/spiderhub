package controllers

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type DefaultController struct {
	Ctx     iris.Context
	Service services.MenuService
}

//加载框架
func (this *DefaultController) Get() mvc.Result {
	return &mvc.View{
		Name: "default/index.html",
		Data: iris.Map{"Title": "SpiderHub 管理系统"},
	}
}

func (this *DefaultController) GetHeader() {
	this.Ctx.ViewData("menus", string(jStr))
	this.Ctx.View("default/head.html")
}

func (this *DefaultController) GetMain() {
	this.Ctx.View("default/main.html")
}
