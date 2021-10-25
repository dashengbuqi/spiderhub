package cleaner

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
	inData       *common.Communication
	bean         *collect.Application
	outLog       chan []byte
	outData      chan map[string]interface{}
	rabbitConn   *queue.Base
	logQueue     *queue.Channel
	dataQueue    *queue.Channel
	container    *otto.Otto
	mainRule     *RuleConfig
	crawlerTable string //爬虫数据表
	logTable     string //清洗数据日志表
	dataTable    string //清洗数据表
}

func NewSchedule(cc common.Communication) *Schedule {
	cc.Token = helper.NewToken(cc.UserId, cc.AppId, cc.DebugId).Clean().ToString()
	lt := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_LOG, cc.Token)
	dt := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_DATA, cc.Token)
	return &Schedule{
		inData:    &cc,
		logTable:  lt,
		dataTable: dt,
		outLog:    make(chan []byte),
		outData:   make(chan map[string]interface{}),
		logQueue: &queue.Channel{
			Exchange:     "Cleaners",
			ExchangeType: "direct",
			RoutingKey:   lt,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		dataQueue: &queue.Channel{
			Exchange:     "Cleaners",
			ExchangeType: "direct",
			RoutingKey:   dt,
			Reliable:     true,
			Durable:      false,
			AutoDelete:   true,
		},
		mainRule:   NewRuleConfig(),
		container:  otto.New(),
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
		err = sp.ModifyCleanToken(this.inData.AppId, this.inData.Token)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}()
	var err error
	//初始化清洗对像
	err = this.mainRule.Container.Set("Clean", this.init)
	if err != nil {
		this.outLog <- common.FmtLog(common.LOG_ERROR, "初始化失败", common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		return
	}
	console := map[string]interface{}{
		"log": func(call otto.FunctionCall) otto.Value {
			out := helper.FmtConsole(call.ArgumentList)
			this.outLog <- common.FmtLog(common.LOG_DEBUG, out, common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
			return otto.Value{}
		},
	}
	err = this.mainRule.Container.Set("console", console)
	if err != nil {
		this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		return
	}
	go func() {
		err = this.mainRule.InitBody(this.inData.Content)
		if err != nil {
			this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
	}()

	for {
		select {
		case l := <-this.outLog:
			//已经结结束
			var res common.LogLevel
			err := json.Unmarshal(l, &res)
			if err != nil {
				goto Loop
			}
			debug := this.inData.Method == common.SCHEDULE_METHOD_DEBUG
			err = this.pushLogger(l, debug)
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
	key := helper.NewToken(this.inData.UserId, this.inData.AppId, this.inData.DebugId).Pool().ToString()
	if CleanPool.Exist(key) {
		this.outLog <- common.FmtLog(common.LOG_INFO, "任务正在执行中...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
		return otto.Value{}
	}
	sp := collect.NewApplication()
	this.bean, _ = sp.GetRowByID(this.inData.AppId)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			CleanPool.Delete(key)
			wg.Done()
		}()
		err := sp.ModifyStatus(this.inData.AppId, collect.STATUS_RUNNING)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		if this.inData.Method == common.SCHEDULE_METHOD_EXECUTE && this.bean.Method == collect.METHOD_INSERT {
			dataObj := spiderhub_data.NewCollectData(this.dataTable)
			err := dataObj.RemoveRows()
			if err != nil {
				this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
			}
			//清空日志
			logObj := spiderhub_data.NewCollectLog(this.logTable)
			err = logObj.RemoveRows()
			if err != nil {
				this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
			}
		}
		cn := NewCleaner(this.inData.AppId, this.inData.Token, this.bean.CrawlerToken, this.inData.Method, this.mainRule.Rules, this.mainRule.Container, this.outLog, this.outData)
		CleanPool.Start(key, cn)
	}()

	wg.Wait()
	return otto.Value{}
}

func (this *Schedule) pushLogger(data []byte, debug bool) error {
	if debug {
		err := this.rabbitConn.Publish(this.logQueue, data)
		return err
	}
	//消息持久化
	res := make(map[string]interface{})
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	res["app_id"] = this.inData.AppId
	obj := spiderhub_data.NewCollectLog(this.logTable)
	if _, err := obj.Build(res); err != nil {
		return err
	}
	return nil
}

func (this *Schedule) pushData(body map[string]interface{}, debug bool) error {
	if debug {
		data := make(map[string]interface{})
		for key, val := range body {
			for _, v := range val.(map[bool]*common.FieldData) {
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
	data := make(map[string]map[bool]*common.FieldData)
	for field, value := range body {
		data[field] = value.(map[bool]*common.FieldData)
	}
	data["app_id"] = map[bool]*common.FieldData{
		false: {
			Alias: "应用",
			Type:  TYPE_INT,
			Value: this.inData.AppId,
		},
	}
	data["user_id"] = map[bool]*common.FieldData{
		false: {
			Alias: "用户",
			Type:  TYPE_INT,
			Value: this.inData.UserId,
		},
	}
	data["created_at"] = map[bool]*common.FieldData{
		false: {
			Alias: "创建时间",
			Type:  TYPE_INT,
			Value: time.Now().Unix(),
		},
	}
	obj := spiderhub_data.NewCollectData(this.dataTable)
	if err := obj.Build(data, this.bean.Method); err != nil {
		return err
	}
	return nil
}
