package test

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	scan := "http://www.msensor.com.cn/product/{1-300}.html"

	reg := regexp.MustCompile(`{(\d+)-(\d+)}`)
	matRes := reg.FindStringSubmatch(scan)
	var start, end int
	var old string
	if len(matRes) > 0 {
		old = matRes[0]
		start, _ = strconv.Atoi(matRes[1])
		end, _ = strconv.Atoi(matRes[2])
	}
	t.Log(start, end)

	for i := start; i <= end; i++ {
		t.Log(strings.Replace(scan, old, strconv.Itoa(i), -1))
	}
}

func TestInterface(t *testing.T) {
	var m interface{}
	m = map[string]interface{}{
		"key1": map[bool]interface{}{
			false: "val1",
		},
		"key2": map[bool]map[string]interface{}{
			false: {
				"hello": "world",
			},
		},
	}
	for _, val := range m.(map[string]interface{}) {
		t.Log(val)
	}
}
