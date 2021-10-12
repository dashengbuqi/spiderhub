package services

import (
	"encoding/json"
	"errors"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
	"strconv"
)

type MenuService interface {
	GetLevelMenu(parant_id int64) string
	PostMenuList(post *helper.RequestParams) string
	GetRowBy(id int64) *system.SystemMenu
	ModifyMenuItem(id int64, form map[string][]string) error
	RemoveMenu(id int64) error
}

type menuService struct {
	repo *system.Menu
}

func NewMenuService() MenuService {
	return &menuService{
		repo: system.NewMenu(),
	}
}

func (this *menuService) RemoveMenu(id int64) error {
	if id == 0 {
		return errors.New("暂不支持")
	}
	return this.repo.RemoveItem(id)
}

//更新数据
func (this *menuService) ModifyMenuItem(id int64, form map[string][]string) error {
	var task_name, full_name, path, icon string
	var sort, tp, parent_id int
	if _, ok := form["task_name"]; ok {
		task_name = form["task_name"][0]
	}
	if _, ok := form["full_name"]; ok {
		full_name = form["full_name"][0]
	}
	if _, ok := form["path"]; ok {
		path = form["path"][0]
	}
	if _, ok := form["icon"]; ok {
		icon = form["icon"][0]
	}
	if _, ok := form["sort"]; ok {
		sort, _ = strconv.Atoi(form["sort"][0])
	}
	if _, ok := form["type"]; ok {
		tp, _ = strconv.Atoi(form["type"][0])
	}
	if _, ok := form["parent_id"]; ok {
		parent_id, _ = strconv.Atoi(form["parent_id"][0])
	}
	return this.repo.ModifyItem(id, &system.SystemMenu{
		TaskName: task_name,
		FullName: full_name,
		Path:     path,
		Icon:     icon,
		Sort:     sort,
		Type:     tp,
		ParentId: int64(parent_id),
	})
}

//获取规范的菜单层级
func (this *menuService) GetLevelMenu(parent_id int64) string {
	data := this.repo.GetRowsData(parent_id)
	jStr, _ := json.Marshal(data)
	return string(jStr)
}

func (this *menuService) PostMenuList(post *helper.RequestParams) string {
	result := this.repo.PostMenuList(post)
	return result
}

func (this *menuService) GetRowBy(id int64) *system.SystemMenu {
	result, _ := this.repo.GetRowBy(id)
	return result
}
