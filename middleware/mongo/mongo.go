package mongo

import (
	"context"
	"github.com/dashengbuqi/spiderhub/configs"
	"github.com/qiniu/qmgo"
)

var client *qmgo.Client
var MongoEngine map[string]*qmgo.Database

const (
	MONGO_MAIN = "spider_main"
	MONGO_DATA = "spider_data"
)

func init() {
	ctx := context.Background()

	params, _ := configs.GetParams("Mongo")

	host := params.(map[interface{}]interface{})["Host"].(string)
	port := params.(map[interface{}]interface{})["Port"].(string)
	user := params.(map[interface{}]interface{})["User"].(string)
	pwd := params.(map[interface{}]interface{})["Password"].(string)
	dbs := params.(map[interface{}]interface{})["Dbs"].([]interface{})
	var uri string = "mongodb://"

	if len(user) > 0 && len(pwd) > 0 {
		uri += user + ":" + pwd + "@"
	}
	if len(host) > 0 {
		uri += host
	}
	if len(port) > 0 {
		uri += ":" + port
	}
	uri += "/?charset=utf8"
	var err error
	client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})
	if err != nil {
		panic(err)
	}
	/*defer func() {
		if err = client.Close(ctx); err != nil {
			panic(err)
		}
	}()*/
	if len(dbs) > 0 {
		MongoEngine = make(map[string]*qmgo.Database)
		//URI example: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
		for _, db := range dbs {
			MongoEngine[db.(string)] = client.Database(db.(string))
		}
	}
}
