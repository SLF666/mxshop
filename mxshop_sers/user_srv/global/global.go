package global

import (
	"gorm.io/gorm"
	"mxshop_sers/user_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.Nacos
)
