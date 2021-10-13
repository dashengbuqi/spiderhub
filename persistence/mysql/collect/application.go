package collect

import (
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
)

const (
	STATUS_NORMAL = iota
	STATUS_RUNNING

	METHOD_INSERT = 1
	METHOD_UPDATE = 2
	METHOD_APPEND = 3
)

type CollectApplication struct {
	Id             int64  `bson:"id"`
	Title          string `bson:"title"`
	UserId         int64  `bson:"user_id"`
	CrawlerToken   string `bson:"crawler_token"`
	CleanToken     string `bson:"clean_token"`
	Status         int    `bson:"status"`   //状态(0完成1执行中)
	Schedule       string `bson:"schedule"` //计划任务
	Storage        int    `bson:"storage"`  //存储附件(0不存1服务器)
	Method         int    `bson:"method"`   //抓取方式(1重新抓取2更新3追加)
	ErrorInfo      string `bson:"error_info"`
	CrawlerContent string `bson:"crawler_content"`
	CleanContent   string `bson:"clean_content"`
	UpdatedAt      int64  `bson:"updated_at"`
	CreatedAt      int64  `bson:"created_at"`
}

func (this CollectApplication) TableName() string {
	return "collect_app"
}

type application struct {
	session *xorm.Engine
}

func NewApplication() *application {
	return &application{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}

//更新爬虫状态
func (this *application) ModifyStatus(id int64, state int) error {
	var item CollectApplication
	item.Status = state
	_, err := this.session.Where("id=?", id).Cols("status").Update(item)
	return err
}

func (this *application) GetRowByID(id int64) (*CollectApplication, error) {
	var item *CollectApplication
	_, err := this.session.Where("id=?", id).Get(&item)
	return item, err
}
