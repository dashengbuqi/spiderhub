package spider_main

import (
	"context"
	"github.com/dashengbuqi/spiderhub/middleware/mongo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Crawler struct {
	Id           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	UserId       int                `bson:"user_id"`
	CrawlerToken string             `bson:"crawler_token"`
	CleanToken   string             `bson:"clean_token"`
	Status       int                `bson:"status"`   //状态(0完成1执行中)
	Schedule     string             `bson:"schedule"` //计划任务
	Storage      int                `bson:"storage"`
	Method       int                `bson:"method"` //抓取方式(1重新抓取2更新3追加)
	ErrorInfo    string             `bson:"error_info"`
	SpiderRule   string             `bson:"spider_rule"`
	CleanRule    string             `bson:"clean_rule"`
	UpdatedAt    int64              `bson:"updated_at"`
	CreatedAt    int64              `bson:"created_at"`
}

type CrawlerImpl struct {
	Collect *qmgo.Collection
	Ctx     context.Context
}

func NewCrawler() *CrawlerImpl {
	return &CrawlerImpl{
		Collect: mongo.MongoEngine[mongo.MONGO_MAIN].Collection("crawler"),
		Ctx:     context.Background(),
	}
}
