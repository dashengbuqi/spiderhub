package spiderhub_data

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectLog struct {
	collect *qmgo.Collection
	ctx     context.Context
}

func NewCollectLog(table string) *CollectLog {
	return &CollectLog{
		collect: mongo.MongoEngine[mongo.MONGO_DATA].Collection(table),
		ctx:     context.TODO(),
	}
}

//创建数据
func (this *CollectLog) Build(doc interface{}) (primitive.ObjectID, error) {
	res, err := this.collect.InsertOne(this.ctx, doc)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

//删除表中数据
func (this *CollectLog) RemoveRows() error {
	_, err := this.collect.RemoveAll(this.ctx, bson.M{})
	return err
}

//删除表
func (this *CollectLog) Delete() error {
	return this.collect.DropCollection(this.ctx)
}
