package backend

import (
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/configs"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func Run() {
	param, err := configs.GetParams("Backend")
	if err != nil {
		panic(err)
	}
	params := param.(map[interface{}]interface{})
	app := iris.New()
	//输入IRIS日志
	app.Logger().SetLevel(params["Level"].(string))
	app.Favicon("./assets/favicon.ico")
	app.RegisterView(iris.HTML("./web/views", ".html").Layout("layout/main.html").Reload(params["Reload"].(bool)))
	app.HandleDir("/static", iris.Dir("./assets"))
	//默认
	mvc.Configure(app.Party("/"), index)
	//登录
	mvc.Configure(app.Party("/login"), login)

	err = app.Run(iris.Addr(params["Addr"].(string)), iris.WithConfiguration(
		iris.Configuration{
			DisableStartupLog:                 false,
			DisableInterruptHandler:           false,
			DisablePathCorrection:             false,
			EnablePathEscape:                  false,
			FireMethodNotAllowed:              false,
			DisableBodyConsumptionOnUnmarshal: false,
			DisableAutoFireStatusCode:         false,
			TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
			Charset:                           "UTF-8"}),
		iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
	}
}

func index(app *mvc.Application) {

}

func login(app *mvc.Application) {

}
