package api

import (
	"github.com/NotFound1911/filestore/api/v1"
	"github.com/gin-gonic/gin"
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
		group.GET("/upload", v1.UploadView)
		group.POST("/upload", v1.UploadFile)
		group.GET("/upload/successful", v1.UploadSuccessful)

	}
	return group
}
