/*
@Time : 2019-06-14 09:32
@Author : yangping
@File : index.go
@Desc :
*/
package tinyUrl

import (
	"github.com/gin-gonic/gin"
	"tinyUrl/common/middleware"
	"tinyUrl/handler/business"
)

func InitTinyUrlRoute(router *gin.RouterGroup) {
	router.Use(middleware.TokenAuthMiddleware())
	router.POST("/group", business.AddTinyGroup)
	router.POST("/list", business.TinyGroupList)
	router.GET("/group/list", business.GroupList)

	tinyUrl := router.Group("/url")
	UrlRoute(tinyUrl)
}
