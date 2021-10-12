package system

import (
	"fmt"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
	"strings"
)

type SystemAdmin struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	RealName   string `json:"real_name"`
	Mobile     string `json:"mobile"`
	AuthKey    string `json:"auth_key"`
	Password   string `json:"password" xorm:"password_hash"`
	Email      string `json:"email"`
	Status     int    `json:"status"`
	LoginTimes int64  `json:"login_times"`
	UpdatedAt  int64  `json:"updated_at"`
	CreatedAt  int64  `json:"created_at"`
}

type Admin struct {
	session *xorm.Engine
}

func NewAdmin() *Admin {
	return &Admin{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}

type SystemAdminBackend struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	RealName    string `json:"real_name"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Status      int    `json:"status"`
	UiStatus    string `json:"ui_status" xorm:"-"`
	LoginTimes  int64  `json:"login_times"`
	UpdatedAt   int64  `json:"updated_at"`
	UiUpdatedAt string `json:"ui_updated_at" xorm:"-"`
	CreatedAt   int64  `json:"created_at"`
	UiCreatedAt string `json:"ui_created_at" xorm:"-"`
}

func (this *SystemAdminBackend) callUI() {
	this.UiCreatedAt = helper.FmtDateTime(this.CreatedAt)
	this.UiStatus = StatusArr[this.Status]
	this.UiUpdatedAt = helper.FmtDateTime(this.UpdatedAt)
}
func (this *Admin) GetRowBy(id int64) (*SystemAdmin, error) {
	var item SystemAdmin
	_, err := this.session.Where("id = ?", id).Get(&item)
	return &item, err
}

//加载菜载列表数据
func (this *Admin) PostMenuList(req *helper.RequestParams) string {
	var query *xorm.Session
	var where string
	if req.Params != nil {
		params := req.Params.(map[string]interface{})
		if _, ok := params["username"]; ok {
			where = fmt.Sprintf("username like '%?%'", strings.TrimSpace(params["username"].(string)))
		}
	}

	query = this.session.Table(new(SystemAdmin)).Where(where)
	result := this.assembleTable(query, req)
	return result.ToJson()
}

func (this *Admin) assembleTable(query *xorm.Session, req *helper.RequestParams) *helper.ResultEasyUItem {
	pages := &helper.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	var sortStr string

	if len(req.Sort) > 0 {
		sortKeys := strings.Split(req.Sort, ",")
		sortValues := strings.Split(req.Order, ",")
		for i, key := range sortKeys {
			sortStr += key + " " + sortValues[i] + ","
		}
	}
	var items []*SystemAdminBackend
	limit := pages.GetLimit()
	offset := pages.GetOffset()
	sortBy := strings.Trim(sortStr, ",")

	total, err := query.OrderBy(sortBy).Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		spiderhub.Logger.Error("%v", err.Error())
	}
	for _, item := range items {
		item.callUI()
	}
	pages.Total = total
	return &helper.ResultEasyUItem{
		Pages:  pages,
		Models: items,
	}
}
