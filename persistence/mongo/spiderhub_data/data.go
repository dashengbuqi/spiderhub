package spiderhub_data

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_main"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type primary struct {
	field string
	value interface{}
}

type CrawlerData struct {
	collect *qmgo.Collection
	ctx     context.Context
}

func NewCrawlerData(table string) *CrawlerData {
	return &CrawlerData{
		collect: mongo.MongoEngine[mongo.MONGO_DATA].Collection(table),
		ctx:     context.Background(),
	}
}

//创建数据
func (this *CrawlerData) Build(data map[string]interface{}, method int) error {
	//格式化数据
	pm, ndata := dataFormat(data)
	//需要更新则检查是否存在
	if method == spiderhub_main.METHOD_UPDATE && len(pm.field) > 0 {
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
func (this *CrawlerData) RemoveRows() error {
	_, err := this.collect.RemoveAll(this.ctx, bson.M{})
	return err
}

func (this *CrawlerData) GetRowsBy(skip int64, limit int64) ([]interface{}, error) {
	var items []interface{}
	err := this.collect.Find(this.ctx, bson.M{}).Skip(skip).Limit(limit).All(&items)
	return items, err
}

//删除表
func (this *CrawlerData) Delete() error {
	return this.collect.DropCollection(this.ctx)
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
