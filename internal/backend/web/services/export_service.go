package services

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
)

type ExportService interface {
	GetRowBy(id int64) *collect.Application
	GetCleanList(post *helper.RequestParams) string
}

type exportService struct {
	repo collect.ApplicationImp
}

func NewExportService() ExportService {
	return &exportService{
		repo: collect.NewApplication(),
	}
}

func (this *exportService) GetRowBy(id int64) *collect.Application {
	result, _ := this.repo.GetRowByID(id)
	return result
}

func (this *exportService) GetCleanList(post *helper.RequestParams) string {
	result := this.repo.PostList(post)
	return result
}
