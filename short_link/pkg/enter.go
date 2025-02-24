package pkg

import (
	"regexp"

	"github.com/zeromicro/go-zero/core/logx"
)

func Inlist(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

// InlistByRegs 根据提供的正则表达式列表，检查给定的 key 是否匹配其中任何一个正则表达式。
// 如果匹配，返回 true；否则返回 false。
//
// 参数:
//   list: 包含多个正则表达式字符串的切片。
//   key: 要匹配的字符串。
//
// 返回值:
//   返回一个布尔值，如果 key 匹配 list 中的任意正则表达式，返回 true；
//   否则返回 false。
func InlistByRegs(list []string, key string) bool {
	for _, s := range list {
		regex, err := regexp.Compile(s)
		if err != nil {
			logx.Error(err)  
			return false
		}
		if regex.MatchString(key) {
			return true
		}
	}
	return false
}



// 去重
func DeduplicationList[T string|int|uint32](req []T) (response []T) {
	m := make(map[T]bool)
	for _, item := range req {
		if _, ok := m[item]; !ok {
			m[item] = true
			response = append(response, item)
		}
	}
	return
}