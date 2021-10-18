package cleaner

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_data"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"github.com/robertkrimen/otto"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

func NewCleaner(appId int64, token string, dataTable string, method int, rule map[string]interface{}, vm *otto.Otto, log chan<- []byte, data chan<- map[string]interface{}) *Cleaner {
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
		err := this.inst.ModifyStatus(this.appId, collect.STATUS_NORMAL)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		this.outLog <- nil
	}()
	err := this.inst.ModifyStatus(this.appId, collect.STATUS_RUNNING)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	this.outLog <- common.FmtLog(common.LOG_INFO, "开始执行任务...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	//初始化数据结构
	if this.method == common.SCHEDULE_METHOD_EXECUTE {
		this.initTable()
	}
	var page int64 = 1
	var limit int64 = 20
	var skip int64
	lost := make(map[string]bool)
	cd := spiderhub_data.NewCrawlerData(this.dataTable)
	for {
		//中断执行
		if this.abort == true {
			goto Loop
		}
		skip = (page - 1) * limit
		//获取爬虫数据
		list, err := cd.GetRowsBy(skip, limit)
		if err != nil {
			this.outLog <- common.FmtLog(common.LOG_ERROR, "获取分页数据错误:"+err.Error(), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
			continue
		}
		if list == nil {
			goto Loop
		}
		for _, item := range list {
			if this.abort == true {
				goto Loop
			}
			this.process(item.(map[string]interface{}))
		}
		page++
	}
Loop:
	if len(lost) == 0 {
		lostStr, _ := json.Marshal(lost)
		this.outLog <- common.FmtLog(common.LOG_ERROR, string(lostStr), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
	}
	this.outLog <- common.FmtLog(common.LOG_INFO, "数据清洗结束", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
}

func (this *Cleaner) Stop() {
	this.abort = true
}

//处理数据
func (this *Cleaner) process(data map[string]interface{}) {
	row := make(map[string]interface{})
	row["extract"] = map[string]interface{}{
		"__id":   data["_id"].(primitive.ObjectID).String(),
		"__url":  data["targetUrl"],
		"__time": data["created_at"],
	}
	delete(data, "_id")
	delete(data, "app_id")
	delete(data, "user_id")
	delete(data, "created_at")
	row["data"] = data

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
	this.outData <- result
}

//打包数据
func (this *Cleaner) packaging(data map[string]interface{}, fields []FieldStash) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		if _, ok := data[field.Name]; ok {
			result[field.Name] = map[bool]*common.FieldData{
				field.Primary: {
					Type:  field.Type,
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
		for isPrimary, data := range item.(map[bool]interface{}) {
			if isPrimary {
				val := data.(*common.FieldData).Value
				return field, val
			}
		}
	}
	return "", nil
}

func (this *Cleaner) initTable() {
	var items []interface{}
	for _, field := range this.rules[crawler.FIELDS].([]crawler.FieldStash) {
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
