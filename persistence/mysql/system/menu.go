package system

import (
	"fmt"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
	"strings"
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

func (this *Menu) AttributeLabels(attribute string) string {
	attributes := map[string]string{
		"id":         "序号",
		"task_name":  "名称",
		"full_name":  "全称",
		"path":       "接口地址",
		"parent_id":  "上级",
		"type":       "类型",
		"icon":       "图标",
		"status":     "状态",
		"sort":       "排序",
		"updated_at": "更新",
		"created_at": "创建",
	}
	return attributes[attribute]
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

//加载菜载列表数据
func (this *Menu) PostMenuList(post map[string]interface{}) string {
	var query *xorm.Session
	var where string
	if _, ok := post["task_name"]; ok {
		where = fmt.Sprintf("task_name like '%?%'", strings.TrimSpace(post["task_name"].(string)))
	}
	query = this.session.Table(new(SystemMenu)).Where(where)
	result := this.assembleTable(query, post)
	return result.ToJson()
}

func (this *Menu) assembleTable(query *xorm.Session, params map[string]interface{}) *helper.ResultEasyUItem {
	page := params["page"].(int)
	pageSize := params["pageSize"].(int)

	pages := &helper.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
	var sortStr string
	if _, ok := params["sort"]; ok {
		if len(params["sort"].(string)) > 0 {
			sortKeys := strings.Split(params["sort"].(string), ",")
			sortValues := strings.Split(params["order"].(string), ",")
			for i, key := range sortKeys {
				sortStr += key + " " + sortValues[i] + ","
			}
		}
	}
	var items []SystemMenu
	limit := pages.GetLimit()
	offset := pages.GetOffset()
	sortBy := strings.Trim(sortStr, ",")

	total, err := query.OrderBy(sortBy).Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		spiderhub.Logger.Error("%v", err.Error())
	}
	pages.Total = total
	return &helper.ResultEasyUItem{
		Pages:  pages,
		Models: items,
	}
}
