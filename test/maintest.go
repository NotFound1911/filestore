package main

import (
	"github.com/NotFound1911/filestore/api"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := api.NewRouter1()
	r.Run(":8888")
}
