package middleware

import "github.com/gin-gonic/gin"

func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		c.Keys["user"] = service[0]
		c.Next()
	}
}
