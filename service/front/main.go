package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const staticPath = "./service/front/static/view"

func main() {
	server := gin.Default()
	server.HTMLRender = loadTemplates(staticPath)
	registerRoutes(server)
	if err := server.Run(":8080"); err != nil {
		fmt.Println("err:", err)
	}
}
