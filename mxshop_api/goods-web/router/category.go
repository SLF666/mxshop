package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods-web/api/category"
	"mxshop_api/goods-web/middlewares"
)

func InitCategoryRouter(Router *gin.RouterGroup, jwt *middlewares.JWTManager) {
	//路由组
	CategoryRouter := Router.Group("category")
	{
		CategoryRouter.GET("list", category.List)                                                      //分类列表
		CategoryRouter.POST("", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), category.New)          //此接口需要管理员权限，增加分类
		CategoryRouter.GET("/:id", category.Detail)                                                    //分类详情
		CategoryRouter.DELETE("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), category.Delete) //删除分类
		CategoryRouter.PUT("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), category.Update)    //更新分类
	}
}
