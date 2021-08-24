package spider

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler/rules"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/robertkrimen/otto"
	"net/http"
	"sync"
	"time"
)

type Spider struct {
	appId      int
	method     int
	container  *otto.Otto
	log        chan<- []byte
	data       chan<- map[string]interface{}
	abort      bool
	runTimes   int
	startedAt  int64
	mu         sync.RWMutex
	maxLimit   int
	rules      map[string]interface{}
	queue      *queue.Queue
	failure    map[string]int
	httpClient *http.Client
	token      string
}

func NewSpider(appId int, method int, rule map[string]interface{}, token string, vm *otto.Otto, lc chan<- []byte, dc chan<- map[string]interface{}) *Spider {
	return &Spider{
		appId:     appId,
		container: vm,
		log:       lc,
		data:      dc,
		token:     token,
		method:    method,
		startedAt: time.Now().Unix(),
		rules:     rule,
		failure:   make(map[string]int),
	}
}

//蜘蛛开始行动
func (this *Spider) Run() {
	defer func() {
		p := recover()
		if p != nil {
			this.log <- helper.FmtLog(common.LOG_ERROR, p.(error).Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}

		//@todo 更新状态
		this.log <- helper.FmtLog(common.LOG_INFO, "执行完成", common.LOG_LEVEL_INFO, common.LOG_TYPE_FINISH)
	}()

	this.log <- helper.FmtLog(common.LOG_INFO, "开始执行任务...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)

	//初始化数据结构
	if this.method == common.SCHEDULE_METHOD_EXECUTE {
		this.initTable()
	}
}

func (this *Spider) initTable() {
	var items []interface{}
	for _, field := range this.rules[rules.FIELDS].([]rules.FieldStash) {
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

	}

}
