package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop_api/goods-web/config"
	pb "mxshop_api/goods-web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	Trans ut.Translator

	GoodsSrvClient pb.GoodsClient
)
