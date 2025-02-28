package global

import (
	"gorm.io/gorm"
	"mxshop_sers/goods_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.Nacos
)
