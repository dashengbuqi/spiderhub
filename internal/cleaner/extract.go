package cleaner

import (
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/robertkrimen/otto"
)

//提取字段
type Extract struct {
	container *otto.Otto
	outLog    chan<- []byte
	ov        otto.Value
	fields    []FieldStash
}

func NewExtract(ov otto.Value, fields []FieldStash, vm *otto.Otto, out chan<- []byte) *Extract {
	return &Extract{
		container: vm,
		outLog:    out,
		ov:        ov,
		fields:    fields,
	}
}

func (this *Extract) Run() interface{} {
	defer func() {
		p := recover()
		if p != nil {
			this.outLog <- helper.FmtLog(common.LOG_ERROR, p.(error).Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
	}()
	result := make(map[string]interface{})
	for _, field := range this.fields {
		value, _ := this.ov.Object().Get(field.Name)
		result[field.Name] = map[bool]interface{}{
			field.Primary: this.recursExtract(value, field),
		}
	}
	return result
}

//递归提取字段
func (this *Extract) recursExtract(value otto.Value, field FieldStash) *common.FieldData {
	var result *common.FieldData
	//待提取字段类型
	fieldType := field.Type
	if len(fieldType) == 0 {
		fieldType = TYPE_STRING
	}

	if fieldType == TYPE_STRING {
		result = &common.FieldData{
			Alias: field.Alias,
			Type:  fieldType,
			Value: value.String(),
		}
	} else if fieldType == TYPE_MAP {
		//检查是否有子项
		if len(field.Children) > 0 {
			subResult := make(map[string]interface{})
			for _, subField := range field.Children {
				subValue, _ := value.Object().Get(subField.Name)
				if subValue.IsObject() {
					subResult[subField.Name] = this.recursExtract(subValue, subField)
				}
			}
			result = &common.FieldData{
				Alias: field.Alias,
				Type:  TYPE_MAP,
				Value: subResult,
			}
		} else {
			if value.IsObject() {
				var err error
				subResult := make(map[string]interface{})
				for _, key := range value.Object().Keys() {
					val, _ := value.Object().Get(key)
					subResult[key], err = val.ToString()
					if err != nil {
						spiderhub.Logger.Error("%s", err.Error())
					}
				}
				result = &common.FieldData{
					Alias: field.Alias,
					Type:  TYPE_MAP,
					Value: subResult,
				}
			}
		}
	}

}
