package router

import "github.com/gin-gonic/gin"

func Init() (router *gin.Engine) {
	router = gin.New()
	router.Use(gin.Recovery())
	blockRouter := router.Group("/block")
}
