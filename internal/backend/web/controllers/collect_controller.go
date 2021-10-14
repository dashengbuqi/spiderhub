package controllers

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
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
