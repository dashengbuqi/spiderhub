package controllers

import (
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
		Name:     "form[type]",
		Value:    system.TypeArr[model.Type],
		Data:     model.GetTypeComboList(),
		Multiple: false,
	})
	m := system.NewMenu()
	this.Ctx.ViewData("parents", widgets.Combobox{
		Id:       "parent_id",
		Name:     "form[parent_id]",
		Value:    model.ParentId,
		Editable: false,
		Data:     m.GetMenuTreeList(),
	})
	this.Ctx.View("menu/edit.html")
	/*return &mvc.View{
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
	}*/
}
