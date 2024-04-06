package ioc

import (
	"fmt"
	"github.com/NotFound1911/filestore/app/account/repository/dao"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

func InitDb() *gorm.DB {
	return initOrm()
}
func initOrm() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		"127.0.0.1",
		"postgres",
		"123456",
		"postgres",
		"5432",
	)

	gormConfig := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	gormConfig = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		Logger: logger.New(getGormLogWriter(), logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Info,            // 日志级别
			IgnoreRecordNotFoundError: false,                  // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                  // 禁用彩色打印
		}), // 使用自定义 Logger
	}

	// gorm将类名转换成数据库表名的逻辑
	if gormConfig.NamingStrategy == nil {
		gormConfig.NamingStrategy = schema.NamingStrategy{
			SingularTable: true,
		}
	}
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	// 执行数据库脚本建表
	initSqlTables(db)
	return db
}

// 自定义 接管gorm日志，打印到文件 or 控制台
func getGormLogWriter() logger.Writer {
	var writer io.Writer

	writer = os.Stdout
	return log.New(writer, "\r\n", log.LstdFlags)
}

// 数据库表初始化
func initSqlTables(db *gorm.DB) {
	err := dao.InitTables(db)
	if err != nil {
		panic(err.Error())
	}
}
