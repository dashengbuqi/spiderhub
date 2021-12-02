spiderhub 支持 js 脚本编写规则实时调试抓取目标网站 


### 配置文件设置
- configs/ymls/params.yml

### 数据库表结构
- doc/dbs

### 初始化

> go mod tidy

> go mod vendor

### 应用启动

> 爬虫
- cd cmd/crawler
- go build
- nohup ./crawler &

> 清洗
- cd cmd/cleaner
- go build
- nohup ./cleaner &

> 后台管理
- cd cmd/backend
- go build
- nohup ./backend &

> 后台访问地址
- http://127.0.0.1:8080
- 用户名：test   密码：123456

### demo
- 爬虫： /crawler/demo/...
- 清洗： /cleaner/demo/...
