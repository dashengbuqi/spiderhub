package services

import (
	"errors"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"strconv"
)

type CollectService interface {
	GetRowBy(id int64) *collect.Application
	GetCollectList(post *helper.RequestParams) string
	ModifyCollectItem(id int64, form map[string][]string) error
	ModifyCrawler(id int64, content string) error
	Remove(id int64) error
}

type collectService struct {
	repo collect.ApplicationImp
}

func NewCollectService() CollectService {
	return &collectService{
		repo: collect.NewApplication(),
	}
}

func (this *collectService) GetRowBy(id int64) *collect.Application {
	result, _ := this.repo.GetRowByID(id)
	return result
}

func (this *collectService) GetCollectList(post *helper.RequestParams) string {
	result := this.repo.PostList(post)
	return result
}
func (this *collectService) ModifyCrawler(id int64, content string) error {
	if id == 0 {
		return errors.New("请先创建采集任务")
	}
	return this.repo.ModifyCrawlerContent(id, content)
}

//更新数据
func (this *collectService) ModifyCollectItem(id int64, form map[string][]string) error {
	var title, schedule, storage, method string
	if _, ok := form["title"]; ok {
		title = form["title"][0]
	}
	if _, ok := form["schedule"]; ok {
		schedule = form["schedule"][0]
	}
	if _, ok := form["storage"]; ok {
		storage = form["storage"][0]
	}
	if _, ok := form["method"]; ok {
		method = form["method"][0]
	}
	storageInt, _ := strconv.Atoi(storage)
	methodInt, _ := strconv.Atoi(method)
	return this.repo.ModifyItem(id, &collect.Application{
		Title:    title,
		Schedule: schedule,
		Storage:  storageInt,
		Method:   methodInt,
	})
}

func (this *collectService) Remove(id int64) error {
	if id == 0 {
		return errors.New("暂不支持")
	}
	return this.repo.Remove(id)
}
