package services

import (
	"errors"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/system"
)

type UserService interface {
	GetRowBy(id int64) *system.SystemAdmin
	ModifyMenuItem(id int64, form map[string][]string) error
	GetUserList(post *helper.RequestParams) string
	RemoveUser(id int64) error
	GetByUsernameAndPwd(username, password string) (*system.SystemAdmin, error)
}

type userService struct {
	repo *system.Admin
}

func NewUserService() UserService {
	return &userService{
		repo: system.NewAdmin(),
	}
}

func (this *userService) GetByUsernameAndPwd(username string, pwd string) (*system.SystemAdmin, error) {
	u, err := this.repo.GetByName(username)
	if err != nil {
		return nil, err
	}
	if u.Status == system.ADMIN_STATUS_DISABLE {
		return nil, errors.New("账号已被冻结")
	}
	p := helper.MakeHash(pwd, u.AuthKey)
	if p != u.Password {
		return nil, errors.New("密码错误")
	}
	return u, nil
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
	return this.repo.ModifyItem(id, &system.SystemAdmin{
		Username: username,
		Mobile:   mobile,
		Email:    email,
		Pwd:      password,
	})
}

func (this *userService) GetUserList(post *helper.RequestParams) string {
	result := this.repo.PostMenuList(post)
	return result
}

func (this *userService) GetRowBy(id int64) *system.SystemAdmin {
	result, _ := this.repo.GetRowBy(id)
	return result
}

func (this *userService) RemoveUser(id int64) error {
	if id == 0 {
		return errors.New("暂不支持")
	}
	return this.repo.RemoveItem(id)
}
