package cleaner

import (
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spider_main"
	"github.com/robertkrimen/otto"
)

const (
	//设置前缀
	PREFIX_LOG  = "cleanerLog"
	PREFIX_DATA = "cleanerData"
)

type Schedule struct {
	inData     *common.Communication
	bean       *spider_main.Crawler
	outLog     chan []byte
	outData    chan map[string]interface{}
	rabbitConn *queue.Base
	logQueue   *queue.Channel
	dataQueue  *queue.Channel
	container  *otto.Otto
	mainRule   *Application
}

func NewSchedule(cc common.Communication) *Schedule {
	cc.Token = helper.NewToken(cc.UserId, cc.AppId, cc.DebugId).Clean().ToString()
	return &Schedule{
		inData:  &cc,
		outLog:  make(chan []byte),
		outData: make(chan map[string]interface{}),
		logQueue: &queue.Channel{
			Exchange:     "Cleaners",
			ExchangeType: "direct",
			RoutingKey:   PREFIX_LOG + cc.Token,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		dataQueue: &queue.Channel{
			Exchange:     "Cleaners",
			ExchangeType: "direct",
			RoutingKey:   PREFIX_DATA + cc.Token,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		mainRule:   NewApplication(),
		container:  otto.New(),
		rabbitConn: queue.RabbitConn,
	}
}

func (this *Schedule) Run() {
	defer func() {
		sp := spider_main.NewCrawler()
		err := sp.ModifyStatus(this.inData.AppId, spider_main.STATUS_NORMAL)
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
	//初始化清洗对像
	this.mainRule.Container.Set("Clean", this.init)
}

//清洗初始化
func (this *Schedule) init(call otto.FunctionCall) otto.Value {
	if call.Argument(0).IsObject() {
		config := call.Argument(0).Object()
		this.mainRule.LazyLoad(config)
		obj, _ := this.container.Object(`Clean = {}`)
		err := obj.Set("start", this.start)
		if err != nil {
			spiderhub.Logger.Error("%s", err.Error())
		}
		res, _ := this.container.ToValue(obj)
		return res
	}
	return otto.Value{}
}

//开始清洗
func (this *Schedule) start(call otto.FunctionCall) otto.Value {

}
