package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name           string   `mapstructure:"name" json:"name"`
	Host           string   `mapstructure:"host" json:"host"`
	Port           int      `mapstructure:"port" json:"port"`
	Tags           []string `mapstructure:"tags" json:"tags"`
	GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JWTConfig      `mapstructure:"jwt" json:"jwt"`
	RedisConfig    `mapstructure:"redis" json:"redis"`
	ConsulConfig   `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
