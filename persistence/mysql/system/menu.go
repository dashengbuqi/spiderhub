package system

import (
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
)

const (
	MENU_TYPE_CATALOG = 1
	MENU_TYPE_COLUMN  = 2
	MENU_TYPE_MENU    = 3
	MENU_TYPE_BUTTON  = 4
	MENU_TYPE_EVENT   = 5

	MENU_STATUS_DISABLE = 0
	MENU_STATUS_ENABLE  = 1
)

type SystemMenu struct {
	Id        int64  `json:"id"`
	TaskName  string `json:"task_name"`
	FullName  string `json:"full_name"`
	Path      string `json:"path"`
	ParentId  int64  `json:"parent_id"`
	Type      int    `json:"type"`
	Icon      string `json:"icon"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `json:"created_at"`
}

type Menu struct {
	session *xorm.Engine
}

func NewMenu() *Menu {
	return &Menu{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}

//获取菜单
func (this *Menu) GetRowsData(parent_id int64) (result []map[string]interface{}) {
	var items []SystemMenu
	err := this.session.Where("parent_id = ? AND status = ? AND type <> ?", parent_id, MENU_STATUS_ENABLE, MENU_TYPE_BUTTON).
		OrderBy("sort").
		Find(&items)
	if err != nil {
		return result
	}
	for _, item := range items {
		temp := map[string]interface{}{
			"id":        item.Id,
			"task_name": item.TaskName,
			"icon":      item.Icon,
			"task_url":  item.Path,
			"children":  this.GetRowsData(item.Id),
		}
		result = append(result, temp)
	}
	return result
}
