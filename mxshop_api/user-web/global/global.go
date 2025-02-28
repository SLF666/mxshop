package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop_api/user-web/config"
	pb "mxshop_api/user-web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	Trans ut.Translator

	UserSrvClien pb.UserClient
)
