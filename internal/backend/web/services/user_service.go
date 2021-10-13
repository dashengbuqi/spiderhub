package services

import (
	"fmt"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
)

type UserService interface {
	GetRowBy(id int64) *system.SystemAdmin
	ModifyMenuItem(id int64, form map[string][]string) error
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

//更新数据
func (this *userService) ModifyMenuItem(id int64, form map[string][]string) error {
	var username, mobile, email, password string
	if _, ok := form["username"]; ok {
		username = form["username"][0]
	}
	if _, ok := form["mobile"]; ok {
		mobile = form["mobile"][0]
	}
	if _, ok := form["email"]; ok {
		email = form["email"][0]
	}
	if _, ok := form["password"]; ok {
		password = form["password"][0]
	}

	fmt.Println(username, mobile, email, password)
	return nil
	/*return this.repo.ModifyItem(id, &system.SystemMenu{
		TaskName: task_name,
		FullName: full_name,
		Path:     path,
		Icon:     icon,
		Sort:     sort,
		Type:     tp,
		ParentId: int64(parent_id),
	})*/
}

func (this *userService) GetUserList(post *helper.RequestParams) string {
	result := this.repo.PostMenuList(post)
	return result
}

func (this *userService) GetRowBy(id int64) *system.SystemAdmin {
	result, _ := this.repo.GetRowBy(id)
	return result
}
