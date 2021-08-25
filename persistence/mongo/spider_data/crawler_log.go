package spider_data

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CrawlerLog struct {
	collect *qmgo.Collection
	ctx     context.Context
}

func NewCrawlerLog(table string) *CrawlerLog {
	return &CrawlerLog{
		collect: mongo.MongoEngine[mongo.MONGO_DATA].Collection(table),
		ctx:     context.TODO(),
	}
}

//创建数据
func (this *CrawlerLog) Build(doc interface{}) (primitive.ObjectID, error) {
	res, err := this.collect.InsertOne(this.ctx, doc)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
