package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type CleanController struct {
	Ctx     iris.Context
	Service services.CleanService
	Session *sessions.Session
}

//加载列表视图
func (this *CleanController) GetList() mvc.Result {
	return &mvc.View{
		Name: "clean/list.html",
	}
}

func (this *CleanController) PostList() string {
	page, _ := this.Ctx.PostValueInt("page")
	pageSize, _ := this.Ctx.PostValueInt("rows")
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.GetCleanList(&helper.RequestParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	})
	return result
}

//正式开始
func (this *CleanController) PutStart() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.CleanStart(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("启动成功", nil)
}

//检查状态
func (this *CleanController) GetStatus() string {
	id, _ := this.Ctx.URLParamInt64("id")
	err := this.Service.CleanStatus(id)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("执行完成", nil)
}

//加载调试视图
func (this *CleanController) GetDebug() mvc.Result {
	id, _ := this.Ctx.URLParamInt64("id")
	model := this.Service.GetRowBy(id)
	return &mvc.View{
		Name: "clean/debug.html",
		Data: iris.Map{"id": id, "status": model.Status, "content": model.CleanContent},
	}
}

//调试开始
func (this *CleanController) PostBegin() string {
	id, _ := this.Ctx.URLParamInt64("id")
	code := this.Ctx.FormValue("code")
	debug_id, err := this.Service.CleanBegin(id, code)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("开始执行", iris.Map{"debug_id": debug_id})
}

//心跳检测
func (this *CleanController) GetHeart() string {
	id, _ := this.Ctx.URLParamInt64("id")
	debug_id, _ := this.Ctx.URLParamInt64("debug_id")
	res := this.Service.CleanHeart(id, debug_id, 0)
	return helper.ResultSuccess("SUCCESS", res)
}

//调试结束
func (this *CleanController) PutEnd() string {
	id, _ := this.Ctx.URLParamInt64("id")
	debug_id, _ := this.Ctx.URLParamInt64("debug_id")
	err := this.Service.CleanEnd(id, debug_id, 0)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("调试正在终止...", nil)
}

func (this *CleanController) PostSave() string {
	id, _ := this.Ctx.URLParamInt64("id")
	code := this.Ctx.FormValue("code")
	err := this.Service.ModifyClean(id, code)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("保存成功", nil)
}

func (this *CleanController) GetData() mvc.Result {
	id, _ := this.Ctx.URLParamInt64("id")
	th := this.Service.GetCleanHead(id)
	return &mvc.View{
		Name: "clean/data.html",
		Data: iris.Map{"id": id, "head": th},
	}
}

func (this *CleanController) PostData() string {
	id, _ := this.Ctx.URLParamInt64("id")
	page, _ := this.Ctx.PostValueInt("page")
	pageSize, _ := this.Ctx.PostValueInt("rows")
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.GetCleanData(&helper.RequestParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
		Id:       id,
	})
	return result
}
