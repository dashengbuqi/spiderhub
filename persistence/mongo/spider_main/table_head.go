package spider_main

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableHead struct {
	Id        primitive.ObjectID `bson:"_id"`
	Target    int32              `bson:"target"`
	TargetId  int64              `bson:"target_id"`
	Content   string             `bson:"content"`
	UpdatedAt int64              `bson:"updated_at"`
	CreatedAt int64              `bson:"created_at"`
}

type TableHeadImpl struct {
	Collect *qmgo.Collection
	Ctx     context.Context
}

func NewTableHead() *TableHeadImpl {
	return &TableHeadImpl{
		Collect: mongo.MongoEngine[mongo.MONGO_MAIN].Collection("table_head"),
		Ctx:     context.Background(),
	}
}
