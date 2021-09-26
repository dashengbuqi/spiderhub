package tool

import (
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/go-xorm/xorm"
	"strings"
)

//组装管理系统表单列表数据
func AssembleTable(query *xorm.Session, params map[string]interface{}) map[string]interface{} {
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
	var items []interface{}
	limit := pages.GetLimit()
	offset := pages.GetOffset()
	sortBy := strings.Trim(sortStr, ",")

	total, err := query.OrderBy(sortBy).Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		spiderhub.Logger.Error("%v", err.Error())
	}
	pages.Total = total
	return map[string]interface{}{
		"pages":  pages,
		"models": items,
	}
}
