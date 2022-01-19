package cleaner

import (
	"encoding/json"
	"fmt"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_data"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"github.com/robertkrimen/otto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type Cleaner struct {
	appId     int64
	rules     map[string]interface{}
	inst      collect.ApplicationImp
	method    int
	token     string
	container *otto.Otto
	outLog    chan<- []byte
	outData   chan<- map[string]interface{}
	abort     bool
	runTimes  int
	startBy   int64
	dataTable string
	mu        sync.RWMutex
}

func NewCleaner(appId int64, token string, crawlerToken string, method int, rule map[string]interface{}, vm *otto.Otto, log chan<- []byte, data chan<- map[string]interface{}) *Cleaner {
	dataTable := fmt.Sprintf("%s%s", common.PREFIX_CRAWL_DATA, crawlerToken)
	return &Cleaner{
		appId:     appId,
		inst:      collect.NewApplication(),
		token:     token,
		rules:     rule,
		container: vm,
		outLog:    log,
		outData:   data,
		method:    method,
		dataTable: dataTable,
	}
}

func (this *Cleaner) Run() {
	defer func() {
		p := recover()
		if p != nil {
			this.outLog <- common.FmtLog(common.LOG_ERROR, p.(error).Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
		err := this.inst.ModifyStatus(this.appId, common.STATUS_NORMAL)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		this.outLog <- common.FmtLog(common.LOG_INFO, "执行完成", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
		this.outLog <- common.FmtLog(common.LOG_INFO, "", common.LOG_LEVEL_INFO, common.LOG_TYPE_FINISH)
	}()
	err := this.inst.ModifyStatus(this.appId, common.STATUS_RUNNING)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	this.outLog <- common.FmtLog(common.LOG_INFO, "开始执行任务...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	//初始化数据结构
	if this.method == common.SCHEDULE_METHOD_EXECUTE {
		this.initTable()
	}
	var page int64 = 1
	var limit int64 = 50
	var skip int64
	lost := make(map[string]bool)
	cd := spiderhub_data.NewCollectData(this.dataTable)
	for {
		//中断执行
		if this.abort == true {
			goto Loop
		}
		skip = (page - 1) * limit
		//获取爬虫数据
		list, err := cd.GetRowsByPage(skip, limit)
		if err != nil {
			this.outLog <- common.FmtLog(common.LOG_ERROR, "获取分页数据错误:"+err.Error(), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
			continue
		}
		if len(list) == 0 {
			goto Loop
		}
		for _, item := range list {
			this.process(item)
		}
		page++
	}
Loop:
	if len(lost) > 0 {
		lostStr, _ := json.Marshal(lost)
		this.outLog <- common.FmtLog(common.LOG_ERROR, string(lostStr), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
	}
}

func (this *Cleaner) Stop() {
	this.abort = true
}

//处理数据
func (this *Cleaner) process(data interface{}) {
	mu.Lock()
	defer mu.Unlock()
	row := make(map[string]interface{})
	row["extract"] = map[string]interface{}{
		"__id":   data.(map[string]interface{})["_id"].(primitive.ObjectID).String(),
		"__url":  data.(map[string]interface{})["target_url"],
		"__time": data.(map[string]interface{})["created_at"],
	}
	delete(data.(map[string]interface{}), "_id")
	delete(data.(map[string]interface{}), "app_id")
	delete(data.(map[string]interface{}), "user_id")
	delete(data.(map[string]interface{}), "created_at")
	row["data"] = data.(map[string]interface{})

	//回调
	var result map[string]interface{}
	if res, err := this.container.Call(FUNC_ON_EACH_ROW, nil, row); err == nil {
		if res.IsDefined() == true {
			result = NewExtract(res, this.rules[FIELDS].([]FieldStash), this.container, this.outLog).Run()
		} else {
			result = this.packaging(row["data"].(map[string]interface{}), this.rules[FIELDS].([]FieldStash))
		}
	} else {
		result = this.packaging(row["data"].(map[string]interface{}), this.rules[FIELDS].([]FieldStash))
	}
	//下载附件
	if len(result) > 0 {
		_, primaryValue := this.searchPrimary(result)
		if primaryValue == nil {
			this.outLog <- common.FmtLog(common.LOG_WARNING, "下载附件需要指定主键", common.LOG_LEVEL_WARN, common.LOG_TYPE_SYSTEM)
		} else {
			for _, field := range this.rules[FIELDS].([]FieldStash) {
				if field.Download {
					err := NewDownload(result[field.Name].(map[bool]*common.FieldData)[field.Primary], primaryValue, field, this.token, this.container).Run()
					if err != nil {
						this.outLog <- common.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
					}
				}
			}
		}
	}
	result["target_url"] = map[bool]*common.FieldData{
		false: {
			Type:  TYPE_STRING,
			Value: data.(map[string]interface{})["target_url"],
			Alias: "目标地址",
		}}
	//data.(map[string]interface{})["target_url"]
	this.outData <- result
}

//打包数据
func (this *Cleaner) packaging(data map[string]interface{}, fields []FieldStash) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		tp := field.Type
		if len(tp) == 0 {
			tp = "string"
		}
		if _, ok := data[field.Name]; ok {
			result[field.Name] = map[bool]*common.FieldData{
				field.Primary: {
					Type:  tp,
					Value: data[field.Name],
					Alias: field.Alias,
				}}
		}
	}
	return result
}

//获取主键对
func (this *Cleaner) searchPrimary(result map[string]interface{}) (string, interface{}) {
	for field, item := range result {
		for isPrimary, data := range item.(map[bool]*common.FieldData) {
			if isPrimary {
				val := data.Value
				return field, val
			}
		}
	}
	return "", nil
}

func (this *Cleaner) initTable() {
	var items []interface{}
	for _, field := range this.rules[FIELDS].([]FieldStash) {
		table := make(map[string]string)
		alias := field.Alias
		if len(alias) == 0 {
			alias = field.Name
		}
		table = map[string]string{
			"name":  field.Name,
			"alias": alias,
		}
		items = append(items, table)
	}
	items = append(items, map[string]string{
		"name":  "target_url",
		"alias": "目标地址",
	})
	if len(items) > 0 {
		itemStr, err := json.Marshal(items)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		th := collect.NewAppField()
		err = th.Modify(common.TARGET_TYPE_CLEAN, this.appId, itemStr)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}

}
