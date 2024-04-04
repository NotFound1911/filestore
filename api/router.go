package api

import (
	"fmt"
	"github.com/NotFound1911/filestore/api/v1"
	"github.com/NotFound1911/filestore/api/v1/jwt"
	"github.com/NotFound1911/filestore/api/v1/middleware"
	"github.com/NotFound1911/filestore/repository"
	"github.com/NotFound1911/filestore/repository/dao"
	"github.com/NotFound1911/filestore/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	setApiGroupRoutes(router)
	return router
}
func setApiGroupRoutes(
	router *gin.Engine,
) *gin.RouterGroup {
	group := router.Group("/api/storage/v1")
	{
		group.GET("/file/upload", v1.UploadView)
		group.POST("/file/upload", v1.UploadFile)
		group.GET("/file/upload/successful", v1.UploadSuccessful)
		group.GET("/file/meta", v1.GetFileMeta)

		group.POST("/file/download", v1.DownLoad)

		group.POST("/file/update", v1.UpdateFile)
		group.POST("/file/delete", v1.DeleteFile)
	}
	return group
}
func NewRouter1() *gin.Engine {
	server := gin.Default()
	// gorm
	orm := initOrm()
	// dao
	userDao := dao.NewOrmUser(orm)
	// repository
	userRepo := repository.NewCachedUserRepository(userDao)
	// service
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	userService := service.NewUserService(userRepo)
	hdl := jwt.NewRedisJWTHandler(rdb)
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin())
	userHandler := v1.NewUserHandler(userService, hdl)
	userHandler.RegisterUserRoutes(server)

	return server
}

// todo
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
			//TablePrefix:   "t_",
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
	initMySqlTables(db)
	return db
}

// 自定义 接管gorm日志，打印到文件 or 控制台
func getGormLogWriter() logger.Writer {
	var writer io.Writer

	writer = os.Stdout
	return log.New(writer, "\r\n", log.LstdFlags)
}

// 数据库表初始化
func initMySqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		dao.FileMeta{},
		dao.UserInfo{},
	)
	if err != nil {
		panic(err.Error())
	}
}
