package services

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
)

type UserService interface {
	GetRowBy(id int64) *system.SystemAdmin
	GetUserList(post *helper.RequestParams) string
}

type userService struct {
	repo *system.Admin
}

func NewUserService() UserService {
	return &userService{
		repo: system.NewAdmin(),
	}
}

func (this *userService) GetUserList(post *helper.RequestParams) string {
	result := this.repo.PostMenuList(post)
	return result
}

func (this *userService) GetRowBy(id int64) *system.SystemAdmin {
	result, _ := this.repo.GetRowBy(id)
	return result
}
