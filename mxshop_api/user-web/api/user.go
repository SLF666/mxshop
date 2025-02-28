package api

import (
	"context"
	"fmt"
	"mxshop_api/user-web/middlewares"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"mxshop_api/user-web/forms"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/response"
	pb "mxshop_api/user-web/proto"
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

// 获取用户列表
func GetUserList(c *gin.Context) {

	claims, _ := c.Get("jwt_claims")
	currentUser := claims.(*middlewares.CustomClaims)
	zap.S().Infof("欢迎用户：%v\n", currentUser.ID)

	pn := c.DefaultQuery("pn", "1")
	pnInt, err := strconv.Atoi(pn)
	if err != nil {
		panic(err)
	}
	pSize := c.DefaultQuery("pSize", "10")
	pSizeInt, err := strconv.Atoi(pSize)
	if err != nil {
		panic(err)
	}

	rsp, err := global.UserSrvClien.GetUserList(context.Background(), &pb.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表页】 失败:" + err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]response.UserResponse, 0)
	for _, v := range rsp.Data {
		data := &response.UserResponse{
			Id:       v.Id,
			NickName: v.NickName,
			Birthday: time.Unix(int64(v.BirthDay), 0).Format("2006-01-02"),
			Mobile:   v.Mobile,
			Gender:   v.Gender,
		}
		result = append(result, *data)
	}
	c.JSON(http.StatusOK, result)
}

// 密码登录
func PasswordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	//连接grpc
	//userConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", global.ServerConfig.UserSrvConfig.Host,
	//	global.ServerConfig.UserSrvConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 连接 [用户服务失败]",
	//		"msg", err.Error(),
	//	)
	//}
	//userClient := pb.NewUserClient(userConn)

	if rsp, err := global.UserSrvClien.GetUserByMobile(context.Background(), &pb.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户不存在",
				})
			}
		}
	} else { //找到了，验证密码
		if pwdRsp, pwdErr := global.UserSrvClien.CheckPassWord(context.Background(), &pb.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); pwdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登录失败:" + pwdErr.Error(),
			})
		} else {
			if pwdRsp.Success {
				jwt := middlewares.NewJWTManager()
				token, err := jwt.GenerateToken(int(rsp.Id), rsp.NickName, int(rsp.Role))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				//返回当前用户信息
				c.JSON(http.StatusOK, gin.H{
					"id":        rsp.Id,
					"nick_name": rsp.NickName,
					"token":     "Bearer " + token,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败，密码错误",
				})
			}
		}
	}
}

// 注册用户
func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", global.ServerConfig.UserSrvConfig.Host,
		global.ServerConfig.UserSrvConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]",
			"msg", err.Error(),
		)
	}
	userClient := pb.NewUserClient(userConn)

	userInfo, err := userClient.CreateUser(context.Background(), &pb.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[register] 注册用户失败：%v\n", err.Error())
		HandleValidatorError(c, err)
		return
	}

	//自动登录
	jwt := middlewares.NewJWTManager()
	token, err := jwt.GenerateToken(int(userInfo.Id), userInfo.NickName, int(userInfo.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	//返回当前用户信息
	c.JSON(http.StatusOK, gin.H{
		"id":        userInfo.Id,
		"nick_name": userInfo.NickName,
		"token":     "Bearer " + token,
	})
}
