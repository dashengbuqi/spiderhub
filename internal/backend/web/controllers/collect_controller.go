package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/dashengbuqi/spiderhub/internal/backend/widgets"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type CollectController struct {
	Ctx     iris.Context
	Service services.CollectService
}

//加载列表视图
func (this *CollectController) GetList() mvc.Result {
	return &mvc.View{
		Name: "collect/list.html",
	}
}

func (this *CollectController) PostList() string {
	page, _ := this.Ctx.PostValueInt("page")
	pageSize, _ := this.Ctx.PostValueInt("rows")
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.GetCollectList(&helper.RequestParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	})
	return result
}

func (this *CollectController) GetEdit() {
	id, _ := this.Ctx.URLParamInt64("id")
	model := this.Service.GetRowBy(id)
	this.Ctx.ViewData("id", model.Id)
	this.Ctx.ViewData("title", model.Title)
	this.Ctx.ViewData("schedule", model.Schedule)
	this.Ctx.ViewData("storages", widgets.Combobox{
		Id:       "storage",
		Name:     "storage",
		Value:    model.Storage,
		Data:     model.GetStorageComboList(),
		Multiple: false,
	})
	this.Ctx.ViewData("methods", widgets.Combobox{
		Id:       "method",
		Name:     "method",
		Value:    model.Method,
		Data:     model.GetMethodComboList(),
		Multiple: false,
	},
	)
	this.Ctx.View("collect/edit.html")
}

func (this *CollectController) PostEdit() string {
	id, _ := this.Ctx.URLParamInt64("id")
	form := this.Ctx.FormValues()
	err := this.Service.ModifyCollectItem(id, form)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("操作成功", nil)
}

func (this *CollectController) GetRemove() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.Remove(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("操作成功", nil)
}

func (this *CollectController) GetDebug() mvc.Result {
	id, _ := this.Ctx.URLParamInt64("id")
	model := this.Service.GetRowBy(id)
	return &mvc.View{
		Name: "collect/debug.html",
		Data: iris.Map{"id": id, "status": model.Status, "content": model.CrawlerContent},
	}
}

func (this *CollectController) PostSave() string {
	id, _ := this.Ctx.URLParamInt64("id")
	code := this.Ctx.FormValue("code")
	err := this.Service.ModifyCrawler(id, code)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("保存成功", nil)
}
