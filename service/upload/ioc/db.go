package ioc

import (
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/service/upload/repository/dao"
	"gorm.io/gorm"
)

func InitDb(conf *config.Configuration) *gorm.DB {
	db := config.InitDb(conf)
	// 执行数据库脚本建表
	initSqlTables(db)
	return db
}

// 数据库表初始化
func initSqlTables(db *gorm.DB) {
	err := dao.InitTables(db)
	if err != nil {
		panic(err.Error())
	}
}
