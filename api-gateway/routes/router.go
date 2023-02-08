package routes

import (
	"api-gateway/internal/handler"
	"api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(service))
	v1 := ginRouter.Group("/douyin")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register/", handler.Register)
		v1.POST("/user/login/", handler.Login)
		v1.GET("/user/", handler.UserInfo)
	}

	return ginRouter
}
