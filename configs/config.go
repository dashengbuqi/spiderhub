package configs

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_inst  *Config
	osType = runtime.GOOS
)

type Config struct {
	data map[string]interface{}
}

func init() {
	initConfig()
}

func initConfig() {
	path := CurDir()
	sp := GetSeparator()
	fp := path + sp + "ymls" + sp + "params.yml"
	_inst = New(fp)
}

func New(filepath string) *Config {
	//检查文件是否存在
	if ext := FileExist(filepath); ext == false {
		panic("yml config file path not found")
	}
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	d := make(map[string]interface{})
	if err := yaml.Unmarshal(raw, &d); err != nil {
		panic(err)
	}
	return &Config{
		data: d,
	}
}

//获取应用配置
func (this *Config) GetParams(app string) (interface{}, error) {
	env, err := this.GetEnv()
	if err != nil {
		return nil, err
	}
	if _, ok := this.data[app]; ok == false {
		return nil, errors.New("App Config not Fount")
	}
	params := this.data[app].(map[interface{}]interface{})
	if _, ok := params[env]; ok == false {
		return nil, errors.New("not Found Env:" + env)
	}
	return params[env], nil
}

func (this *Config) GetParam(app string, field string) (interface{}, error) {
	params, err := this.GetParams(app)
	if err != nil {
		return nil, err
	}
	if len(field) == 0 {
		return params, nil
	}
	if _, ok := params.(map[interface{}]interface{})[field]; ok == false {
		return nil, errors.New("Not Found Field:" + field)
	}
	return params.(map[interface{}]interface{})[field], nil
}

func (this *Config) GetEnv() (string, error) {
	if _, ok := this.data["Env"]; ok == false {
		return "", errors.New("Config Not Found Env,Env is Require")
	}
	return this.data["Env"].(string), nil
}

func GetParams(app string) (interface{}, error) {
	if _inst == nil {
		initConfig()
	}
	return _inst.GetParams(app)
}

func GetParamsByField(app string, field string) (interface{}, error) {
	if _inst == nil {
		initConfig()
	}
	return _inst.GetParam(app, field)
}

func GetEnv() (string, error) {
	if _inst == nil {
		initConfig()
	}
	return _inst.GetEnv()
}

func FileExist(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	} else {
		return true
	}
}

func CurDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func GetSeparator() string {
	if osType == "windows" {
		return "\\"
	} else if osType == "linux" {
		return "/"
	}
	return "/"
}
