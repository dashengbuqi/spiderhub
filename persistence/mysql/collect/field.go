package collect

import (
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
	"strings"
	"time"
)

const (
	TARGET_CRAWLER = 1
	TARGET_CLEAN   = 2
)

type AppFields struct {
	Id        int64  `bson:"id"`
	Target    int    `bson:"target"`
	TargetId  int64  `bson:"target_id"`
	Content   string `bson:"content"`
	UpdatedAt int64  `bson:"updated_at"`
	CreatedAt int64  `bson:"created_at"`
}

func (AppFields) TableName() string {
	return "collect_app_fields"
}

type field struct {
	session *xorm.Engine
}

func NewAppField() *field {
	return &field{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}
func (this *field) GetRowByID(target int, id int64) (*AppFields, error) {
	var item AppFields
	_, err := this.session.Where("target = ? AND target_id = ?", target, id).Get(&item)
	return &item, err
}

//更新数据
func (this *field) Modify(target int, target_id int64, content []byte) error {
	var item AppFields

	has, err := this.session.Where("target=? AND target_id=?", target, target_id).Get(&item)
	if err != nil {
		return err
	}
	tm := time.Now().Unix()
	if has {
		item.Content = string(content)
		item.UpdatedAt = tm
		_, err = this.session.Where("id=?", item.Id).Cols("content", "updated_at").Update(item)
	} else {
		item.Target = target
		item.TargetId = target_id
		item.Content = string(content)
		item.UpdatedAt = tm
		item.CreatedAt = tm
		_, err = this.session.InsertOne(item)
	}
	return err
}

func (this *field) PostList(req *helper.RequestParams) string {
	var query *xorm.Session

	query = this.session.Table(new(AppFields))
	result := this.assembleTable(query, req)
	return result.ToJson()
}

func (this *field) assembleTable(query *xorm.Session, req *helper.RequestParams) *helper.ResultEasyUItem {
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
	var items []*AppFields
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
