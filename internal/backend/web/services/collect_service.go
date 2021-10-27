package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_data"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"strconv"
	"time"
)

type CollectService interface {
	GetRowBy(id int64) *collect.Application
	GetCollectList(post *helper.RequestParams) string
	GetCollectData(post *helper.RequestParams) string
	ModifyCollectItem(id int64, user_id int64, form map[string][]string) error
	ModifyCrawler(id int64, content string) error
	Remove(id int64) error
	CrawlerBegin(id int64, code string) (int64, error)
	CrawlerHeart(id int64, debug_id int64) interface{}
	CrawlerEnd(id int64, debug_id int64) error
	CrawlerStart(id int64) error
	CrawlerStatus(id int64) error
	GetCrawlerHead(id int64) []*common.TableHead
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

func (this *collectService) GetCollectData(post *helper.RequestParams) string {
	if post.Id == 0 {
		return "{\"total\":0,\"rows\":{}}"
	}
	item, _ := this.repo.GetRowByID(post.Id)
	if len(item.CrawlerToken) == 0 {
		return "{\"total\":0,\"rows\":{}}"
	}
	dataTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, item.CrawlerToken)
	d := spiderhub_data.NewCollectData(dataTable)
	result := d.PostList(post)
	return result
}

func (this *collectService) ModifyCrawler(id int64, content string) error {
	if id == 0 {
		return errors.New("请先创建采集任务")
	}
	return this.repo.ModifyCrawlerContent(id, content)
}

//更新数据
func (this *collectService) ModifyCollectItem(id int64, user_id int64, form map[string][]string) error {
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
	return this.repo.ModifyItem(id, user_id, &collect.Application{
		Title:    title,
		Schedule: schedule,
		Storage:  storageInt,
		Method:   methodInt,
	})
}
func (this *collectService) CrawlerStatus(id int64) error {
	ca := collect.NewApplication()
	row, _ := ca.GetRowByID(id)
	if row.Status == common.STATUS_RUNNING {
		return errors.New("正在执行中")
	}
	return nil
}
func (this *collectService) Remove(id int64) error {
	if id == 0 {
		return errors.New("暂不支持")
	}
	//删除对应的数据
	item, _ := this.repo.GetRowByID(id)
	if item.Id > 0 {
		if len(item.CrawlerToken) > 0 {
			dataTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, item.CrawlerToken)
			cd := spiderhub_data.NewCollectData(dataTable)
			cd.Delete()
			//日志表
			logTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_LOG, item.CrawlerToken)
			cc := spiderhub_data.NewCollectLog(logTable)
			cc.Delete()
		}
		if len(item.CleanToken) > 0 {
			dataTable := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_DATA, item.CleanToken)
			cd := spiderhub_data.NewCollectData(dataTable)
			cd.Delete()
			//日志表
			logTable := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_LOG, item.CleanToken)
			cc := spiderhub_data.NewCollectLog(logTable)
			cc.Delete()
		}
		return this.repo.Remove(id)
	}
	return errors.New("采集数据不存在")
}

//正式执行
func (this *collectService) CrawlerStart(id int64) error {
	if id == 0 {
		return errors.New("缺少参数")
	}
	ca := collect.NewApplication()
	row, _ := ca.GetRowByID(id)
	if len(row.CrawlerContent) == 0 {
		return errors.New("还没有爬虫规则,点击调试创建吧！")
	}
	cm := &common.Communication{
		AppId:   row.Id,
		UserId:  row.UserId,
		DebugId: 0,
		Method:  common.METHOD_EXCUTE,
		Content: row.CrawlerContent,
	}
	str, err := json.Marshal(cm)
	if err != nil {
		return err
	}
	err = queue.RabbitConn.Publish(&queue.CrawlerChannel, str)
	return err
}

func (this *collectService) GetCrawlerHead(id int64) []*common.TableHead {
	af := collect.NewAppField()
	item, _ := af.GetRowByID(collect.TARGET_CRAWLER, id)
	var content []*common.TableHead
	if len(item.Content) > 0 {
		err := json.Unmarshal([]byte(item.Content), &content)
		if err != nil {
			return nil
		}
	}
	return content
}

//调试开始
func (this *collectService) CrawlerBegin(id int64, code string) (int64, error) {
	if id == 0 {
		return 0, errors.New("缺少参数")
	}
	if len(code) == 0 {
		return 0, errors.New("未获取到采集规则")
	}
	row, _ := this.repo.GetRowByID(id)
	debug_id := time.Now().Unix()
	cm := &common.Communication{
		AppId:   id,
		UserId:  row.UserId,
		DebugId: debug_id,
		Method:  common.METHOD_DEBUG,
		Content: code,
	}
	str, err := json.Marshal(cm)
	if err != nil {
		return debug_id, err
	}
	err = queue.RabbitConn.Publish(&queue.CrawlerChannel, str)
	return debug_id, err
}

func (this *collectService) CrawlerHeart(id, debug_id int64) interface{} {
	if id == 0 || debug_id == 0 {
		return ""
	}
	model, _ := this.repo.GetRowByID(id)
	var logList []common.LogLevel
	var dataList []map[string]interface{}
	token := helper.NewToken(model.UserId, id, debug_id).Crawler().ToString()
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

	logOut := make(chan []byte)
	dataOut := make(chan []byte)
	var err error
	go queue.RabbitConn.Consume(lcnl, logOut)
	go queue.RabbitConn.Consume(dcnl, dataOut)
	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case l := <-logOut:
			var item common.LogLevel
			err = json.Unmarshal(l, &item)
			if err == nil {
				logList = append(logList, item)
			}
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
		Logs    []common.LogLevel        `json:"logs"`
		Rows    []map[string]interface{} `json:"rows"`
	}{
		Id:      id,
		Status:  model.Status,
		DebugId: debug_id,
		Logs:    logList,
		Rows:    dataList,
	}
	return res
}

//终止调试
func (this *collectService) CrawlerEnd(id int64, debug_id int64) error {
	if id == 0 || debug_id == 0 {
		return errors.New("缺少参数")
	}
	row, _ := this.repo.GetRowByID(id)
	cm := &common.Communication{
		AppId:   id,
		UserId:  row.UserId,
		DebugId: debug_id,
		Method:  common.METHOD_DEBUG,
		Abort:   true,
	}
	str, _ := json.Marshal(cm)
	err := queue.RabbitConn.Publish(&queue.CrawlerChannel, str)
	if err != nil {
		return err
	}
	return nil
}
