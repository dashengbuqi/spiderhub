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
	"time"
)

type CleanService interface {
	GetRowBy(id int64) *collect.Application
	GetCleanList(post *helper.RequestParams) string
	GetCleanData(post *helper.RequestParams) string
	ModifyClean(id int64, content string) error
	CleanBegin(id int64, code string) (int64, error)
	CleanHeart(id int64, debug_id int64) interface{}
	CleanEnd(id int64, debug_id int64) error
	CleanStart(id int64) error
	CleanStatus(id int64) error
	GetCleanHead(id int64) []*common.TableHead
}

type cleanService struct {
	repo collect.ApplicationImp
}

func NewCleanService() CleanService {
	return &cleanService{
		repo: collect.NewApplication(),
	}
}

func (this *cleanService) GetRowBy(id int64) *collect.Application {
	result, _ := this.repo.GetRowByID(id)
	return result
}

func (this *cleanService) GetCleanList(post *helper.RequestParams) string {
	result := this.repo.PostList(post)
	return result
}

func (this *cleanService) GetCleanData(post *helper.RequestParams) string {
	if post.Id == 0 {
		return "{\"total\":0,\"rows\":{}}"
	}
	item, _ := this.repo.GetRowByID(post.Id)
	if len(item.CrawlerToken) == 0 {
		return "{\"total\":0,\"rows\":{}}"
	}
	dataTable := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_DATA, item.CleanToken)
	d := spiderhub_data.NewCollectData(dataTable)
	result := d.PostList(post)
	return result
}

func (this *cleanService) ModifyClean(id int64, content string) error {
	if id == 0 {
		return errors.New("请先创建采集任务")
	}
	return this.repo.ModifyCleanContent(id, content)
}

func (this *cleanService) CleanStatus(id int64) error {
	ca := collect.NewApplication()
	row, _ := ca.GetRowByID(id)
	if row.Status == common.STATUS_RUNNING {
		return errors.New("正在执行中")
	}
	return nil
}

//正式执行
func (this *cleanService) CleanStart(id int64) error {
	if id == 0 {
		return errors.New("缺少参数")
	}
	ca := collect.NewApplication()
	row, _ := ca.GetRowByID(id)
	if len(row.CleanContent) == 0 {
		return errors.New("还没有爬虫规则,点击调试创建吧！")
	}
	cm := &common.Communication{
		AppId:   row.Id,
		UserId:  row.UserId,
		DebugId: 0,
		Method:  common.METHOD_EXCUTE,
		Content: row.CleanContent,
	}
	str, err := json.Marshal(cm)
	if err != nil {
		return err
	}
	err = queue.RabbitConn.Publish(&queue.CleanerChannel, str)
	return err
}

func (this *cleanService) GetCleanHead(id int64) []*common.TableHead {
	af := collect.NewAppField()
	item, _ := af.GetRowByID(collect.TARGET_CLEAN, id)
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
func (this *cleanService) CleanBegin(id int64, code string) (int64, error) {
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
	err = queue.RabbitConn.Publish(&queue.CleanerChannel, str)
	return debug_id, err
}

func (this *cleanService) CleanHeart(id, debug_id int64) interface{} {
	if id == 0 || debug_id == 0 {
		return ""
	}
	model, _ := this.repo.GetRowByID(id)
	var logList []common.LogLevel
	var dataList []map[string]interface{}
	token := helper.NewToken(model.UserId, id, debug_id).Clean().ToString()
	dataTable := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_DATA, token)
	logTable := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_LOG, token)
	lcnl := &queue.Channel{
		Exchange:     "Cleaners",
		ExchangeType: "direct",
		RoutingKey:   logTable,
		Reliable:     true,
		Durable:      false,
		AutoDelete:   true,
	}
	dcnl := &queue.Channel{
		Exchange:     "Cleaners",
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
func (this *cleanService) CleanEnd(id int64, debug_id int64) error {
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
	err := queue.RabbitConn.Publish(&queue.CleanerChannel, str)
	if err != nil {
		return err
	}
	return nil
}
