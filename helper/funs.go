package helper

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/axgle/mahonia"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/oliveagle/jsonpath"
	"github.com/robertkrimen/otto"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func FmtConsole(argumentList []otto.Value) string {
	output := []string{}
	for _, argument := range argumentList {
		output = append(output, fmt.Sprintf("%v", argument))
	}
	return strings.Join(output, " ")
}

func FmtUrl(urlStr string) string {
	urlStr = strings.Replace(urlStr, "\\", "", -1)
	urlStr = strings.Replace(urlStr, "\"", "", -1)
	urlStr, _ = url.QueryUnescape(urlStr)
	return urlStr
}

func StringToM5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

var weekDayArr = map[string]string{
	"Monday":    "星期一",
	"Tuesday":   "星期二",
	"Wednesday": "星期三",
	"Thursday":  "星期四",
	"Friday":    "星期五",
	"Saturday":  "星期六",
	"Sunday":    "星期天",
}

//转换成时间，星期，年月日
func FmtDateTime(t int64) string {
	tm := time.Now()
	now := tm.Unix()
	yesterday_start := time.Date(tm.Year(), tm.Month(), tm.Day()-1, 0, 0, 0, 0, tm.Location()).Unix()
	yesterday_end := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()).Unix()
	diff := now - t

	day := int64(3600 * 24)
	week := day * 7
	if t > yesterday_start && t < yesterday_end {
		hm := time.Unix(t, 0).Format("15:04")
		return "昨天" + "" + hm
	} else {
		if diff <= day {
			if diff <= 300 {
				return "刚刚"
			} else if diff < 3600 {
				return fmt.Sprintf("%d 分钟前", diff/60)
			}
			return fmt.Sprintf("%d 小时前", diff/3600)
		} else if diff <= week {
			wd := time.Unix(t, 0).Weekday().String()
			hm := time.Unix(t, 0).Format("15:04")
			return weekDayArr[wd] + " " + hm
		} else {
			res := time.Unix(t, 0).Format("2006-01-02 15:04")
			return res
		}
	}
}

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func ValidUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func AutoFindLinkUrls(body string) (result []string) {
	urls := Extracts(body, "//a/@href", common.SELECTORTYPE_XPATH)
	if len(urls) > 0 {
		for _, url := range urls {
			//if strings.Contains(url, "http") || strings.Contains(url, "https") {
			result = append(result, url.(string))
			//}
		}
	}
	return result
}

func Extracts(content, pathRule string, ruleType int) []interface{} {
	if ruleType == common.SELECTORTYPE_JSONPATH {
		var jsonData interface{}
		json.Unmarshal([]byte(content), &jsonData)
		var rs []interface{}
		pat, _ := jsonpath.Compile(pathRule)
		res, err := pat.Lookup(jsonData)

		if err != nil {
			return rs
		}
		return res.([]interface{})
	} else {
		doc, err := htmlquery.Parse(strings.NewReader(content))
		if err != nil {
			return nil
		}
		expr := xpath.MustCompile(pathRule)
		iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)

		var rs []interface{}
		for iter.MoveNext() {
			if v := iter.Current().Value(); len(v) > 0 {
				rs = append(rs, v)
			}
		}
		return rs
	}
}

func ExtractHtml(content, pathRule string, ruleType int) string {
	if ruleType == common.SELECTORTYPE_JSONPATH {
		var jsonData interface{}
		json.Unmarshal([]byte(content), &jsonData)

		res, err := jsonpath.JsonPathLookup(jsonData, pathRule)

		if err != nil {
			return ""
		}
		return res.(string)
	} else {
		root, err := htmlquery.Parse(strings.NewReader(content))
		if err != nil {
			return ""
		}
		node := htmlquery.FindOne(root, pathRule)
		if node == nil {
			return ""
		}
		return htmlquery.OutputHTML(node, true)
	}
}

func ExtractItem(content, pathRule string, ruleType int) interface{} {
	if ruleType == common.SELECTORTYPE_JSONPATH {
		var jsonData interface{}
		json.Unmarshal([]byte(content), &jsonData)

		res, err := jsonpath.JsonPathLookup(jsonData, pathRule)

		if err != nil {
			return ""
		}
		return res
	} else {
		var rs interface{}
		node, err := htmlquery.Parse(strings.NewReader(content))
		if err != nil {
			return nil
		}
		expr := xpath.MustCompile(pathRule)
		iter := expr.Evaluate(htmlquery.CreateXPathNavigator(node)).(*xpath.NodeIterator)

		for iter.MoveNext() {
			if v := iter.Current().Value(); len(v) > 0 {
				rs = v
				break
			}
		}
		return rs
	}

}
