package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile); !ok {
		return false
	}
	return true
}
