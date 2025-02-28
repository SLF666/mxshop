package initialize

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_api/user-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "user-web/config-pro.yaml"
	if debug {
		configFileName = "user-web/config-debug.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("Nacos信息：%v", global.NacosConfig)
	zap.S().Infof("namespace = %v", global.NacosConfig.Namespace)
	zap.S().Infof("group = %v", global.NacosConfig.Group)
	zap.S().Infof("dataid = %v", global.NacosConfig.DataId)

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,         // Nacos的服务地址
			Port:   uint64(global.NacosConfig.Port), // Nacos的服务端口
		},
	}

	// 创建clientConfig
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	//获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	//zap.S().Infof("content: %v", content)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(content), &global.ServerConfig); err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)
}
