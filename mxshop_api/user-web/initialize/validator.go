package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"mxshop_api/user-web/global"
)

// 初始化翻译
func InitTrans(locale string) (err error) {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uniT := ut.New(enT, zhT)

		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		global.Trans, ok = uniT.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(validate, global.Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(validate, global.Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(validate, global.Trans)
		}
		return
	} else {
		fmt.Printf("初始化翻译出错，bind：%v\n", err)
	}
	return
}
