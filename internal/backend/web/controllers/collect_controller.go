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
	step := this.Ctx.URLParamDefault("step", "")
	return &mvc.View{
		Name: "collect/list.html",
		Data: iris.Map{"step": step},
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
	})
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

//调试开始
func (this *CollectController) PostBegin() string {
	id, _ := this.Ctx.URLParamInt64("id")
	code := this.Ctx.FormValue("code")
	debug_id, err := this.Service.CrawlerBegin(id, code)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("开始执行", iris.Map{"debug_id": debug_id})
}

//心跳检测
func (this *CollectController) GetHeart() string {
	id, _ := this.Ctx.URLParamInt64("id")
	debug_id, _ := this.Ctx.URLParamInt64("debug_id")
	res := this.Service.CrawlerHeart(id, debug_id, 0)
	return helper.ResultSuccess("SUCCESS", res)
}

//调试结束
func (this *CollectController) PutEnd() string {
	id, _ := this.Ctx.URLParamInt64("id")
	debug_id, _ := this.Ctx.URLParamInt64("debug_id")
	err := this.Service.CrawlerEnd(id, debug_id, 0)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("调试正在终止...", nil)
}

//正式开始
func (this *CollectController) PutStart() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.CrawlerStart(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("启动成功", nil)
}

func (this *CollectController) GetStatus() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.CrawlerStatus(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("执行完成", nil)
}

func (this *CollectController) GetData() mvc.Result {
	id, _ := this.Ctx.URLParamInt64("id")
	th := this.Service.GetCrawlerHead(id)
	return &mvc.View{
		Name: "collect/data.html",
		Data: iris.Map{"id": id, "head": th},
	}
}

func (this *CollectController) PostData() string {
	id, _ := this.Ctx.URLParamInt64("id")
	page, _ := this.Ctx.PostValueInt("page")
	pageSize, _ := this.Ctx.PostValueInt("rows")
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.GetCollectData(&helper.RequestParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
		Id:       id,
	})
	return result
}
