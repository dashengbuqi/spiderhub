package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type CleanController struct {
	Ctx     iris.Context
	Service services.CleanService
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

func (this *CleanController) PostSave() string {
	id, _ := this.Ctx.URLParamInt64("id")
	code := this.Ctx.FormValue("code")
	err := this.Service.ModifyClean(id, code)
	if err != nil {
		return helper.ResultError(err.Error())
	}
	return helper.ResultSuccess("保存成功", nil)
}
