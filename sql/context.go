package sql

import (
	"log"
	"time"

	cValidator "ginBlog/valid"

	models "ginBlog/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	db                      *gorm.DB /* 当前连接的数据库 */
	*cValidator.ValidConfig          /* 参数枚举验证 */
}

func ConnectDB(conf *models.ConfDB) *DB {
	db_adr := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(db_adr), &gorm.Config{
		// QueryFields:            false,
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	sqlDB.SetMaxOpenConns(conf.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLive))
	validator := cValidator.Init()
	var instance *DB
	if err != nil {
		instance = &DB{
			db:          nil,
			ValidConfig: validator,
		}
	} else {
		instance = &DB{
			db:          db,
			ValidConfig: validator,
		}
	}
	for _, model := range models.Models {
		instance.db.AutoMigrate(model)
	}
	return instance
}

func (d *DB) GetDB() *gorm.DB {
	if d.db == nil {
		log.Println("请先指定数据库连接")
		return nil
	}
	return d.db
}
