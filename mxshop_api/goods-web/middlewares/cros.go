package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置CORS头部信息
		c.Header("Access-Control-Allow-Origin", "*") // 允许所有来源
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin")

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent) // 返回204 No Content
			return
		}

		c.Next() // 继续处理请求
	}
}
