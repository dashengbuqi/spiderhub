package services

import (
	"encoding/json"
	"errors"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
)

type CleanService interface {
	GetRowBy(id int64) *collect.Application
	GetCleanList(post *helper.RequestParams) string
	GetCleanData(post *helper.RequestParams) string
	ModifyClean(id int64, content string) error
	CleanBegin(id int64, code string) (int64, error)
	CleanHeart(id int64, debug_id int64, user_id int64) interface{}
	CleanEnd(id int64, debug_id int64, user_id int64) error
	CleanStart(id int64) error
	CleanStatus(id int64) error
	GetCleanHead(id int64) []*common.TableHead
}

type cleanService struct {
	repo collect.ApplicationImp
}

func NewCleanService() CleanService {
	return &cleanService{
		repo: collect.NewApplication(),
	}
}

func (this *cleanService) GetRowBy(id int64) *collect.Application {
	result, _ := this.repo.GetRowByID(id)
	return result
}

func (this *cleanService) GetCleanList(post *helper.RequestParams) string {
	result := this.repo.PostList(post)
	return result
}

func (this *cleanService) GetCleanData(post *helper.RequestParams) string {
	return ""
}

func (this *cleanService) ModifyClean(id int64, content string) error {
	if id == 0 {
		return errors.New("请先创建采集任务")
	}
	return this.repo.ModifyCleanContent(id, content)
}

func (this *cleanService) CleanStatus(id int64) error {
	ca := collect.NewApplication()
	row, _ := ca.GetRowByID(id)
	if row.Status == collect.STATUS_RUNNING {
		return errors.New("正在执行中")
	}
	return nil
}

//正式执行
func (this *cleanService) CleanStart(id int64) error {
	if id == 0 {
		return errors.New("缺少参数")
	}
	ca := collect.NewApplication()
	row, _ := ca.GetRowByID(id)
	if len(row.CleanContent) == 0 {
		return errors.New("还没有爬虫规则,点击调试创建吧！")
	}
	cm := &common.Communication{
		AppId:   row.Id,
		UserId:  row.UserId,
		DebugId: 0,
		Method:  common.METHOD_EXCUTE,
		Content: row.CleanContent,
	}
	str, err := json.Marshal(cm)
	if err != nil {
		return err
	}
	err = queue.RabbitConn.Publish(&queue.CleanerChannel, str)
	return err
}

func (this *cleanService) GetCleanHead(id int64) []*common.TableHead {
	return nil
}

//调试开始
func (this *cleanService) CleanBegin(id int64, code string) (int64, error) {
	return 0, nil
}

func (this *cleanService) CleanHeart(id, debug_id, user_id int64) interface{} {
	return nil
}

//终止调试
func (this *cleanService) CleanEnd(id int64, debug_id int64, user_id int64) error {

	return nil
}
