package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/backend/web/services"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/persistence/mongo/spiderhub_data"
	"github.com/dashengbuqi/spiderhub/persistence/mysql/collect"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

type ExportController struct {
	Ctx     iris.Context
	Service services.ExportService
	Session *sessions.Session
}

//加载列表视图
func (this *ExportController) GetList() mvc.Result {
	return &mvc.View{
		Name: "export/list.html",
	}
}

func (this *ExportController) PostList() string {
	page, _ := this.Ctx.PostValueInt("page")
	pageSize, _ := this.Ctx.PostValueInt("rows")
	sort := this.Ctx.PostValueDefault("sort", "id")
	order := this.Ctx.PostValueDefault("order", "desc")
	result := this.Service.GetCleanList(&helper.RequestParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	})
	return result
}

func (this *ExportController) GetDownload() {
	id, _ := this.Ctx.URLParamInt64("id")
	af := collect.NewAppField()
	item, _ := af.GetRowByID(collect.TARGET_CLEAN, id)
	var headers []*common.TableHead
	if len(item.Content) > 0 {
		json.Unmarshal([]byte(item.Content), &headers)
	}
	app := this.Service.GetRowBy(id)
	table := fmt.Sprintf("%s%s", common.PREFIX_CLEAN_DATA, app.CleanToken)
	cd := spiderhub_data.NewCollectData(table)
	var page int64 = 1
	var pageSize int64 = 100
	f := excelize.NewFile()
	sheet := f.NewSheet(app.Title)
	//设置头
	headMap := make(map[string]string)
	for i, head := range headers {
		c := string(rune(65 + i))
		f.SetCellValue(app.Title, c+"1", head.Alias+"|"+head.Name)
		headMap[head.Name] = c
	}
	for {
		data, _ := cd.GetRowsBy((page-1)*pageSize, pageSize)
		if len(data) == 0 {
			break
		}
		for i, item := range data {
			for name, value := range item {
				if _, ok := headMap[name]; ok {
					if value != nil {
						tp := value.(map[string]interface{})["type"]
						if tp == "map" || tp == "array" {
							vStr, _ := json.Marshal(value.(map[string]interface{})["value"])
							f.SetCellValue(app.Title, headMap[name]+strconv.Itoa(i+2), string(vStr))
						} else {
							f.SetCellValue(app.Title, headMap[name]+strconv.Itoa(i+2), value.(map[string]interface{})["value"])
						}
					} else {
						f.SetCellValue(app.Title, headMap[name]+strconv.Itoa(i+2), "")
					}

				}
			}
		}
		page++
	}
	f.SetActiveSheet(sheet)
	this.Ctx.ResponseWriter().Header().Set("Content-Type", "application/octet-stream")
	dsp := fmt.Sprintf("attachment;filename=%s%s.xlsx", app.Title, time.Unix(time.Now().Unix(), 0).Format("20060102"))
	this.Ctx.ResponseWriter().Header().Set("Content-Disposition", dsp)
	_ = f.Write(this.Ctx.ResponseWriter())
}
