package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/user-web/api"
	"mxshop_api/user-web/middlewares"
)

/*
所有的具体路由，由于有多个文件，需要传一个gin.RouterGroup进来，统一生成
*/

func InitUserRouter(r *gin.RouterGroup, jwt *middlewares.JWTManager) {
	UserRouterGroup := r.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouterGroup.GET("list", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouterGroup.POST("pwd_login", api.PasswordLogin)
		UserRouterGroup.POST("register", api.Register)
	}
}
