package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("jwt_claims")
		userClaims := claims.(*CustomClaims)

		if userClaims.Role != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
