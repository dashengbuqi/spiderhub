package spiderhub_main

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AppFields struct {
	Id        primitive.ObjectID `bson:"_id"`
	Target    int                `bson:"target"`
	TargetId  primitive.ObjectID `bson:"target_id"`
	Content   string             `bson:"content"`
	UpdatedAt int64              `bson:"updated_at"`
	CreatedAt int64              `bson:"created_at"`
}

type AppFieldsImpl struct {
	Collect *qmgo.Collection
	Ctx     context.Context
}

func NewAppField() *AppFieldsImpl {
	return &AppFieldsImpl{
		Collect: mongo.MongoEngine[mongo.MONGO_MAIN].Collection("fields"),
		Ctx:     context.Background(),
	}
}

//更新数据
func (this *AppFieldsImpl) Modify(target int, appid primitive.ObjectID, content []byte) error {
	var item AppFields

	err := this.Collect.Find(this.Ctx, bson.M{"target": target, "target_id": appid}).One(&item)
	if err != nil {
		return err
	}
	tm := time.Now().Unix()
	if item.Id == primitive.NilObjectID {
		item.Id = primitive.NewObjectID()
		item.Target = target
		item.TargetId = appid
		item.Content = string(content)
		item.UpdatedAt = tm
		item.CreatedAt = tm
		_, err := this.Collect.InsertOne(this.Ctx, item)
		if err != nil {
			return err
		}
	} else {
		err := this.Collect.UpdateOne(this.Ctx,
			bson.M{"_id": item.Id},
			bson.M{"$set": bson.M{"content": string(content), "updated_at": tm}},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
