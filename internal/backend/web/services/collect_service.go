package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"strconv"
	"time"
)

type CollectService interface {
	GetRowBy(id int64) *collect.Application
	GetCollectList(post *helper.RequestParams) string
	ModifyCollectItem(id int64, form map[string][]string) error
	ModifyCrawler(id int64, content string) error
	Remove(id int64) error
	CrawlerBegin(id int64, code string) (int64, error)
	CrawlerHeart(id int64, debug_id int64, user_id int64) interface{}
	CrawlerEnd(id int64, debug_id int64, user_id int64) error
}

type collectService struct {
	repo collect.ApplicationImp
}

func NewCollectService() CollectService {
	return &collectService{
		repo: collect.NewApplication(),
	}
}

func (this *collectService) GetRowBy(id int64) *collect.Application {
	result, _ := this.repo.GetRowByID(id)
	return result
}

func (this *collectService) GetCollectList(post *helper.RequestParams) string {
	result := this.repo.PostList(post)
	return result
}
func (this *collectService) ModifyCrawler(id int64, content string) error {
	if id == 0 {
		return errors.New("请先创建采集任务")
	}
	return this.repo.ModifyCrawlerContent(id, content)
}

//更新数据
func (this *collectService) ModifyCollectItem(id int64, form map[string][]string) error {
	var title, schedule, storage, method string
	if _, ok := form["title"]; ok {
		title = form["title"][0]
	}
	if _, ok := form["schedule"]; ok {
		schedule = form["schedule"][0]
	}
	if _, ok := form["storage"]; ok {
		storage = form["storage"][0]
	}
	if _, ok := form["method"]; ok {
		method = form["method"][0]
	}
	storageInt, _ := strconv.Atoi(storage)
	methodInt, _ := strconv.Atoi(method)
	return this.repo.ModifyItem(id, &collect.Application{
		Title:    title,
		Schedule: schedule,
		Storage:  storageInt,
		Method:   methodInt,
	})
}

func (this *collectService) Remove(id int64) error {
	if id == 0 {
		return errors.New("暂不支持")
	}
	return this.repo.Remove(id)
}

//调试开始
func (this *collectService) CrawlerBegin(id int64, code string) (int64, error) {
	if id == 0 {
		return 0, errors.New("缺少参数")
	}
	if len(code) == 0 {
		return 0, errors.New("未获取到采集规则")
	}
	debug_id := time.Now().Unix()
	cm := &common.Communication{
		AppId:   id,
		UserId:  0,
		DebugId: debug_id,
		Method:  0,
		Content: code,
	}
	str, err := json.Marshal(cm)
	if err != nil {
		return debug_id, err
	}
	fmt.Println("进队列")
	err = queue.RabbitConn.Publish(&common.CrawlerChannel, str)
	fmt.Println(err)
	return debug_id, err
}

func (this *collectService) CrawlerHeart(id, debug_id, user_id int64) interface{} {
	if id == 0 || debug_id == 0 {
		return ""
	}
	model, _ := this.repo.GetRowByID(id)
	var logList []common.LogLevel
	var dataList []map[string]interface{}
	token := helper.NewToken(user_id, id, debug_id).Crawler().ToString()
	dataTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, token)
	logTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_LOG, token)

	lcnl := &queue.Channel{
		Exchange:     "Crawlers",
		ExchangeType: "direct",
		RoutingKey:   logTable,
		Reliable:     true,
		Durable:      false,
		AutoDelete:   true,
	}
	dcnl := &queue.Channel{
		Exchange:     "Crawlers",
		ExchangeType: "direct",
		RoutingKey:   dataTable,
		Reliable:     true,
		Durable:      false,
		AutoDelete:   true,
	}

	logOut := make(chan []byte, 5)
	dataOut := make(chan []byte)
	var err error
	go queue.RabbitConn.Consume(lcnl, logOut)
	go queue.RabbitConn.Consume(dcnl, dataOut)
	t := time.NewTicker(time.Second)
	for {
		select {
		case l := <-logOut:
			if len(logList) == 5 {
				goto Loop
			}
			var item common.LogLevel
			err = json.Unmarshal(l, &item)
			if err != nil {
				goto Loop
			}
			logList = append(logList, item)
			if item.Type == common.LOG_TYPE_FINISH {
				goto Loop
			}
		case d := <-dataOut:
			var item map[string]interface{}
			err = json.Unmarshal(d, &item)
			if err == nil {
				dataList = append(dataList, item)
			}
		case <-t.C:
			goto Loop
		}
	}
Loop:
	res := struct {
		Id      int64                    `json:"id"`
		Status  int                      `json:"status"`
		DebugId int64                    `json:"debug_id"`
		Log     []common.LogLevel        `json:"log"`
		Rows    []map[string]interface{} `json:"rows"`
	}{
		Id:      id,
		Status:  model.Status,
		DebugId: debug_id,
		Log:     logList,
		Rows:    dataList,
	}
	return res
}

//终止调试
func (this *collectService) CrawlerEnd(id int64, debug_id int64, user_id int64) error {
	if id == 0 || debug_id == 0 {
		return errors.New("缺少参数")
	}
	cm := &common.Communication{
		AppId:   id,
		UserId:  0,
		DebugId: debug_id,
		Method:  0,
		Abort:   true,
	}
	str, _ := json.Marshal(cm)
	err := queue.RabbitConn.Publish(&common.CrawlerChannel, str)
	if err != nil {
		return err
	}
	return nil
}
