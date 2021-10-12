package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/dashengbuqi/spiderhub/internal/backend/widgets"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type MenuController struct {
	Ctx     iris.Context
	Service services.MenuService
}

//加载列表视图
func (this *MenuController) GetList() mvc.Result {
	return &mvc.View{
		Name: "menu/list.html",
	}
}

//返回列表数据
func (this *MenuController) PostList() string {
	page := this.Ctx.PostValueIntDefault("page", 1)
	pageSize := this.Ctx.PostValueIntDefault("rows", 15)
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.PostMenuList(map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
		"sort":     sort,
		"order":    order,
	})
	return result
}

func (this *MenuController) PostEdit() string {
	id, _ := this.Ctx.URLParamInt64("id")
	form := this.Ctx.FormValues()
	err := this.Service.ModifyMenuItem(id, form)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("操作成功", nil)
}

func (this *MenuController) GetEdit() {
	id, _ := this.Ctx.URLParamInt64("id")
	model := this.Service.GetRowBy(id)
	this.Ctx.ViewData("id", model.Id)
	this.Ctx.ViewData("task_name", model.TaskName)
	this.Ctx.ViewData("full_name", model.FullName)
	this.Ctx.ViewData("path", model.Path)
	this.Ctx.ViewData("icon", model.Icon)
	this.Ctx.ViewData("sort", model.Sort)
	this.Ctx.ViewData("type", model.Type)
	this.Ctx.ViewData("parent_id", model.ParentId)
	this.Ctx.ViewData("types", widgets.Combobox{
		Id:       "type",
		Name:     "type",
		Value:    model.Type,
		Data:     model.GetTypeComboList(),
		Multiple: false,
		OnChange: "DefaultOnChange",
	})
	m := system.NewMenu()
	this.Ctx.ViewData("parents", widgets.Combobox{
		Id:       "parent_id",
		Name:     "parent_id",
		Value:    model.ParentId,
		Editable: false,
		Data:     m.GetMenuTreeList(),
		Width:    200,
	})
	this.Ctx.View("menu/edit.html")
}

func (this *MenuController) Delete() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.RemoveMenu(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("操作成功", nil)
}
