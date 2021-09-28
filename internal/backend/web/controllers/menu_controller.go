package controllers

import (
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
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

func (this *MenuController) GetEdit() mvc.Result {
	id, _ := this.Ctx.URLParamInt64("id")
	model := this.Service.GetRowBy(id)
	return &mvc.View{
		Name: "menu/edit.html",
		Data: iris.Map{
			"id":        model.Id,
			"task_name": model.TaskName,
			"full_name": model.FullName,
			"path":      model.Path,
			"icon":      model.Icon,
			"sort":      model.Sort,
			"type":      model.Type,
			"parent_id": model.ParentId,
		},
	}
}
