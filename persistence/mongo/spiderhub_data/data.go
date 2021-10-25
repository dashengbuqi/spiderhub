package spiderhub_data

import (
	"context"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type primary struct {
	field string
	value interface{}
}

type CollectData struct {
	collect *qmgo.Collection
	ctx     context.Context
}

func NewCollectData(table string) *CollectData {
	return &CollectData{
		collect: mongo.MongoEngine[mongo.MONGO_DATA].Collection(table),
		ctx:     context.Background(),
	}
}

//创建数据
func (this *CollectData) Build(data map[string]interface{}, method int) error {
	//格式化数据
	pm, ndata := dataFormat(data)
	//需要更新则检查是否存在
	if method == common.METHOD_UPDATE && len(pm.field) > 0 {
		cond := bson.M{pm.field: pm.value}
		amount, _ := this.collect.Find(this.ctx, cond).Count()
		if amount > 0 {
			err := this.collect.UpdateOne(this.ctx, cond, ndata)
			if err != nil {
				return err
			}
			return nil
		}
	}
	_, err := this.collect.InsertOne(this.ctx, ndata)
	if err != nil {
		return err
	}
	return nil
}

//创建数据
func (this *CollectData) BuildClean(data map[string]map[bool]*common.FieldData, method int) error {
	//格式化数据
	pm, ndata := dataCleanFormat(data)
	//需要更新则检查是否存在
	if method == common.METHOD_UPDATE && len(pm.field) > 0 {
		cond := bson.M{pm.field: pm.value}
		amount, _ := this.collect.Find(this.ctx, cond).Count()
		if amount > 0 {
			err := this.collect.UpdateOne(this.ctx, cond, ndata)
			if err != nil {
				return err
			}
			return nil
		}
	}
	_, err := this.collect.InsertOne(this.ctx, ndata)
	if err != nil {
		return err
	}
	return nil
}

//删除表中数据
func (this *CollectData) RemoveRows() error {
	_, err := this.collect.RemoveAll(this.ctx, bson.M{})
	return err
}

func (this *CollectData) GetRowsBy(skip int64, limit int64) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	err := this.collect.Find(this.ctx, bson.M{}).Skip(skip).Limit(limit).All(&items)
	return items, err
}

//删除表
func (this *CollectData) Delete() error {
	return this.collect.DropCollection(this.ctx)
}

func (this *CollectData) Has() bool {
	num, _ := this.collect.Find(this.ctx, bson.M{}).Count()
	return num > 0
}

func (this *CollectData) PostList(req *helper.RequestParams) string {
	var query qmgo.QueryI

	query = this.collect.Find(this.ctx, bson.M{})
	result := this.assembleTable(query, req)
	return result.ToJson()
}

func (this *CollectData) assembleTable(query qmgo.QueryI, req *helper.RequestParams) *helper.ResultEasyUItem {
	pages := &helper.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	var sortStr string

	if len(req.Sort) > 0 {
		sortKeys := strings.Split(req.Sort, ",")
		sortValues := strings.Split(req.Order, ",")
		for i, key := range sortKeys {
			if sortValues[i] == "desc" {
				sortStr += "-" + key + ","
			} else {
				sortStr += "-" + key + ","
			}
		}
	}
	var items []map[string]interface{}
	limit := pages.GetLimit()
	offset := pages.GetOffset()
	sortBy := strings.Trim(sortStr, ",")
	total, _ := query.Count()
	err := query.Limit(int64(limit)).Skip(int64(offset)).Sort(sortBy).All(&items)
	//total, err := query.OrderBy(sortBy).Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		spiderhub.Logger.Error("%v", err.Error())
	}

	pages.Total = total
	return &helper.ResultEasyUItem{
		Pages:  pages,
		Models: items,
	}
}

func dataFormat(data map[string]interface{}) (*primary, interface{}) {
	p := new(primary)
	n := make(map[string]interface{})

	for key, value := range data {
		for isPrimary, val := range value.(map[bool]interface{}) {
			if isPrimary {
				p.field = key
				p.value = val
			}
			n[key] = val
		}
	}
	return p, n
}

func dataCleanFormat(data map[string]map[bool]*common.FieldData) (*primary, interface{}) {
	p := new(primary)
	n := make(map[string]interface{})

	for key, value := range data {
		for isPrimary, val := range value {
			if isPrimary {
				p.field = key
				p.value = val
			}
			n[key] = val
		}
	}
	return p, n
}
