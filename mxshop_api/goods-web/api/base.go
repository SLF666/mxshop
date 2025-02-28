package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_api/goods-web/global"
	"net/http"
	"strings"
)

// 去除字段中的点以及前缀，只保留后面的
func removeTopStruce(fields map[string]string) map[string]string {
	rsp := make(map[string]string)
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 如果有错误，将错误翻译
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruce(errs.Translate(global.Trans)),
	})
}

// 将grpc的状态转换为http的状态，并直接返回
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}
