package crawler

import (
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/dashengbuqi/spiderhub"
	"github.com/oliveagle/jsonpath"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

const (
	SELECTORTYPE_XPATH int64 = iota
	SELECTORTYPE_JSONPATH
	SELECTORTYPE_REGEX
)

const (
	APP_NAME      = "name"
	MAX_LIMIT     = "max_limit"
	DOMAIN        = "domains"
	SCAN_URLS     = "scan_urls"
	CONTENT_REGEX = "content_regex"
	HELPER_REGEX  = "helper_regex"
	FIELDS        = "fields"
	TIMEOUT       = "timeout"
	COOKIE        = "cookie"
	//CONFIG_JS_ENGINE        = "js_engine"
	//CONFIG_ENABLEJS         = "enable_js"
	USERAGENT          = "user_agent"
	ACCEPT_HTTP_STATUS = "accept_http_status"
	AUTOFIND_URLS      = "auto_find_urls"
	DELAY              = "delay"
	//CONFIG_VISIT_TYPE       = "visit_type"

	//function initCrawl(site)
	//@param site 内置对象，参考内置对象site
	//爬虫初始化的时候被调用，多节点运行时，只在第一个节点（又称主节点）中被调用，其他节点等待主节点的initCrawl方法执行完之后才继续往下执行。
	//建议在此回调中做添加入口页的操作。
	FUNC_INIT_CRAWL = "initCrawl"
	//function beforeCrawl(site)
	//@param site 内置对象，参考内置对象site
	//initCrawl方法之后被调用，在所有节点中都会被调用。
	//全局的User-Agent设置、Cookie设置建议放到此回调函数中。
	FUNC_BEFORE_CRAWL = "beforeCrawl"
	//function onDataReceived(data, site)
	//@param data 数据对象，Pipeline中前一个应用的数据
	//@param site 内置对象，参考内置对象site
	//Pipeline中的爬虫通过此回调来获取前一个应用的数据，在beforeCrawl之后被回调。
	//是否需要在内容中获取新链接
	FUNC_ON_DATA_RECEIVED = "onDataReceived"
	//function afterCrawl(site)
	//@param site 内置对象，参考内置对象site
	//爬虫结束时调用，每个节点都会回调，在beforeExit之前被回调。
	FUNC_AFTER_CRAWL = "afterCrawl"
	//function beforeExit(site)
	//@param site 内置对象，参考内置对象site
	//爬虫结束时回调，只有最后一个结束的节点会回调此方法，在afterCrawl之后被回调。
	FUNC_BEFORE_EXIT = "beforeExit"
	//function beforeDownloadPage(page, site)
	//@param page 内置对象，参考内置对象page
	//@param site 内置对象，参考内置对象site
	//@return page 内置对象，参考内置对象page。不重写此函数时，默认返回原page对象
	//当链接调度器从待爬队列中调度出来一个链接的时候，回调此函数。在此回调函数中可以修改链接地址page.url，修改完之后需要return page。
	// 常见的场景是链接中有时间戳，而添加链接和处理链接的时间通常是不确定的，这时可以在此回调函数中更新链接中的时间戳。
	FUNC_BEFORE_DOWNLOAD_PAGE = "beforeDownloadPage"
	//function onChangeProxy(site, page)
	//@param site 内置对象，参考内置对象site
	//@param page 内置对象，参考内置对象page
	//当获取到一个新的代理的时候，回调此函数。切换代理之后，之前的cookie会被清空，一般在此回调中做一些cookie的加载。
	FUNC_ON_CHANGE_PROXY = "onChangeProxy"
	//function isAntiSpider(url, content, page)
	//@param url 当前正在处理的链接地址
	//@param content 当前下载的网页内容
	//@param page 内置对象，参考内置对象page
	//@return boolean 是否反爬，true表示反爬，false表示没有反爬。不重写此函数时，默认返回false。
	//每个被调度的链接下载完成之后，会先判断返回的状态码是否403，如果403，则直接会触发切换代理；
	// 如果不是403，则回调此函数，开发者一般需要在此函数中判断返回码或者网页内容，给出是否反爬的判断，如果判断为反爬，
	// 需要返回true，否则返回false。
	//configs.isAntiSpider = function(url, content, page) {
	//  if (page.raw && page.raw.indexOf("请求太快了，请休息一会") >= 0) {
	//    return true;
	//  }
	//  return false;
	//}
	FUNC_IS_ANTI_SPIDER = "isAntiSpider"
	//function afterDownloadPage(page, site)
	//@param page 内置对象，参考内置对象page
	//@param site 内置对象，参考内置对象site
	//@return page 内置对象，参考内置对象page。不重写此函数时，默认返回原page对象。
	//每个被调度的链接下载完成之后回调该函数。在该函数中可以修改page.url和page.raw，修改之后，修改之后的内容会一直持续到该链接的生命周期结束。
	// 修改page.raw后会影响后续的数据抽取，所以一般可以在这个回调函数中发一些请求，把获取的数据拼接到page.raw中，以便后续抽取。
	FUNC_AFTER_DOWNLOAD_PAGE = "afterDownloadPage"
	//function onProcessScanPage(page, content, site)
	//@param page 内置对象，参考内置对象page
	//@param content 网页内容，content与page.raw的区别在于，content中的链接都是绝对地址（以http开头）
	//@param site 内置对象，参考内置对象site
	//@return boolean 是否还需要自动发现链接，true表示还需要自动发现，false表示不需要自动发现。不重写此函数时，默认返回configs.autoFindUrls的值。
	//网页在下载完之后，如果当前链接是入口页，则回调此函数。一般在此函数中实现手动链接发现，一般是发现帮助页，也可以直接发现内容页。
	FUNC_ON_PROCESS_SCAN_PAGE = "onProcessScanPage"
	//function onProcessHelperPage(page, content, site)
	//@param page 内置对象，参考内置对象page
	//@param content 网页内容，content与page.raw的区别在于，content中的链接都是绝对地址（以http开头）
	//@param site 内置对象，参考内置对象site
	//@return boolean 是否还需要自动发现链接，true表示还需要自动发现，false表示不需要自动发现。不重写此函数时，默认返回configs.autoFindUrls的值。
	//入口页判断以及onProcessScanPage回调之后，会继续判断当前链接是否是帮助页，如果是，则回调此函数。
	// 一般在此函数中实现手动链接发现，多数情况是发现内容页链接以及下一页帮助页的链接。
	FUNC_ON_PROCESS_HELPER_PAGE = "onProcessHelperPage"
	//function onProcessContentPage(page, content, site)
	//@param page 内置对象，参考内置对象page
	//@param content 网页内容，content与page.raw的区别在于，content中的链接都是绝对地址（以http开头）
	//@param site 内置对象，参考内置对象site
	//@return boolean 是否还需要自动发现链接，true表示还需要自动发现，false表示不需要自动发现。不重写此函数时，默认返回configs.autoFindUrls的值。
	//帮助页判断以及onProcessHelperPage回调之后，会继续判断当前链接是否是内容页，如果是，则回调此函数。一般内容页不需要再做链接发现，所以此函数多数情况下直接返回false。
	//onProcessXxxPage小结：
	//这个系列的三个函数主要用来控制链接的发现，如果想要提高爬虫爬取效率，或者需要精确地控制爬虫的爬取路径，需要重点实现这三个函数，并禁用自动链接发现。
	//一个链接可能同时是入口页和帮助页，也可能同时是帮助页和内容页，甚至可能同时是入口页、帮助页和内容页，这种情况下，这个链接产生的onProcessXxxPage回调，
	// 必须同时都返回false，才能禁用自动链接发现。
	FUNC_ON_PROCESS_CONTENT_PAGE = "onProcessContentPage"
	//function afterDownloadAttachedPage(page, site)
	//@param page 内置对象，参考内置对象page
	//@param site 内置对象，参考内置对象site
	//@return page 内置对象，参考内置对象page。不重写此函数时，默认返回原page对象。
	//attachedUrl下载完成之后会回调此函数。可以在此函数中修改page.raw的值，从而影响attachedUrl的后续抽取。
	// 多数场景是，attachedUrl返回的数据是jsonp格式，这时需要在此回调中把数据处理成json数据，以便后续用JsonPath来抽取。
	FUNC_AFTER_DOWNLOAD_ATTACHED_PAGE = "afterDownloadAttachedPage"
	//function afterExtractField(fieldName, data, page, site, index)
	//@param fieldName 抽取项名
	//@param data 当前抽取项抽取出的数据
	//@param page 内置对象，参考内置对象page
	//@param site 内置对象，参考内置对象site
	//@param index 当前项是在父抽取项的第几个结果中进行抽取，从0开始。
	//@return 数据对象 返回此项对应的数据。当不重写此函数时，默认返回原data对象。
	//在每个抽取项抽取到内容时回调此函数，一个网页的抽取过程中，会多次回调此函数。在此函数中，可以对抽取到的数据做进一步的处理，然后返回处理后的数据。
	FUNC_AFTER_EXTRACT_FIELD = "afterExtractField"
	//function extractTemporaryField(url,content)
	//@return
	FUNC_AFTER_EXTRACT_TEMPORARY_FIELD = "afterExtractTemporaryField"
	//function beforeHandleImg(fieldName, img)
	//@param fieldName 抽取项名，同afterExtractField
	//@param img 一个完整的img标签
	//@return String 处理后的img
	//在抽取的内容中发现标签时，回调此函数。一般在此函数中修改src，使src指向真实的图片地址。
	FUNC_BEFORE_HANDLE_IMG = "beforeHandleImg"
	//function beforeHostFile(fieldName, url)
	//@param fieldName 抽取项名，同afterExtractField
	//@param url 待托管的文件链接
	//@return newUrl 处理后的托管链接
	//在托管文件之前回调此函数，在此函数中可以对文件地址做修改。常用的场景是，在图片托管中，修改链接地址来获取分辨率更高的图片。
	FUNC_BEFORE_HOST_FILE = "beforeHostFile"
	//function afterHostFile(fieldName, hostedUrl)
	//@param fieldName 抽取项名，同afterExtractField
	//@param hostedUrl 托管后的链接地址
	//在托管后的文件链接计算结束之后回调此函数，在此函数中可以对托管后的链接进行收集。
	FUNC_AFTER_HOST_FILE = "afterHostFile"
	//function afterExtractPage(page, data, site)
	//@param page 内置对象，参考内置对象page
	//@param data 整个页面抽取出的数据
	//@param site 内置对象，参考内置对象site
	//@return 数据对象 返回处理后的抽取数据。当不重写此函数时，默认返回原data对象。
	//当整个网页完成抽取时回调此函数。一般在此回调中做一些数据整理的操作，也可以继续发送网络请求，把返回的数据整理后放到data中返回。
	FUNC_AFTER_EXTRACT_PAGE = "afterExtractPage"

	VALUE_STRING = "string"
	VALUE_ARRAY  = "array"
	VALUE_MAP    = "map"
	VALUE_INT    = "int"

	BODY_TEXT = 0
	BODY_HTML = 1
)

var (
	ValueArr = []string{
		VALUE_STRING,
		VALUE_MAP,
		VALUE_ARRAY,
		VALUE_INT,
	}
)

const (
	EXTRACT_NORMAL      = iota //正常提取内容 Extract
	EXTRACT_URLCONTEXT         //发现内容链接
	EXTRACT_ATTACHEDURL        //内容在其它链接上
)

type FieldStash struct {
	//抽取项的名字。
	Name string `json:"name"`
	//抽取项的别名。一般起中文名，方便查看数据。只影响网页的上显示，可随意修改。
	Alias string `json:"alias"`
	//存储容器 string,array,map
	Container string `json:"container"`
	//标识当前抽取项的值是否必须（不能为空）。默认是false，可以为空。
	Required bool `json:"required"`
	//需要下载
	Download bool `json:"download"`
	//抽取规则的类型。默认值是SelectorType.XPath
	//SelectorType.XPath 一般针对html网页或xml，查看教程
	//SelectorType.JsonPath 针对json数据，查看教程
	//SelectorType.Regex 可以针对一切文本，查看教程
	SelectorType int `json:"selector_type"`
	//抽取规则 如果selector为空或者未设置，则抽取的值为null，在进行required的判定之前，仍会进行afterExtractField回调。
	Selector string `json:"selector"`
	//内容类型 0:text 1:html
	BodyType int `json:"body_type"`
	// 抽取项的子抽取项。
	//field支持子项，可以设置多层级，方便数据以本身的层级方式存储，而不用全部展开到第一层级。
	//注意：
	//第一层field默认从当前网页内容中抽取，而子项默认从父项的内容中抽取
	Children []FieldStash `json:"children"`
	//当前抽取项是否作为整条数据的主键组成部分。默认是false。
	Primary bool `json:"primary"`
	// 数据抽取源。无默认值，不设置时，抽取源默认时当前网页或父项内容。
	// 访问方式:0 UrlContext  1 AttachedUrl
	// SourceType.UrlContext
	// 表示从当前链接的附加内容中抽取数据。在添加链接的时候，可以同时给该链接附加一些数据，通常的使用场景是，列表页展示的一些内容页没有的数据，
	// 那么在做链接发现的时候，可以直接把这部分数据附加到对应的内容页链接上。
	// SourceType.AttachedUrl
	// 表示需要的数据在另外一个链接（我们叫attachedUrl）的请求结果里面，需要额外再发一次请求来获取数据。
	//只有当sourceType为SourceType.AttachedUrl时，下面的attachedUrl、attachedMethod、attachedParams、attachedHeaders设置才有意义。
	// 多说两句：在还没有site.requestUrl接口的时候，很多复杂的场景（比如文章内容分页）只能用attachedUrl来实现；
	// 相比于attachedUrl，site.requestUrl具有更高的灵活性，所以对于复杂的场景，我们建议在回调函数中通过site.requestUrl来实现，
	// attachedUrl只用来实现一些简单的单请求场景。
	ExtractMethod int `json:"extract_method"`
	// attachedUrl请求地址。
	//attachedUrl支持变量替换，变量可引用上下文中已经抽取的字段。
	//同一层级的field或者第一层级的field，引用方式为”{fieldName}”。
	//不同层级需要从根field开始，$表示根，引用方式为”{$.fieldName.fieldName}”。
	//特殊变量$$url表示当前网页的url，引用方式为”{$$url}”。
	//比如抽取到字段item_id，attachedUrl形式为https://item.example.com/item.htm?id=1000，则attachedUrl可以写为：
	//"https://item.example.com/item.htm?id={item_id}"
	AttachedUrl string `json:"attached_url"`
	// HTTP请求是”GET”还是”POST”。默认是”GET”。
	AttachedMethod string `json:"attached_method"`
	// HTTP请求的POST参数。如果请求是”GET”，参数将会被忽略。
	// 参数形如a=b&c=d，支持变量替换。与attachedUrl的变量引用方式相同。
	AttachedUrlParams string `json:"attached_url_params"`
	// HTTP请求的headers。
	AttachedHeaders map[string]string `json:"attached_headers"`
	// 抽取项是否是临时的。默认是false。临时的抽取项，数据存储的时候，不会存储其值。
	Temporary bool `json:"temporary"`
	// 回调函数
	Func string `json:"func"`
}

const (
	USERAGENT_COMPUTER int = iota
	USERAGENT_ANDROID
	USERAGENT_IOS
	USERAGENT_MOBILE
	USERAGENT_EMPTY

	HTTP_STATUS_SUCCESS               = 200
	HTTP_STATUS_MOVED_PERMANENTLY     = 301
	HTTP_STATUS_MOVE_TEMPORARILY      = 302
	HTTP_STATUS_BAD                   = 400
	HTTP_STATUS_UNAUTHORIZED          = 401
	HTTP_STATUS_PAYMENT_REQUIRED      = 402
	HTTP_STATUS_FORBIDDEN             = 403
	HTTP_STATUS_NOT_FOUND             = 404
	HTTP_STATUS_METHOD_NOT            = 405
	HTTP_STATUS_INTERNAL_SERVER_ERROR = 500
	HTTP_STATUS_BAD_GATEWAY           = 502
	HTTP_STATUS_SERVICE_UNAVAILABLE   = 503
	HTTP_STATUS_GATEWAY_TIMEOUT       = 504
)

type CrawlerConfig struct {
	Name             string           `json:"name"`
	Domains          []interface{}    `json:"domains"`
	MaxLimit         int              `json:"max_limit"`
	ScanUrls         []interface{}    `json:"scan_urls"`
	HelperRegex      []*regexp.Regexp `json:"helper_regex"`
	ContentRegex     []*regexp.Regexp `json:"content_regex"`
	Fields           []*FieldStash    `json:"fields"`
	Timeout          int64            `json:"timeout"`
	Delay            int64            `json:"delay"`
	UserAgent        int              `json:"user_agent"`
	AcceptHttpStatus []int            `json:"accept_http_status"`
	AutoFindUrls     bool             `json:"auto_find_urls"`
}

type Application struct {
	Rules     map[string]interface{}
	Container *otto.Otto
	RuleName  string
	oo        *otto.Object
	Start     bool
}

func NewApplication() *Application {
	return &Application{
		Rules:     make(map[string]interface{}),
		Container: otto.New(),
	}
}

func (this *Application) Init(body string) error {
	var err error
	SelectorType, _ := this.Container.Object(`({Xpath:0,JsonPath:1,Regex:2})`)
	err = this.Container.Set("SelectorType", SelectorType)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	UserAgent, _ := this.Container.Object(`({Computer:0,Android:1,IOS:2,Mobile:3,Empty:4,})`)
	err = this.Container.Set("UserAgent", UserAgent)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	SourceType, _ := this.Container.Object(`({Normal:0,UrlContext:1,AttachedUrl:2})`)
	err = this.Container.Set("ExtractType", SourceType)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	Typ, _ := this.Container.Object(`({String: "string",Array:"array",Map:"map"})`)
	err = this.Container.Set("Container", Typ)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	_, err = this.Container.Run(body)
	return err
}

func (this *Application) unpackBoolean(key string) {
	res, _ := this.oo.Get(key)
	if res.IsUndefined() {
		this.Rules[key] = false
		return
	}
	value, _ := res.ToBoolean()
	this.Rules[key] = value
}

func (this *Application) unpackInt(key string) {
	res, _ := this.oo.Get(key)
	if res.IsUndefined() {
		this.Rules[key] = 0
		return
	}
	value, _ := res.ToInteger()
	this.Rules[key] = value
}

func (this *Application) unpackArray(key string) {
	res, _ := this.oo.Get(key)
	if res.IsUndefined() {
		this.Rules[key] = []interface{}{}
		return
	}
	if res.IsObject() {
		vals, _ := res.Export()
		this.Rules[key] = vals.([]interface{})
	}
}

func (this *Application) unpackString(key string) {
	res, _ := this.oo.Get(key)
	if res.IsUndefined() {
		this.Rules[key] = ""
		return
	}
	value, _ := res.ToString()
	if len(value) > 0 && value != "undefined" {
		this.Rules[key] = value
	}
}

func (this *Application) setCallBack(fn string) {
	if fnVal, _ := this.oo.Get(fn); fnVal.IsDefined() {
		err := this.Container.Set(fn, fnVal)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}
}

func (this *Application) LazyLoad(oo *otto.Object) {
	var err error
	this.oo = oo

	this.unpackString(APP_NAME)
	this.unpackString(COOKIE)
	this.unpackArray(DOMAIN)
	this.unpackArray(SCAN_URLS)
	this.unpackArray(CONTENT_REGEX)
	this.unpackArray(HELPER_REGEX)
	this.unpackInt(TIMEOUT)
	this.unpackInt(USERAGENT)
	this.unpackArray(ACCEPT_HTTP_STATUS)
	this.unpackBoolean(AUTOFIND_URLS)
	this.unpackInt(DELAY)
	this.unpackInt(MAX_LIMIT)

	//回调方法
	this.setCallBack(FUNC_INIT_CRAWL)
	this.setCallBack(FUNC_BEFORE_CRAWL)
	this.setCallBack(FUNC_ON_DATA_RECEIVED)
	this.setCallBack(FUNC_AFTER_CRAWL)
	this.setCallBack(FUNC_BEFORE_EXIT)
	this.setCallBack(FUNC_BEFORE_DOWNLOAD_PAGE)
	this.setCallBack(FUNC_ON_CHANGE_PROXY)
	this.setCallBack(FUNC_IS_ANTI_SPIDER)
	this.setCallBack(FUNC_AFTER_DOWNLOAD_PAGE)
	this.setCallBack(FUNC_ON_PROCESS_SCAN_PAGE)
	this.setCallBack(FUNC_ON_PROCESS_HELPER_PAGE)
	this.setCallBack(FUNC_ON_PROCESS_CONTENT_PAGE)
	this.setCallBack(FUNC_AFTER_DOWNLOAD_ATTACHED_PAGE)
	this.setCallBack(FUNC_AFTER_EXTRACT_FIELD)
	this.setCallBack(FUNC_BEFORE_HANDLE_IMG)
	this.setCallBack(FUNC_BEFORE_HOST_FILE)
	this.setCallBack(FUNC_AFTER_HOST_FILE)
	this.setCallBack(FUNC_AFTER_EXTRACT_PAGE)
	this.setCallBack(FUNC_AFTER_EXTRACT_TEMPORARY_FIELD)

	//生成js方法
	err = this.Container.Set("extracts", this.extractList)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	err = this.Container.Set("extract", this.extract)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	//解析字段
	if valField, err := this.Container.Get(FIELDS); err == nil {
		fields := valField.Object().Keys()
		if len(fields) > 0 {
			var items []FieldStash
			for _, fd := range fields {
				item, err := valField.Object().Get(fd)
				if err != nil {
					spiderhub.Logger.Error("%v", err)
					continue
				}
				var fs FieldStash
				fieldItem, err := item.Export()
				if err != nil {
					spiderhub.Logger.Error("%v", err)
					continue
				}
				field, err := json.Marshal(fieldItem)
				if err != nil {
					spiderhub.Logger.Error("%v", err)
					continue
				}
				err = json.Unmarshal(field, &fs)
				if err != nil {
					spiderhub.Logger.Error("%v", err)
					continue
				}
				items = append(items, fs)
			}
			this.Rules[FIELDS] = items
		}
	}
	this.Start = true
}

func (this *Application) extractList(call otto.FunctionCall) (result otto.Value) {
	body := call.Argument(0).String()
	r := call.Argument(1).String()
	t, _ := call.Argument(2).ToInteger()

	if t == SELECTORTYPE_JSONPATH {
		var item interface{}
		err := json.Unmarshal([]byte(body), &item)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
			return
		}
		res, _ := jsonpath.JsonPathLookup(item, r)
		result, _ = this.Container.ToValue(res)
	} else if t == SELECTORTYPE_XPATH {
		doc, err := htmlquery.Parse(strings.NewReader(body))
		if err != nil {
			spiderhub.Logger.Error("%v", err)
			return
		}
		expr := xpath.MustCompile(r)
		iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
		var items []interface{}
		for iter.MoveNext() {
			if item := iter.Current().Value(); len(item) > 0 {
				items = append(items, item)
			}
		}
		result, _ = this.Container.ToValue(items)
	}
	return
}

func (this *Application) extract(call otto.FunctionCall) (result otto.Value) {
	body := call.Argument(0).String()
	r := call.Argument(1).String()
	t, _ := call.Argument(2).ToInteger()

	if t == SELECTORTYPE_JSONPATH {
		var item interface{}
		err := json.Unmarshal([]byte(body), &item)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
			return
		}
		res, _ := jsonpath.JsonPathLookup(item, r)
		result, _ = this.Container.ToValue(res)
	} else if t == SELECTORTYPE_XPATH {
		doc, err := htmlquery.Parse(strings.NewReader(body))
		if err != nil {
			spiderhub.Logger.Error("%v", err)
			return
		}
		expr := xpath.MustCompile(r)
		iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
		var item interface{}
		for iter.MoveNext() {
			if v := iter.Current().Value(); len(v) > 0 {
				item = v
				break
			}
		}
		result, _ = this.Container.ToValue(item)
	}
	return
}
