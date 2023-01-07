package cValidator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var CustomVerify = map[string]validator.Func{}

func VerifyEmailFormat(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyMobileFormat(fl validator.FieldLevel) bool {
	mobileNum := fl.Field().String()
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|195|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func VerifyUrlFormat(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	regular := "(http|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(url)
}

// 任意数字字母长度为6-12
func VerifyPwdFormat(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	reg := regexp.MustCompile(`^[a-zA-Z0-9]{6,12}$`)
	return reg.MatchString(password)
}

// 检测用户名称
func VerifyNameFormat(fl validator.FieldLevel) bool {
	return VerifyMobileFormat(fl) || VerifyEmailFormat(fl)
}

func init() {
	CustomVerify["VerifyEmail"] = VerifyEmailFormat
	CustomVerify["VerifyMobile"] = VerifyMobileFormat
	CustomVerify["VerifyPassword"] = VerifyPwdFormat
	CustomVerify["VerifyName"] = VerifyNameFormat
}
