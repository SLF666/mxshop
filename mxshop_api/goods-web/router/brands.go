package router

import (
	"mxshop_api/goods-web/api/brands"
	"mxshop_api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBrandsRouter(Router *gin.RouterGroup, jwt *middlewares.JWTManager) {
	//路由组
	brandsRouter := Router.Group("brands")
	{
		brandsRouter.GET("list", brands.BrandList)                                                      //品牌列表
		brandsRouter.POST("", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), brands.NewBrand)          //此接口需要管理员权限，新增品牌
		brandsRouter.DELETE("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), brands.DeleteBrand) //删除品牌
		brandsRouter.PUT("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), brands.UpdateBrand)    //更新品牌
	}

	CategoryBrandRouter := Router.Group("categorybrand")
	{
		CategoryBrandRouter.GET("list", brands.CategoryBrandList)                                                      //分类品牌列表
		CategoryBrandRouter.POST("", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), brands.NewCategoryBrand)          //此接口需要管理员权限，新增分类品牌
		CategoryBrandRouter.DELETE("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), brands.DeleteCategoryBrand) //删除品牌分类
		CategoryBrandRouter.PUT("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), brands.UpdateCategoryBrand)    //更新品牌分类
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)                                                   //通过分类id获取其下所有品牌
	}
}
