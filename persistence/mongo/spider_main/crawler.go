package spider_main

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	STATUS_NORMAL = iota
	STATUS_RUNNING

	METHOD_INSERT = 1
	METHOD_UPDATE = 2
	METHOD_APPEND = 3
)

type Crawler struct {
	Id           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	UserId       int                `bson:"user_id"`
	CrawlerToken string             `bson:"crawler_token"`
	CleanToken   string             `bson:"clean_token"`
	Status       int                `bson:"status"`   //状态(0完成1执行中)
	Schedule     string             `bson:"schedule"` //计划任务
	Storage      int                `bson:"storage"`  //存储附件(0不存1服务器)
	Method       int                `bson:"method"`   //抓取方式(1重新抓取2更新3追加)
	ErrorInfo    string             `bson:"error_info"`
	SpiderRule   string             `bson:"spider_rule"`
	CleanRule    string             `bson:"clean_rule"`
	UpdatedAt    int64              `bson:"updated_at"`
	CreatedAt    int64              `bson:"created_at"`
}

type CrawlerImpl struct {
	collect *qmgo.Collection
	ctx     context.Context
}

func NewCrawler() *CrawlerImpl {
	return &CrawlerImpl{
		collect: mongo.MongoEngine[mongo.MONGO_MAIN].Collection("crawler"),
		ctx:     context.Background(),
	}
}

//更新爬虫状态
func (this *CrawlerImpl) ModifyStatus(id primitive.ObjectID, state int) error {
	return this.collect.UpdateOne(this.ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": state}})
}

func (this *CrawlerImpl) GetRowByID(id primitive.ObjectID) (*Crawler, error) {
	var item Crawler
	err := this.collect.Find(this.ctx, bson.M{"_id": id}).One(&item)
	if err != nil {
		return &item, err
	}
	return &item, nil
}
