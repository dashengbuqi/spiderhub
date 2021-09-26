package controllers

import (
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
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
	menus := this.Service.GetLevelMenu(0)
	this.Ctx.ViewData("menus", menus)
	this.Ctx.View("default/head.html")
}

func (this *DefaultController) GetMain() mvc.Result {
	return &mvc.View{
		Name: "default/main.html",
	}
}
