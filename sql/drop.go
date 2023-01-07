package sql

import (
	"reflect"

	"gorm.io/gorm"
)

func (d *DB) DropIndex(dest any, key string) {
	if d.GetDB() != nil && d.db.Migrator().HasIndex(dest, key) {
		d.db.Migrator().DropIndex(dest, key)
	}
}

func (d *DB) RenameIndex(dest any, oKey string, nKey string) {
	if d.GetDB() != nil && d.db.Migrator().HasIndex(dest, oKey) {
		d.db.Migrator().RenameIndex(dest, oKey, nKey)
	}
}

func (d *DB) DropTable(dest ...any) error {
	if d.GetDB() != nil {
		err := d.db.Migrator().DropTable(dest...)
		return err
	}
	return nil
}

func (d *DB) DelCols(dest any, colNames any) error {
	if d.GetDB() != nil {
		value, ok := colNames.(string)
		if ok {
			err := d.db.Migrator().DropColumn(dest, value)
			if err != nil {
				return err
			}
			return nil
		}
		slices, status := colNames.([]string)
		if status {
			for _, col := range slices {
				err := d.db.Migrator().DropColumn(dest, col)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *DB) DelRecords(dest any, value any, cond ...any) bool { /* 条件删除多条 */
	if d.ValidConfig.IsStructPtr(dest) {
		var result *gorm.DB
		if len(cond) > 0 {
			result = d.db.Where(cond[0], cond[1:]).Delete(dest)
		} else {
			result = d.db.Delete(dest, value) /* 主键key */
		}
		return int(result.RowsAffected) > 0
	}
	return false
}

func (d *DB) DelRecord(dest any, cond ...any) bool { /* 从查询记录中条件删除单条记录 */
	t := reflect.TypeOf(dest)
	if d.ValidConfig.IsPtr(dest) {
		var result *gorm.DB
		if t.Elem().Kind() == reflect.Slice { /* 条件删除 */
			result = d.db.Where(cond[0], cond[1:]).Delete(dest)
		} else {
			result = d.db.Delete(dest, cond...)
		}
		return int(result.RowsAffected) > 0
	}
	return false
}

func (d *DB) DelReal(dest any) bool { /* 完全删除 */
	if d.ValidConfig.IsPtr(dest) {
		result := d.db.Unscoped().Delete(dest)
		return int(result.RowsAffected) > 0
	} else {
		return false
	}
}

func (d *DB) GetFake(tableName string, cond ...any) *gorm.DB { /* 获取到软删除的记录(delete_at字段) */
	if len(cond) > 0 {
		return d.db.Table(tableName).Unscoped().Where(cond[0], cond[1:])
	} else {
		return d.db.Table(tableName).Unscoped().Not("deleted is null")
	}
}
