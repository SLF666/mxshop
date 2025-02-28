package router

import (
	"mxshop_api/goods-web/api/banner"
	"mxshop_api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup, jwt *middlewares.JWTManager) {
	//路由组
	BannersRouter := Router.Group("banners")
	{
		BannersRouter.GET("list", banners.List)                                                      //轮播图列表
		BannersRouter.POST("", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), banners.New)          //此接口需要管理员权限，增加轮播图
		BannersRouter.DELETE("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), banners.Delete) //删除轮播图
		BannersRouter.PUT("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), banners.Update)    //更新轮播图
	}
}
