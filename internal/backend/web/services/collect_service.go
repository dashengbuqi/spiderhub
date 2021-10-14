package services

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
)

type CollectService interface {
	GetRowBy(id int64) *collect.Application
	GetCollectList(post *helper.RequestParams) string
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
	return nil
}

func (this *collectService) GetCollectList(post *helper.RequestParams) string {
	result := this.repo.PostList(post)
	return result
}
