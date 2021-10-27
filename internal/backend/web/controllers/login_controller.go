package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type LoginController struct {
	Ctx     iris.Context
	Service services.UserService
	Session *sessions.Session
}

func (this *LoginController) Get() mvc.Result {
	uid := this.Session.GetInt64Default(common.USER_ID, 0)
	if uid > 0 {
		this.Ctx.Redirect("/")
	}
	return &mvc.View{
		Name:   "default/login.html",
		Layout: "layout/login.html",
	}
}

func (this *LoginController) Post() string {
	name := this.Ctx.FormValue("username")
	pwd := this.Ctx.FormValue("password")

	u, err := this.Service.GetByUsernameAndPwd(name, pwd)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	this.Session.Set(common.USER_ID, u.Id)
	this.Session.Set(common.USER, u)
	return helper.ResultSuccess("成功", nil)
}

//退出
func (this *LoginController) GetDestroy() {
	this.Session.Destroy()
	this.Ctx.Redirect("/login")
}
