package main

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	model "mxshop_sers/goods_srv/model"
	"os"
	"time"
)

// md5加密
func genMd5(psw string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, psw)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:slf666@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8&parseTime=True&loc=Local"

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
	db.AutoMigrate(&model.Category{},
		&model.Brands{},
		&model.GoodsCategoryBrand{},
		&model.Banner{},
		&model.Goods{})

	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("admin123", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//
	//for i := 0; i < 10; i++ {
	//	user := &model.User{
	//		NickName: fmt.Sprintf("user%d", i),
	//		Mobile:   fmt.Sprintf("1234567891%d", i),
	//		Password: newPassword,
	//	}
	//	db.Save(&user)
	//}
}
