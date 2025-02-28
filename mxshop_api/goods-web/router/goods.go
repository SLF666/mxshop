package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/goods-web/api/goods"
	"mxshop_api/goods-web/middlewares"
)

/*
所有的具体路由，由于有多个文件，需要传一个gin.RouterGroup进来，统一生成
*/

func InitGoodsRouter(r *gin.RouterGroup, jwt *middlewares.JWTManager) {
	GoodsRouter := r.Group("goods")
	zap.S().Info("配置用户相关的url")
	{
		//获取商品列表
		GoodsRouter.GET("", goods.List)
		//新增商品
		GoodsRouter.POST("", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), goods.New)
		//获取商品详细信息
		GoodsRouter.GET("/:id", goods.Detail)
		//删除商品
		GoodsRouter.DELETE("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), goods.Delete)
		//获取商品库存
		GoodsRouter.GET("/:id/stocks", goods.Stocks)
		//更新商品状态
		GoodsRouter.PATCH("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), goods.UpdateStatus)
		//更新商品信息
		GoodsRouter.PUT("/:id", jwt.JWTMiddleware(), middlewares.IsAdminAuth(), goods.Update)
	}
}
