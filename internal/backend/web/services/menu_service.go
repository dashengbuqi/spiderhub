package services

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
)

type MenuService interface {
	GetLevelMenu(parant_id int64) string
	PostMenuList(post map[string]interface{}) string
	GetRowBy(id int64) *system.SystemMenu
}

type menuService struct {
	repo *system.Menu
}

func NewMenuService() MenuService {
	return &menuService{
		repo: system.NewMenu(),
	}
}

//获取规范的菜单层级
func (this *menuService) GetLevelMenu(parent_id int64) string {
	data := this.repo.GetRowsData(parent_id)
	jStr, _ := json.Marshal(data)
	return string(jStr)
}

func (this *menuService) PostMenuList(post map[string]interface{}) string {
	result := this.repo.PostMenuList(post)
	return result
}

func (this *menuService) GetRowBy(id int64) *system.SystemMenu {
	result, _ := this.repo.GetRowBy(id)
	return result
}
