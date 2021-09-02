package spider

import (
	"crypto/tls"
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spider_main"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"github.com/robertkrimen/otto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Spider struct {
	appId      primitive.ObjectID
	method     int
	container  *otto.Otto
	outLog     chan<- []byte
	outData    chan<- map[string]interface{}
	abort      bool
	runTimes   int
	tm         int64
	mu         sync.RWMutex
	maxLimit   int
	params     map[string]interface{}
	queue      *queue.Queue
	failure    map[string]int
	httpClient *http.Client
	token      string
	inst       *spider_main.CrawlerImpl
}

func NewSpider(appId primitive.ObjectID, method int, rule map[string]interface{}, token string, vm *otto.Otto, lc chan<- []byte, dc chan<- map[string]interface{}) *Spider {
	return &Spider{
		appId:     appId,
		container: vm,
		outLog:    lc,
		outData:   dc,
		token:     token,
		method:    method,
		tm:        time.Now().Unix(),
		params:    rule,
		failure:   make(map[string]int),
		inst:      spider_main.NewCrawler(),
	}
}

//蜘蛛开始行动
func (this *Spider) Run() {
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
	sp := colly.NewCollector()
	ts := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sp.WithTransport(ts)
	//不可重复抓取
	sp.AllowURLRevisit = false
	var timeout int64 = 30
	if _, ok := this.params[crawler.TIMEOUT]; ok {
		timeout = this.params[crawler.TIMEOUT].(int64)
	}
	if _, ok := this.params[crawler.MAX_LIMIT]; ok {
		this.maxLimit = this.params[crawler.MAX_LIMIT].(int)
	}
	//请求超时
	timeouts := time.Duration(timeout) * time.Second
	sp.SetRequestTimeout(timeouts)

	//限速
	var delay int64 = 1
	if _, ok := this.params[crawler.DELAY]; ok {
		delay = this.params[crawler.DELAY].(int64)
	}
	delays := time.Duration(delay) * time.Second
	err = sp.Limit(
		&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 2,
			Delay:       delays,
			RandomDelay: time.Second * 5,
		})
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}

	this.queue, _ = queue.New(2, nil)

	this.container.Call(crawler.FUNC_INIT_CRAWL, nil, this.queue)
	//随机UA
	extensions.RandomUserAgent(sp)

	//请求
	sp.OnRequest(this.onRequest)
	//错误
	sp.OnError(this.onError)
	//响应
	sp.OnResponse(this.onResponse)
	//完成
	sp.OnScraped(this.onScraped)

	//加载入口
	if _, ok := this.params[crawler.SCAN_URLS]; ok {
		for _, u := range this.params[crawler.SCAN_URLS].([]interface{}) {
			this.outLog <- helper.FmtLog(common.LOG_INFO, u.(string), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
			uf, _ := url.QueryUnescape(u.(string))
			err := this.queue.AddURL(uf)
			if err != nil {
				this.outLog <- helper.FmtLog(common.LOG_ERROR, u.(string), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
			}
		}
	}
	err = this.queue.Run(sp)
	if err != nil {
		this.outLog <- helper.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
	}
}

func (this *Spider) Stop() {
	this.outLog <- helper.FmtLog(common.LOG_INFO, "爬虫停止中...", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	this.abort = true
}

func (this *Spider) Finish() {
	this.outLog <- helper.FmtLog(common.LOG_INFO, "爬虫已停止运行", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	//this.bean.ModifyStatus(this.appId,spider_main.CRAWLER_STATUS_NORMAL)
}

func (this *Spider) initTable() {
	var items []interface{}
	for _, field := range this.params[crawler.FIELDS].([]crawler.FieldStash) {
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
		err = th.Modify(common.TARGET_TYPE_CRAWLER, this.appId, itemStr)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}

}
