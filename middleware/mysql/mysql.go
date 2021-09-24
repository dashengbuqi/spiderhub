package mysql

import (
	"github.com/dashengbuqi/spiderhub/configs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var Engine map[string]*xorm.Engine

const (
	DATABASE_SPIDERHUB = "spiderhub"
)

func init() {
	params, _ := configs.GetParams("Mysql")

	host := params.(map[interface{}]interface{})["Host"].(string)
	user := params.(map[interface{}]interface{})["User"].(string)
	pwd := params.(map[interface{}]interface{})["Password"].(string)
	dbs := params.(map[interface{}]interface{})["Dbs"].([]interface{})

	if len(dbs) > 0 {
		Engine = make(map[string]*xorm.Engine)
		var engin *xorm.Engine
		var err error
		for _, db := range dbs {
			engin, err = xorm.NewEngine("mysql", user+":"+pwd+"@tcp("+host+")/"+db.(string)+"?charset=utf8")
			if err != nil {
				panic(err)
			}
			Engine[db.(string)] = engin
		}
	}
}
