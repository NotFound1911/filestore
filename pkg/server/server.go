package server

import "github.com/gin-gonic/gin"

type Server struct {
	Addr string
	*gin.Engine
}
