var configs = {
    fields: [
        /*
        {
            name:"字段名", //必须唯一,推荐使用字母标识
            alias:"别称",  //一般使用中文
            container:"容器", //暂支持  string,array,map
            required:"必须有值", //标识当前抽取项的值是否必须（不能为空）。默认是false，可以为空。
            primary:"主键",//当前抽取项是否作为整条数据的主键组成部分。默认是false。
            attached_method:"请求方式", //如果有下载的请求[download:true] HTTP请求是”GET”还是”POST”。默认是”GET”。
            attached_url_params:"请求参数", // 参数形如a=b&c=d，支持变量替换。与attachedUrl的变量引用方式相同。
            attached_headers:{  // HTTP请求的headers,如：h["User_Agent"] = "xxxx"
                "UserAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
                ......
            },
            download:"下载", //附件是否需要下载  bool   true || false
            children:[
                {
                    name:"..."
                    .......
                }
            ],
            func:"回调函数" //回调函数
        },
        */
        {
            name: "article_title",
            alias: "标题",
            primary: true,
        },
        {
            name: "article_content",
            alias: "内容",
        },
        {
            name: "article_publish_time",
            alias: "发布时间",
        },
        {
            name: "article_author",
            alias: "作者",
        }
    ]
};
//每条数据回调此函数
configs.onEachRow = function(row) {
    //extracts("内容","规则","方式") 获取列表
    //extract("内容","规则","方式")  获取单条数据
    //extract_count("内容","规则","方式")  获取数量
    //方式:0:XPATH,1:JSONPATH
    //e.g  var urls = extracts(body,"//ul[@class='articleList']//div[@class='txt']/a/@href",0);
    var result = {};
    for (var i in row.data) {
        result[i] = row.data[i];
    }
    return result;
};

var clean = new Clean(configs);
clean.start();