package utils

import (
	"github.com/asaskevich/govalidator"
)

func HasLower(s string) bool {

	for _, c := range s {
		if 'a' <= c && c <= 'z' {
			return true
		}
	}
	return false
}

func HasUpper(s string) bool {

	for _, c := range s {
		if 'A' <= c && c <= 'Z' {
			return true
		}
	}
	return false
}

func HasShuzi(s string) bool {

	for _, c := range s {
		if '0' <= c && c <= '9' {
			return true
		}
	}
	return false
}

func ValidatePassPolicy(str string) {

	//自定义验证函数是否包含大小写和数字 并且长度不小于8
	govalidator.TagMap["PassPolicy"] = govalidator.Validator(func(str string) bool {
		if HasLower(str) && HasUpper(str) && HasShuzi(str) && len(str) >= 8 {
			return true
		}
		return false
	})
}

func ValidateUserActive(str string) {

	//自定义验证用户激活状态是T或者F

}

func ValidateUserInput() {
	govalidator.TagMap["PassPolicy"] = govalidator.Validator(func(str string) bool {
		if HasLower(str) && HasUpper(str) && HasShuzi(str) && len(str) >= 8 {
			return true
		}
		return false
	})
	govalidator.TagMap["ActiveValueValidate"] = govalidator.Validator(func(str string) bool {
		if str == "T" || str == "F" {
			return true
		}
		return false
	})
	//govalidator.TagMap["mutil_emailUser"] = govalidator.Validator(func(str string) bool {
	//	for _, v := range strings.Split(str, ";") {
	//		isValid := govalidator.IsEmail(v)
	//		if isValid == false {
	//			return false
	//		}
	//	}
	//	return true
	//})

}







