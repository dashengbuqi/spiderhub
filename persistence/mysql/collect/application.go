package collect

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_data"
	"github.com/go-xorm/xorm"
	"strings"
	"time"
)

var (
	statusArr = map[int]string{
		common.STATUS_NORMAL:  "正常",
		common.STATUS_RUNNING: "执行中",
	}
	methodArr = map[int]string{
		common.METHOD_INSERT: "重新抓取",
		common.METHOD_UPDATE: "更新",
		common.METHOD_APPEND: "追加",
	}
	storageArr = map[int]string{
		common.STORAGE_NORMAL: "",
		common.STORAGE_SERVER: "服务器",
		common.STORAGE_PAN:    "云盘",
	}
)

type Application struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	UserId         int64  `json:"user_id"`
	CrawlerToken   string `json:"crawler_token"`
	CleanToken     string `json:"clean_token"`
	UiCleanData    bool   `json:"ui_clean_data" xorm:"-"`
	Status         int    `json:"status"` //状态(0完成1执行中)
	UiStatus       string `json:"ui_status" xorm:"-"`
	Schedule       string `json:"schedule"` //计划任务
	Storage        int    `json:"storage"`  //存储附件(0不存1服务器)
	UiStorage      string `json:"ui_storage" xorm:"-"`
	Method         int    `json:"method"` //抓取方式(1重新抓取2更新3追加)
	UiMethod       string `json:"ui_method" xorm:"-"`
	ErrorInfo      string `json:"error_info"`
	UiErrorInfo    string `json:"ui_error_info" xorm:"-"`
	CrawlerContent string `json:"crawler_content"`
	CleanContent   string `json:"clean_content"`
	UpdatedAt      int64  `json:"updated_at"`
	UiUpdatedAt    string `json:"ui_updated_at" xorm:"-"`
	CreatedAt      int64  `json:"created_at"`
	UiCreatedAt    string `json:"ui_created_at" xorm:"-"`
}

func (Application) TableName() string {
	return "collect_app"
}

func (this *Application) callUI() {
	this.UiUpdatedAt = helper.FmtDateTime(this.UpdatedAt)
	this.UiCreatedAt = helper.FmtDateTime(this.CreatedAt)
	this.UiStatus = statusArr[this.Status]
	this.UiMethod = methodArr[this.Method]
	this.UiStorage = storageArr[this.Storage]
	this.UiErrorInfo = ""
	if len(this.ErrorInfo) > 0 {
		this.UiErrorInfo = ""
	}
	//检查是否存在数据
	dataTable := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_DATA, this.CleanToken)
	sd := spiderhub_data.NewCollectData(dataTable)
	this.UiCleanData = sd.Has()
}

func (this *Application) GetStorageComboList() string {
	items := []helper.ComboData{
		{
			Id:   common.STORAGE_NORMAL,
			Text: "请选择附件存储",
		},
		{
			Id:   common.STORAGE_SERVER,
			Text: "服务器",
		},
		{
			Id:   common.STORAGE_PAN,
			Text: "云盘",
		},
	}
	result, _ := json.Marshal(items)
	return string(result)
}

func (this *Application) GetMethodComboList() string {
	items := []helper.ComboData{
		{
			Id:   0,
			Text: "请选择数据存储",
		},
		{
			Id:   common.METHOD_INSERT,
			Text: "重新抓取",
		},
		{
			Id:   common.METHOD_UPDATE,
			Text: "数据更新",
		},
		{
			Id:   common.METHOD_APPEND,
			Text: "追加数据",
		},
	}
	result, _ := json.Marshal(items)
	return string(result)
}

type ApplicationImp interface {
	ModifyStatus(id int64, state int) error
	ModifyToken(id int64, token string) error
	ModifyCleanToken(id int64, token string) error
	GetRowByID(id int64) (*Application, error)
	PostList(req *helper.RequestParams) string
	ModifyItem(id int64, item *Application) error
	Remove(id int64) error
	ModifyCrawlerContent(id int64, content string) error
	ModifyCleanContent(id int64, content string) error
}

type application struct {
	session *xorm.Engine
}

func NewApplication() ApplicationImp {
	return &application{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}

func (this *application) ModifyToken(id int64, token string) error {
	var item Application
	item.CrawlerToken = token
	item.UpdatedAt = time.Now().Unix()
	_, err := this.session.Where("id=?", id).Cols("crawler_token", "updated_at").Update(item)
	return err
}
func (this *application) ModifyCleanToken(id int64, token string) error {
	var item Application
	item.CleanToken = token
	item.UpdatedAt = time.Now().Unix()
	_, err := this.session.Where("id=?", id).Cols("clean_token", "updated_at").Update(item)
	return err
}

//更新爬虫状态
func (this *application) ModifyStatus(id int64, state int) error {
	var item Application
	item.Status = state
	_, err := this.session.Where("id=?", id).Cols("status").Update(item)
	return err
}

func (this *application) GetRowByID(id int64) (*Application, error) {
	var item Application
	_, err := this.session.Where("id=?", id).Get(&item)
	return &item, err
}

func (this *application) PostList(req *helper.RequestParams) string {
	var query *xorm.Session
	var where string
	if req.Params != nil {
		params := req.Params.(map[string]interface{})
		if _, ok := params["title"]; ok {
			where = fmt.Sprintf("title like '%?%'", strings.TrimSpace(params["title"].(string)))
		}
	}

	query = this.session.Table(new(Application)).Where(where)
	result := this.assembleTable(query, req)
	return result.ToJson()
}

func (this *application) assembleTable(query *xorm.Session, req *helper.RequestParams) *helper.ResultEasyUItem {
	pages := &helper.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	var sortStr string

	if len(req.Sort) > 0 {
		sortKeys := strings.Split(req.Sort, ",")
		sortValues := strings.Split(req.Order, ",")
		for i, key := range sortKeys {
			sortStr += key + " " + sortValues[i] + ","
		}
	}
	var items []*Application
	limit := pages.GetLimit()
	offset := pages.GetOffset()
	sortBy := strings.Trim(sortStr, ",")

	total, err := query.OrderBy(sortBy).Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		spiderhub.Logger.Error("%v", err.Error())
	}
	for _, item := range items {
		item.callUI()
	}
	pages.Total = total
	return &helper.ResultEasyUItem{
		Pages:  pages,
		Models: items,
	}
}
func (this *application) ModifyItem(id int64, item *Application) error {
	var err error
	if id == 0 {
		if len(item.Title) == 0 {
			return errors.New("任务名称不能为空")
		}
		if item.Method == 0 {
			return errors.New("请选择数据存储方式")
		}
		item.UserId = 0
		item.Status = common.STATUS_NORMAL
		item.CreatedAt = time.Now().Unix()
		_, err = this.session.InsertOne(item)
	} else {
		item.UpdatedAt = time.Now().Unix()
		cols := []string{
			"title", "schedule", "storage", "method", "updated_at",
		}
		_, err = this.session.Where("id=?", id).Cols(cols...).Update(item)
	}
	return err
}

func (this *application) Remove(id int64) error {
	_, err := this.session.Where("id=?", id).Delete(new(Application))
	return err
}

func (this *application) ModifyCrawlerContent(id int64, content string) error {
	var item Application
	item.CrawlerContent = content
	item.UpdatedAt = time.Now().Unix()
	_, err := this.session.Where("id =?", id).Cols("crawler_content", "updated_at").Update(item)
	return err
}

func (this *application) ModifyCleanContent(id int64, content string) error {
	var item Application
	item.CleanContent = content
	item.UpdatedAt = time.Now().Unix()
	_, err := this.session.Where("id =?", id).Cols("clean_content", "updated_at").Update(item)
	return err
}
