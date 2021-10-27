package controllers

import (
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type DefaultController struct {
	Ctx     iris.Context
	Service services.MenuService
	Session *sessions.Session
}

//加载框架
func (this *DefaultController) Get() mvc.Result {
	uid := this.Session.GetInt64Default(common.USER_ID, 0)
	if uid == 0 {
		this.Ctx.Redirect("/login")
	}
	return &mvc.View{
		Name: "default/index.html",
		Data: iris.Map{"Title": "SpiderHub 管理系统"},
	}
}

func (this *DefaultController) GetHeader() {
	user := this.Session.Get(common.USER)
	menus := this.Service.GetLevelMenu(0)
	this.Ctx.ViewData("menus", menus)
	this.Ctx.ViewData("username", user.(*system.SystemAdmin).Username)
	this.Ctx.ViewData("id", user.(*system.SystemAdmin).Id)
	this.Ctx.View("default/head.html")
}

func (this *DefaultController) GetMain() mvc.Result {
	return &mvc.View{
		Name: "default/main.html",
	}
}

func (this *DefaultController) GetDialog() mvc.Result {
	uri := this.Ctx.URLParamDefault("url", "")
	isRead, _ := this.Ctx.URLParamBool("isReadonly")

	return &mvc.View{
		Name:   "default/dialog.html",
		Data:   iris.Map{"url": uri, "is_read": isRead},
		Layout: "layout/dialog.html",
	}
}
