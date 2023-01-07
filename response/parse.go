package response

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func AnyToInt(val any) int {
	value, ok := val.(int)
	if ok {
		return value
	} else {
		v := AnyToString(val)
		values, err := strconv.Atoi(v)
		if IsError(err) {
			return values
		} else {
			return 0
		}
	}
}

func AnyToString(val any) string {
	value, ok := val.(string)
	if ok {
		return value
	} else {
		return ""
	}
}

func AnyToBool(val any) bool {
	value, ok := val.(bool)
	if ok {
		return value
	} else {
		return false
	}
}

func StringToInt(val string) int {
	value, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return value
}

func toRequired(key string) string {
	return key + "为必传参数"
}

/* 任意类型转json */
func ToJson(dest any) string {
	var value any
	t := reflect.TypeOf(dest)
	switch t.Kind() {
	case reflect.Ptr:
		v := reflect.ValueOf(dest)
		value = v.Elem().Interface()
	default:
		value = dest
	}
	result, err := json.Marshal(value)
	if err == nil {
		return string(result)
	} else {
		return ""
	}
}
