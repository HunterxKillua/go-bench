package utils

import (
	"reflect"
	"strings"
)

/* 提取struct中声明的query字段 */
func GetQueryByModel(dest any) map[string]string {
	var query map[string]string
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Struct {
			query = forEachTag(t.Elem(), "uri")
		}
	}
	return query
}

/* 提取struct中声明的form字段 */
func GetFormByModel(dest any) map[string]string {
	var form map[string]string
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Struct {
			form = forEachTag(t.Elem(), "form")
		}
	}
	return form
}

func forEachTag(dest reflect.Type, typ string) map[string]string {
	var query = map[string]string{}
	for i := 0; i < dest.NumField(); i++ {
		field := dest.Field(i)
		v, ok := field.Tag.Lookup(typ)
		if ok {
			value, oks := field.Tag.Lookup("default")
			if oks {
				query[v] = value
			} else {
				query[v] = ""
			}
		}
	}
	return query
}

func Try(userFn func(), catchFn func(err interface{})) {
	defer func() {
		if err := recover(); err != nil {
			catchFn(err)
		}
	}()
	userFn()
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func ToSnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func ToCamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
			d := s[i]
			if !k && d >= 'A' && d <= 'Z' {
					k = true
			}
			if d >= 'a' && d <= 'z' && (j || !k) {
					d = d - 32
					j = false
					k = true
			}
			if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
					j = true
					continue
			}
			data = append(data, d)
	}
	return string(data[:])
}