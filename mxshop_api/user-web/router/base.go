package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/api"
)

func InitBaseRouter(r *gin.RouterGroup) {
	BaseRouter := r.Group("base")
	{
		BaseRouter.GET("/captcha", api.GetCaptcha)
	}
}
