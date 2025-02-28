package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/middlewares"
	router "mxshop_api/user-web/router"
)

/*
不直接在router中生成路由，是因为有多个文件，不能全部去初始化一个Group
从这里调用router目录下的文件，注册所有的路由
*/
func Routers() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.Use(middlewares.CorsMiddleware())
	v1 := r.Group("/v1")
	jwtManager := middlewares.NewJWTManager()
	router.InitUserRouter(v1, jwtManager)
	router.InitBaseRouter(v1)

	return r
}
