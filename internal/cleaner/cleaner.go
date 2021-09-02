package cleaner

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spider_main"
	"github.com/robertkrimen/otto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cleaner struct {
	appId     primitive.ObjectID
	rules     map[string]interface{}
	inst      *spider_main.CrawlerImpl
	method    int
	token     string
	container *otto.Otto
	outLog    chan<- []byte
	outData   chan<- map[string]interface{}
	abort     bool
	runTimes  int
	startBy   int64
}

func NewCleaner(appId primitive.ObjectID, token string, method int, rule map[string]interface{}, vm *otto.Otto, log chan<- []byte, data chan<- map[string]interface{}) *Cleaner {
	return &Cleaner{
		appId:     appId,
		inst:      spider_main.NewCrawler(),
		token:     token,
		rules:     rule,
		container: vm,
		outLog:    log,
		outData:   data,
		method:    method,
	}
}

func (this *Cleaner) Run() {
	defer func() {
		p := recover()
		if p != nil {
			this.outLog <- helper.FmtLog(common.LOG_ERROR, p.(error).Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
		err := this.inst.ModifyStatus(this.appId, spider_main.STATUS_NORMAL)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
		this.outLog <- helper.FmtLog(common.LOG_INFO, "执行完成", common.LOG_LEVEL_INFO, common.LOG_TYPE_FINISH)
	}()
	err := this.inst.ModifyStatus(this.appId, spider_main.STATUS_RUNNING)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	this.outLog <- helper.FmtLog(common.LOG_INFO, "开始执行任务...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	//初始化数据结构
	if this.method == common.SCHEDULE_METHOD_EXECUTE {
		this.initTable()
	}
	var page int64 = 1
	var limit int64 = 20
	var skip int64
	lost := make(map[string]bool)
	for {
		//中断执行
		if this.abort == true {
			goto Loop
		}
		skip = (page - 1) * limit

	}
Loop:
	if len(lost) == 0 {
		lostStr, _ := json.Marshal(lost)
		this.outLog <- helper.FmtLog(common.LOG_ERROR, string(lostStr), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
	}
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
		th := spider_main.NewTableHead()
		err = th.Modify(common.TARGET_TYPE_CLEAN, this.appId, itemStr)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}

}
