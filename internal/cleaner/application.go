package cleaner

import (
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/dashengbuqi/spiderhub"
	"github.com/oliveagle/jsonpath"
	"github.com/robertkrimen/otto"
	"strings"
)

const (
	SELECTORTYPE_XPATH int64 = iota
	SELECTORTYPE_JSONPATH
	SELECTORTYPE_REGEX

	FIELDS = "fields"

	FUNC_ON_EACH_ROW       = "onEachRow"
	FUNC_BEFORE_HANDLE_IMG = "beforeHandleImg"
	FUNC_BEFORE_HOST_FILE  = "beforeHostFile"
	FUNC_AFTER_HOST_FILE   = "afterHostFile"

	TYPE_STRING = "string"
	TYPE_ARRAY  = "array"
	TYPE_MAP    = "map"

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

type FieldStash struct {
	//抽取项的名字。
	Name string `json:"name"`
	//抽取项的别名。一般起中文名，方便查看数据。只影响网页的上显示，可随意修改。
	Alias string `json:"alias"`
	//存储容器 string,array,map
	Type string `json:"type"`
	//标识当前抽取项的值是否必须（不能为空）。默认是false，可以为空。
	Required bool `json:"required"`
	// 抽取项的子抽取项。
	//field支持子项，可以设置多层级，方便数据以本身的层级方式存储，而不用全部展开到第一层级。
	//注意：
	//第一层field默认从当前网页内容中抽取，而子项默认从父项的内容中抽取
	Children []FieldStash `json:"children"`
	//当前抽取项是否作为整条数据的主键组成部分。默认是false。
	Primary bool `json:"primary"`

	// attachedUrl请求地址。
	//attachedUrl支持变量替换，变量可引用上下文中已经抽取的字段。
	//同一层级的field或者第一层级的field，引用方式为”{fieldName}”。
	//不同层级需要从根field开始，$表示根，引用方式为”{$.fieldName.fieldName}”。
	//特殊变量$$url表示当前网页的url，引用方式为”{$$url}”。
	//比如抽取到字段item_id，attachedUrl形式为https://item.example.com/item.htm?id=1000，则attachedUrl可以写为：
	//"https://item.example.com/item.htm?id={item_id}"
	////AttachedUrl string `json:"attached_url"`
	// HTTP请求是”GET”还是”POST”。默认是”GET”。
	AttachedMethod string `json:"attached_method"`
	// HTTP请求的POST参数。如果请求是”GET”，参数将会被忽略。
	// 参数形如a=b&c=d，支持变量替换。与attachedUrl的变量引用方式相同。
	AttachedUrlParams string `json:"attached_url_params"`
	// HTTP请求的headers。
	AttachedHeaders map[string]string `json:"attached_headers"`
	//附件是否需要下载
	Download bool `json:"download"`
	// 抽取项是否是临时的。默认是false。临时的抽取项，数据存储的时候，不会存储其值。
	//Temporary bool `json:"temporary"`
	// 回调函数
	Func string `json:"func"`
}

type Application struct {
	Rules     map[string]interface{}
	Container *otto.Otto
	RuleName  string
	oo        *otto.Object
	Running   bool
}

func NewApplication() *Application {
	return &Application{
		Rules:     make(map[string]interface{}),
		Container: otto.New(),
	}
}

func (this *Application) Init(body string) error {
	st, _ := this.Container.Object(`({XPath:0,JsonPath:1,Regex:2})`)
	this.Container.Set("SelectType", st)
	ua, _ := this.Container.Object(`({Computer:0,Android:1,IOS:2,Mobile:3,Empty:4})`)
	this.Container.Set("UserAgent", ua)
	et, _ := this.Container.Object(`({Normal:0,UrlContext:1,AttachUrl:2})`)
	this.Container.Set("ExtractType", et)
	c, _ := this.Container.Object(`({String:"string",Array:"array",Map:"map"})`)
	this.Container.Set("Type", c)
	_, err := this.Container.Run(body)
	return err
}

func (this *Application) LazyLoad(oo *otto.Object) {
	this.oo = oo

	this.setCallBack(FUNC_ON_EACH_ROW)
	this.setCallBack(FUNC_BEFORE_HANDLE_IMG)
	this.setCallBack(FUNC_BEFORE_HOST_FILE)
	this.setCallBack(FUNC_AFTER_HOST_FILE)

	var err error
	err = this.Container.Set("extracts", this.extractList)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	err = this.Container.Set("extract", this.extract)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
	err = this.Container.Set("count", this.extractCount)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}

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
	this.Running = true
}

func (this *Application) setCallBack(fn string) {
	if fnVal, _ := this.oo.Get(fn); fnVal.IsDefined() {
		err := this.Container.Set(fn, fnVal)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	}
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

func (this *Application) extractCount(call otto.FunctionCall) (result otto.Value) {
	content := call.Argument(0).String()
	xRule := call.Argument(1).String()
	typ, _ := call.Argument(2).ToInteger()

	if typ == SELECTORTYPE_JSONPATH {
		result, _ = this.Container.ToValue(0)
	} else {
		doc, err := htmlquery.Parse(strings.NewReader(content))
		if err != nil {
			return result
		}
		expr := xpath.MustCompile(xRule)
		iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
		var rs int
		for iter.MoveNext() {
			rs++
		}
		result, _ = this.Container.ToValue(rs)
	}
	return result
}
