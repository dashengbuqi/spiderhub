var configs = {
    name:"应用名称",
    domains: ["leiphone.com"],// 网站域名，设置域名后只处理这些域名下的网页
    scan_urls: [
        "https://www.leiphone.com/search?s=vr&site=article",
        //"https://www.leiphone.com/news/201905/0wo21fMOTrp5DJ3K.html"
    ],// 入口页链接，分别从这些链接开始爬取
    content_regex: [
        "^https://www.leiphone.com/news/(\\d)+/(.*)*\\.html$"
    ],// 内容页url的正则，符合这些正则的页面会被当作内容页处理
    helper_regex: [
        "^https://www.leiphone.com/search\\?(.*)*"
    ],// 列表页url的正则，符合这些正则的页面会被当作列表页处理
    //timeout:60,  //请求超时
    //delay:5，    //请求延迟
    //user_agent:"",   //ua
    // accept_http_status:[200],    //允许的返回状态码，默认：200、201、202、203、204、205、206、207、208、226、301、302
    // auto_find_urls:false,     //否允许自动发现连接
    fields: [
        {
            // 抽取内容页的文章标题
            name: "article_title",
            //抽取项的别名。一般起中文名，方便查看数据。只影响网页的上显示，可随意修改。
            alias:"标题",
            //存储容器 string,array,map  默认string
            container:"string",
            //抽取规则 如果selector为空或者未设置，则抽取的值为null，在进行required的判定之前，仍会进行afterExtractField回调。
            selector: "//h1[contains(@class,'headTit')]",
            //抽取规则的类型。默认值是SelectorType.XPath
            //SelectorType.XPath 一般针对html网页或xml
            //SelectorType.JsonPath 针对json数据
            selector_type:"SelectorType.XPath",
            //内容类型 0:text 1:html  默认 text
            //body_type:"text",
            //当前抽取项是否作为整条数据的主键组成部分。默认是false。
            primary:true,
            // required为true表示该项数据不能为空
            required: true,
            //数据提取方式  ExtractType.Normal    ExtractType.UrlContext    ExtractType.AttachedUrl
            // 数据抽取源。ExtractType.Normal ，抽取源默认时当前网页或父项内容。
            // 访问方式:1 UrlContext  2 AttachedUrl
            // SourceType.UrlContext
            // 表示从当前链接的附加内容中抽取数据。在添加链接的时候，可以同时给该链接附加一些数据，通常的使用场景是，列表页展示的一些内容页没有的数据，
            // 那么在做链接发现的时候，可以直接把这部分数据附加到对应的内容页链接上。
            // SourceType.AttachedUrl
            // 表示需要的数据在另外一个链接（我们叫attachedUrl）的请求结果里面，需要额外再发一次请求来获取数据。
            //只有当sourceType为SourceType.AttachedUrl时，下面的attachedUrl、attachedMethod、attachedParams、attachedHeaders设置才有意义。
            // 多说两句：在还没有site.requestUrl接口的时候，很多复杂的场景（比如文章内容分页）只能用attachedUrl来实现；
            // 相比于attachedUrl，site.requestUrl具有更高的灵活性，所以对于复杂的场景，我们建议在回调函数中通过site.requestUrl来实现，
            // attachedUrl只用来实现一些简单的单请求场景。
            extract_method: ExtractType.Normal
            // attachedUrl请求地址。
            //attachedUrl支持变量替换，变量可引用上下文中已经抽取的字段。
            //同一层级的field或者第一层级的field，引用方式为”{fieldName}”。
            //不同层级需要从根field开始，$表示根，引用方式为”{$.fieldName.fieldName}”。
            //特殊变量$$url表示当前网页的url，引用方式为”{$$url}”。
            //比如抽取到字段item_id，attachedUrl形式为https://item.example.com/item.htm?id=1000，则attachedUrl可以写为：
            //"https://item.example.com/item.htm?id={item_id}"
            ////attached_url:"",
            // HTTP请求是”GET”还是”POST”。默认是”GET”。
            ////attached_method:"",
            // HTTP请求的POST参数。如果请求是”GET”，参数将会被忽略。
            // 参数形如a=b&c=d，支持变量替换。与attachedUrl的变量引用方式相同。
            //// attached_url_params:"",
            // HTTP请求的headers key:value。
            //// attached_headers:{"key":value},
            // 抽取项是否是临时的。默认是false。临时的抽取项，数据存储的时候，不会存储其值。
            //// temporary:false,
            // 回调函数
            //// func:"callbackfunction",
            // 抽取项的子抽取项。
            //field支持子项，可以设置多层级，方便数据以本身的层级方式存储，而不用全部展开到第一层级。
            //注意：
            //第一层field默认从当前网页内容中抽取，而子项默认从父项的内容中抽取
            ////children:[
            ////     {
            ////        name:"",
            ////        alias:"",
            ////   },
            ////]
        },
        {
            // 抽取内容页的文章内容
            name: "article_content",
            selector: "//div[contains(@class,'lph-article-comView')]",
            body_type:"html",
            required: true
        },
        {
            // 抽取内容页的文章发布日期
            name: "article_publish_time",
            selector: "//td[contains(@class,'time')]",
            required: true
        },
        {
            // 抽取内容页的文章作者
            name: "article_author",
            selector: "//td[contains(@class,'aut')]/a",
            required: true
        },
        {
            name:"content_images",
            selector:"//div[contains(@class,'lph-article-comView')]//img/@src",
            alias:"内容图片",
            download:true,
            temporary:true, //临时
            type:"array",
            func: "formatImage"
        }
    ]
};

function formatImage(data,url) {
    return data;
}
/*
 * 队列里添加入口链接
 * 不返回
 */
//configs.initCrawl = function(site) {
//site.AddURL("https://www.leiphone.com/search?s=vr&site=article")
//};

/*
 *访问请求之前调用，可以添加header中的user-agent  cookie
 * 不返回
 */
configs.beforeCrawl = function(req) {
    var url = req.URL.String();
    var re = new RegExp("^https://www.leiphone.com/search\\?(.*)*$","i");
    if (re.test(encodeURI(url))) {
        req.Headers.Add("cookie","PHPSESSID=a596fad52de9a011ca0b268f1807da81; SameSite=None; Hm_lvt_0f7e8686c8fcc36f05ce11b84012d5ee=1614741244; Hm_lpvt_0f7e8686c8fcc36f05ce11b84012d5ee=1614741256");
    } else {
        req.Headers.Add("cookie","PHPSESSID=a596fad52de9a011ca0b268f1807da81; SameSite=None; Hm_lvt_0f7e8686c8fcc36f05ce11b84012d5ee=1614741244; Hm_lpvt_0f7e8686c8fcc36f05ce11b84012d5ee=1614741392");
    }
};
/*
 *判断是否为爬虫
 * @return bool
 */
configs.isAntiSpider = function(url,content) {
    //根据内容或者URL判断为反爬虫  爬虫返回true ,正常返回false;
    return false;
};
/*
 * 是否需要在内容中获取新链接
 */
//configs.onDataReceived = function(data,site) {
//site.AddURL("url");
//};
//爬虫结束时调用，每个节点都会回调，在beforeExit之前被回调。
configs.afterCrawl = function(url) {
    console.log(url +"请求结束");
}

/*
 *需要增加时间戳之类
 * @return string
 * 功能无开通

configs.beforeDownloadPage = function(url) {
    return url;
}; */
//暂未启用
/*configs.afterDownloadPage = function() {

}*/


/*
 * 入口页内容是否需要自动发现新链接 true:自动发现 false:不需要
 * @return bool
 */
configs.onProcessScanPage = function(content,site) {
    return false;
};
/*
 * 在内容中获取链接
 * @param body string  下载内容
 * @param url string 当前请求的链接
 * @param site object 当前队列对象
 * @return bool
 */
configs.onProcessHelperPage = function(body,url,site) {
    //extracts("内容","规则","方式") 获取列表
    //extract("内容","规则","方式")  获取单条数据
    //extract_count("内容","规则","方式")  获取数量
    //方式:0:XPATH,1:JSONPATH
    var urls = extracts(body,"//ul[@class='articleList']//div[@class='txt']/a/@href");
    for (var i in urls) {
        site.AddURL(urls[i]);
    }
    return false;
};
/*
 * 在内容中获取链接
 * @param body string  下载内容
 * @param url string 当前请求的链接
 * @param site object 当前队列对象
 * @return bool    true 自动发现连接
 */
configs.onProcessContentPage = function(body,site) {
    //fieldName[0] = "content_images";
    return false;
};

// 多数场景是，attachedUrl返回的数据是jsonp格式，这时需要在此回调中把数据处理成json数据，以便后续用JsonPath来抽取。
// 暂未启用
/*configs.afterDownloadAttachedPage = function() {

}*/
//@return 数据对象 返回此项对应的数据。当不重写此函数时，默认返回原data对象。
//在每个抽取项抽取到内容时回调此函数，一个网页的抽取过程中，会多次回调此函数。
//在此函数中，可以对抽取到的数据做进一步的处理，然后返回处理后的数据。
/*configs.afterExtractField = function(data) {

}*/

//抽取到临时数据后调用此方法
//@return 数据对象 返回此项对应的数据。当不重写此函数时，默认返回原data对象。
// extract_method: ExtractType.Normal 和 container:"string"时支持
/*configs.afterExtractTemporaryField = function(data) {

}*/
//暂未启用
//在抽取的内容中发现标签时，回调此函数。一般在此函数中修改src，使src指向真实的图片地址。
/*configs.beforeHandleImg = function() {

}*/

//整个网页完成抽取时回调此函数。一般在此回调中做一些数据整理的操作,把返回的数据整理后放到data中返回
//@retrun 整理后的数据
/*configs.afterExtractPage = function(data) {

}*/
// 使用以上配置创建一个爬虫对象
var crawler = new Crawler(configs);
// 启动该爬虫
crawler.start();