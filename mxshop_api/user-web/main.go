package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
	"mxshop_api/user-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop_api/user-web/global"
	initi "mxshop_api/user-web/initialize"
	"mxshop_api/user-web/utils"
	myValidator "mxshop_api/user-web/validator"
)

func main() {
	//初始化日志
	logger := initi.InitLogger()
	defer logger.Sync()
	//初始化配置
	initi.InitConfig()
	//初始化翻译
	if err := initi.InitTrans("zh"); err != nil {
		panic(err)
	}
	//初始化连接
	initi.InitSrvConn()

	viper.AutomaticEnv()
	//如果是本地开发环境，端口固定
	//线上环境自动获取可用端口号
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myValidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法手机号!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	port := global.ServerConfig.Port

	r := initi.Routers()

	//服务注册
	srvRegisterClient := consul.NewRegisterClient(global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	serviceId, _ := uuid.NewV4()
	serviceIdstr := fmt.Sprintf("%s", serviceId)
	err := srvRegisterClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceIdstr)
	if err != nil {
		zap.S().Panic("服务注册失败：", err.Error())
	}

	zap.S().Debugf("正在启动服务器，端口：%v\n", port)
	go func() {
		if err := r.Run(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
			zap.S().Panic("启动失败", zap.Error(err))
		}
	}()

	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = srvRegisterClient.Deregister(serviceIdstr); err != nil {
		zap.S().Panic("服务注销失败", err.Error())
	} else {
		zap.S().Infof("服务注销成功")
	}
}
