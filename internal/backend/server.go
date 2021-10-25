package backend

import (
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/configs"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/controllers"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"path"
	"runtime"
)

func Run() {
	base := getCurrentPath()
	param, err := configs.GetParams("Backend")
	if err != nil {
		panic(err)
	}
	params := param.(map[interface{}]interface{})
	app := iris.New()
	//输入IRIS日志
	app.Logger().SetLevel(params["Level"].(string))
	app.Favicon(base + "/assets/favicon.ico")
	tmpl := iris.HTML(base+"/web/views", ".html").Layout("layout/main.html").Reload(params["Reload"].(bool))
	app.RegisterView(tmpl)
	app.HandleDir("/static", base+"/assets")
	//登录
	mvc.Configure(app.Party("/login"), login)
	//默认
	mvc.Configure(app.Party("/"), index)
	//框架
	mvc.Configure(app.Party("/default"), index)
	//菜单
	mvc.Configure(app.Party("/menu"), menu)
	//用户管理
	mvc.Configure(app.Party("/user"), user)
	//采集
	mvc.Configure(app.Party("/collect"), collect)
	//清洗
	mvc.Configure(app.Party("/clean"), clean)
	//导出数据
	mvc.Configure(app.Party("/export"), export)

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
	menuService := services.NewMenuService()
	app.Register(menuService)
	app.Handle(new(controllers.DefaultController))
}

func export(app *mvc.Application) {
	es := services.NewExportService()
	app.Register(es)
	app.Handle(new(controllers.ExportController))
}

func collect(app *mvc.Application) {
	cs := services.NewCollectService()
	app.Register(cs)
	app.Handle(new(controllers.CollectController))
}

func clean(app *mvc.Application) {
	cs := services.NewCleanService()
	app.Register(cs)
	app.Handle(new(controllers.CleanController))
}

func menu(app *mvc.Application) {
	menuService := services.NewMenuService()
	app.Register(menuService)
	app.Handle(new(controllers.MenuController))
}

func user(app *mvc.Application) {
	us := services.NewUserService()
	app.Register(us)
	app.Handle(new(controllers.UserController))
}

func login(app *mvc.Application) {

}

func getCurrentPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
