package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_data"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"github.com/robertkrimen/otto"
	"sync"
	"time"
)

type Schedule struct {
	inData     *common.Communication       //传输过来的数据
	outLog     chan []byte                 //日志输出
	outData    chan map[string]interface{} //数据输出
	rabbitConn *queue.Base                 //rabbit链接
	logQueue   *queue.Channel              //日志队列
	dataQueue  *queue.Channel              //数据队列
	bean       *collect.Application        //数据持久化实例
	container  *otto.Otto                  //JS规则识别容器
	mainRule   *Application
}

func NewSchedule(cc common.Communication) *Schedule {
	cc.Token = helper.NewToken(cc.UserId, cc.AppId, cc.DebugId).Crawler().ToString()
	dataKey := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, cc.Token)
	logKey := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_LOG, cc.Token)
	return &Schedule{
		inData:    &cc,
		outLog:    make(chan []byte),
		outData:   make(chan map[string]interface{}),
		mainRule:  NewApplication(),
		container: otto.New(),
		logQueue: &queue.Channel{
			Exchange:     "Crawlers",
			ExchangeType: "direct",
			RoutingKey:   logKey,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		dataQueue: &queue.Channel{
			Exchange:     "Crawlers",
			ExchangeType: "direct",
			RoutingKey:   dataKey,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		rabbitConn: queue.RabbitConn,
	}
}

func (this *Schedule) Run() {
	defer func() {
		sp := collect.NewApplication()
		err := sp.ModifyStatus(this.inData.AppId, collect.STATUS_NORMAL)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		err = sp.ModifyToken(this.inData.AppId, this.inData.Token)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}()
	err := this.mainRule.Container.Set("Crawler", this.init)
	if err != nil {
		this.outLog <- common.FmtLog(common.LOG_ERROR, "初始化失败", common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		return
	}
	//初始化 js中的console.log
	console := map[string]interface{}{
		"log": func(call otto.FunctionCall) otto.Value {
			res := helper.FmtConsole(call.ArgumentList)
			this.outLog <- common.FmtLog(common.LOG_DEBUG, res, common.LOG_LEVEL_DEBUG, common.LOG_TYPE_SYSTEM)
			return otto.Value{}
		},
	}
	err = this.mainRule.Container.Set("console", console)
	if err != nil {
		spiderhub.Logger.Error("%v", err.Error())
		return
	}
	go func() {
		err := this.mainRule.InitBody(this.inData.Content)
		if err != nil {
			this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
	}()

	for {
		select {
		case m := <-this.outLog:
			//已经结结束
			var res common.LogLevel
			err := json.Unmarshal(m, &res)
			if err != nil {
				goto Loop
			}
			debug := this.inData.Method == common.SCHEDULE_METHOD_DEBUG
			err = this.pushLog(m, debug)
			if err != nil {
				spiderhub.Logger.Error("%v", err.Error())
			}
			if res.Type == 5 {
				goto Loop
			}
		case d := <-this.outData:
			debug := this.inData.Method == common.SCHEDULE_METHOD_DEBUG
			err = this.pushData(d, debug)
			if err != nil {
				spiderhub.Logger.Error("%v", err.Error())
			}
		}
	}
Loop:
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
		this.outLog <- common.FmtLog(common.LOG_INFO, "任务正在执行中...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
		return otto.Value{}
	}
	sc := collect.NewApplication()
	this.bean, _ = sc.GetRowByID(this.inData.AppId)
	if this.bean.Id == 0 {
		this.outLog <- common.FmtLog(common.LOG_ERROR, "未找到应用数据,请先创建爬虫应用", common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		return otto.Value{}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			Spool.Delete(token)
			wg.Done()
		}()
		//非调试模式且数据存储方式是重新
		if this.inData.Method == common.SCHEDULE_METHOD_EXECUTE && this.bean.Method == collect.METHOD_INSERT {
			dataTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, this.inData.Token)
			dataObj := spiderhub_data.NewCollectData(dataTable)
			err := dataObj.RemoveRows()
			if err != nil {
				this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
			}
			//清空日志
			logTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_LOG, this.inData.Token)
			logObj := spiderhub_data.NewCollectLog(logTable)
			err = logObj.RemoveRows()
			if err != nil {
				this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
			}
		}
		//初始化蜘蛛
		spd := NewSpider(this.inData.AppId, this.inData.Method, this.mainRule.Rules, this.inData.Token, this.mainRule.Container, this.outLog, this.outData)
		Spool.Start(token, spd)
	}()
	wg.Wait()
	return otto.Value{}
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
	res["app_id"] = this.inData.AppId
	logDoc := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_LOG, this.inData.Token)
	obj := spiderhub_data.NewCollectLog(logDoc)
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
	dataDoc := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, this.inData.Token)
	obj := spiderhub_data.NewCollectData(dataDoc)
	if err := obj.Build(data, this.bean.Method); err != nil {
		return err
	}
	return nil
}
