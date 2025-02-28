package initialize

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_sers/goods_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 读取配置文件
func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "goods_srv/config-pro.yaml"
	if debug {
		configFileName = "goods_srv/config-debug.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("Nacos配置：%v\n", global.NacosConfig)

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		// 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	zap.S().Infof("配置信息：%v\n", global.ServerConfig)
}
