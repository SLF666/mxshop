package initialize

import (
	"fmt"
	"log"
	"mxshop_sers/goods_srv/global"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	//dsn := "root:slf666@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", global.ServerConfig.MysqlInfo.User,
		global.ServerConfig.MysqlInfo.Password, global.ServerConfig.MysqlInfo.Host, global.ServerConfig.MysqlInfo.Port, global.ServerConfig.MysqlInfo.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	global.DB = db
}
