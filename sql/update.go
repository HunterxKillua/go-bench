package sql

import (
	"fmt"
	"reflect"
)

func (d *DB) SaveOneByMap(dest any, TableName string, values map[string]any, cond ...any) int { /* 修改单条记录 */
	isPtr := d.ValidConfig.CheckDestPtr(dest, "row")
	if isPtr {
		t := reflect.TypeOf(dest).Elem()
		v := reflect.ValueOf(dest).Elem()
		if t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				value, ok := values[t.Field(i).Name]
				if ok {
					vf := v.Field(i)
					switch vf.Kind() {
					case reflect.String:
						val, ok := value.(string)
						if ok && vf.CanSet() {
							vf.SetString(string(val))
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						val, ok := value.(int)
						if ok && vf.CanSet() {
							vf.SetInt(int64(val))
						}
					case reflect.Bool:
						val, ok := value.(bool)
						if ok && vf.CanSet() {
							vf.SetBool(val)
						}
					}
				}
			}
			result := d.db.Save(dest)
			return int(result.RowsAffected)
		}
		if t.Kind() == reflect.Map {
			if len(cond) >= 1 {
				result := d.db.Table(TableName).Where(cond[0], cond[1:]).Limit(1).Updates(values)
				return int(result.RowsAffected)
			} else {
				fmt.Println("缺少where条件")
			}
		}
	}
	return 0
}

func (d *DB) SaveOneByModel(dest any, column string, value any, cond ...any) int {
	if d.ValidConfig.IsStructPtr(dest) {
		if len(cond) >= 1 {
			result := d.db.Model(dest).Where(cond[0], cond[1:]).Limit(1).Update(column, value)
			return int(result.RowsAffected)
		} else { /* 跳过hook */
			result := d.db.Model(dest).UpdateColumn(column, value)
			return int(result.RowsAffected)
		}
	}
	return 0
}

func (d *DB) SaveByModel(dest any, values any, cond ...any) int {
	if d.ValidConfig.IsStructPtr(dest) {
		if len(cond) >= 1 {
			result := d.db.Model(dest).Where(cond[0], cond[1:]).Updates(values)
			return int(result.RowsAffected)
		} else {
			fmt.Println("缺少where条件")
		}
		return 0
	}
	return 0
}
