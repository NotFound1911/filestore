package main

import (
	"github.com/NotFound1911/filestore/api/rest/upload/v1"
	"github.com/NotFound1911/filestore/app/upload/ioc"
	"github.com/NotFound1911/filestore/app/upload/repository"
	"github.com/NotFound1911/filestore/app/upload/repository/dao"
	"github.com/NotFound1911/filestore/app/upload/service"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	"github.com/NotFound1911/filestore/internal/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	gin.SetMode(gin.DebugMode)
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	hdl := jwt.NewRedisJWTHandler(rdb)
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin())

	// **************
	// gorm
	orm := ioc.InitDb()
	// dao
	uploadDao := dao.NewOrmUpload(orm)
	// repository
	uploadRepo := repository.NewUploadRepository(uploadDao)
	// service
	uploadService := service.NewUploadService(uploadRepo)
	uploadHandler := v1.NewHandler(uploadService, hdl)
	uploadHandler.RegisterUploadRoutes(server)

	server.Run(":8889")
}
