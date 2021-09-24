package services

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
)

type MenuService interface {
	GetLevelList(parant_id int64) string
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
func (this *menuService) GetLevelList(parent_id int64) string {
	data := this.repo.GetRowsData(parent_id)
	jStr, _ := json.Marshal(data)
	return string(jStr)
}
