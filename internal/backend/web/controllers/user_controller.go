package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	Ctx     iris.Context
	Service services.UserService
}

//加载列表视图
func (this *UserController) GetList() mvc.Result {
	return &mvc.View{
		Name: "user/list.html",
	}
}

func (this *UserController) PostList() string {
	page, _ := this.Ctx.PostValueInt("page")
	pageSize, _ := this.Ctx.PostValueInt("rows")
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.GetUserList(&helper.RequestParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	})
	return result
}

func (this *UserController) GetEdit() {
	id, _ := this.Ctx.URLParamInt64("id")
	model := this.Service.GetRowBy(id)
	this.Ctx.ViewData("username", model.Username)
	this.Ctx.ViewData("email", model.Email)
	this.Ctx.ViewData("mobile", model.Mobile)
	this.Ctx.ViewData("id", model.Id)
	this.Ctx.View("user/edit.html")
}

func (this *UserController) PostEdit() string {
	id, _ := this.Ctx.URLParamInt64("id")
	form := this.Ctx.FormValues()
	err := this.Service.ModifyMenuItem(id, form)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("操作成功", nil)
}

func (this *UserController) GetRemove() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.RemoveUser(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("操作成功", nil)
}
