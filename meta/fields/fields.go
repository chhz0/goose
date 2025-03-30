package fields

import (
	"fmt"
	"sort"
	"strings"
)

// Fields 存储独立的字段
type Fields interface {
	// 返回字段是否存在
	Has(field string) (exist bool)

	// 返回字段的值
	Get(field string) (value string)
}

type Set map[string]string

func (ls Set) Has(field string) bool {
	_, exist := ls[field]
	return exist
}

func (ls Set) Get(field string) string {
	return ls[field]
}

// String 返回字段的格式化字符串
func (ls Set) String() string {
	selector := make([]string, 0, len(ls))
	for k, v := range ls {
		selector = append(selector, fmt.Sprintf("%s=%s", k, v))
	}

	sort.StringSlice(selector).Sort()
	return strings.Join(selector, ",")
}

func (ls Set) AsSelector() Selector {
	return SelectorFromSet(ls)
}
