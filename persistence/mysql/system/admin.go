package system

import (
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
)

type SystemAdmin struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	RealName  string `json:"real_name"`
	Mobile    string `json:"mobile"`
	AuthKey   string `json:"auth_key"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
	LastLogin int64  `json:"last_login"`
	Logins    int    `json:"logins"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `json:"created_at"`
}

type Admin struct {
	session *xorm.Engine
}

func NewAdmin() *Admin {
	return &Admin{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}
