package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spider_data"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spider_main"
	"github.com/robertkrimen/otto"
	"time"
)

type Schedule struct {
	inData     *common.Communication       //传输过来的数据
	outLog     chan []byte                 //日志输出
	outData    chan map[string]interface{} //数据输出
	rabbitConn *queue.Base                 //rabbit链接
	logQueue   *queue.Channel              //日志队列
	dataQueue  *queue.Channel              //数据队列
	bean       *spider_main.Crawler        //数据持久化实例
	container  *otto.Otto                  //JS规则识别容器
	mainRule   *Application
}

func NewSchedule(cc *common.Communication) *Schedule {
	cc.Token = helper.NewToken(cc.UserId, cc.AppId, cc.DebugId).Crawler().ToString()
	return &Schedule{
		inData:    cc,
		outLog:    make(chan []byte),
		outData:   make(chan map[string]interface{}),
		mainRule:  NewApplication(),
		container: otto.New(),
		logQueue: &queue.Channel{
			Exchange:     "Crawlers",
			ExchangeType: "direct",
			RoutingKey:   common.PREFIX_LOG + cc.Token,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		dataQueue: &queue.Channel{
			Exchange:     "Crawlers",
			ExchangeType: "direct",
			RoutingKey:   common.PREFIX_DATA + cc.Token,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		rabbitConn: queue.RabbitConn,
	}
}

func (this *Schedule) Run() {
	defer func() {
		this.container.Call(FUNC_BEFORE_EXIT, nil)
		sp := spider_main.NewCrawler()
		err := sp.ModifyStatus(this.inData.AppId, spider_main.CRAWLER_STATUS_NORMAL)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		isDebug := this.inData.Method == common.SCHEDULE_METHOD_DEBUG
		err = this.pushLog(helper.FmtLog(common.LOG_INFO, "执行完成", common.LOG_LEVEL_INFO, common.LOG_TYPE_FINISH), isDebug)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		close(this.outLog)
		close(this.outData)
	}()
	err := this.container.Set("Crawler", this.init)
	if err != nil {
		this.outLog <- helper.FmtLog(common.LOG_ERROR, "初始化失败", common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		return
	}

}

func (this *Schedule) init(call otto.FunctionCall) otto.Value {
	if call.Argument(0).IsObject() {
		config := call.Argument(0).Object()
		this.mainRule.LazyLoad(config)
		obj, _ := this.container.Object(`Crawler = {}`)
		obj.Set("start", this.start)
		res, _ := this.container.ToValue(obj)
		return res
	}
	return otto.Value{}
}

func (this *Schedule) start(call otto.FunctionCall) otto.Value {
	token := helper.NewToken(this.inData.UserId, this.inData.AppId, this.inData.DebugId).Pool().ToString()
	if Spool.Exist(token) {
		this.outLog <- helper.FmtLog(common.LOG_INFO, "任务正在执行中...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
		return otto.Value{}
	}
	sc := spider_main.NewCrawler()
	this.bean, _ = sc.GetRowByID(this.inData.AppId)
}

func (this *Schedule) pushLog(body []byte, debug bool) error {
	if debug {
		//存入队列
		err := this.rabbitConn.Publish(this.logQueue, body)
		if err != nil {
			return err
		}
		return nil
	}
	//消息持久化
	res := make(map[string]interface{})
	err := json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	res["app_id"] = this.inData.AppId.String()
	table := fmt.Sprintf("%s%s", common.PREFIX_LOG, this.inData.Token)
	obj := spider_data.NewCrawlerLog(table)
	if _, err := obj.Build(res); err != nil {
		return err
	}
	return nil
}

func (this *Schedule) pushData(body map[string]interface{}, debug bool) error {
	if debug {
		data := make(map[string]interface{})
		for key, val := range body {
			for _, v := range val.(map[bool]interface{}) {
				data[key] = v
			}
		}
		res, _ := json.MarshalIndent(&data, "", "\t")
		err := this.rabbitConn.Publish(this.dataQueue, res)
		if err != nil {
			return err
		}
		return nil
	}
	//数据持久化
	data := make(map[string]interface{})
	for field, value := range body {
		data[field] = value.(map[bool]interface{})
	}
	data["app_id"] = map[bool]interface{}{
		false: this.inData.AppId,
	}
	data["user_id"] = map[bool]interface{}{
		false: this.inData.UserId,
	}
	data["created_at"] = map[bool]interface{}{
		false: time.Now().Unix(),
	}
	table := fmt.Sprintf("%s%s", common.PREFIX_DATA, this.inData.Token)
	obj := spider_data.NewCrawlerData(table)
	if err := obj.Build(data, this.bean.Method); err != nil {
		return err
	}
	return nil
}
