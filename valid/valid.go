package cValidator

import (
	"ginBlog/pkg/logger"
	utils "ginBlog/util"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

/*
len：length 等于，长度相等
max：小于等于
min：大于等于
eq：等于，字符串相等
ne：不等于
gt：大于
gte：大于等于
lt：小于
lte：小于等于，例如lte=10；
oneof：值中的一个，例如oneof=1 2

支持时间范围的比较lte
时间 RegTime time.Time `validate:"lte"` 小于等于当前时间

跨字段约束
eqfield=ConfirmPassword
eqcsfield=InnerStructField.Field

字符串规则
contains=：包含参数子串
containsany：包含参数中任意的 UNICODE 字符
containsrune：包含参数表示的 rune 字符
excludes：不包含参数子串
excludesall：不包含参数中任意的 UNICODE 字符
excludesrune：不包含参数表示的 rune 字符
startswith：以参数子串为前缀
endswith：以参数子串为后缀

使用unqiue来指定唯一性约束，对不同类型的处理如下：

对于数组和切片，unique约束没有重复的元素；
对于map，unique约束没有重复的值；
对于元素类型为结构体的切片，unique约束结构体对象的某个字段不重复，通过unqiue=name

特殊规则

-：跳过该字段，不检验；
|：使用多个约束，只需要满足其中一个，例如rgb|rgba；
required：字段必须设置，不能为默认值；
omitempty：如果字段未设置，则忽略它。
*/

type DestEnum string

const (
	row  DestEnum = "row"
	rows DestEnum = "rows"
)

type ValidConfig struct {
	Locales       ut.Translator
	ValidInstance *validator.Validate
}

// 生成一个validator实例, 语言为zh
func Init() *ValidConfig {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	for name, validFunc := range CustomVerify { /* 注册自定义验证器 */
		validate.RegisterValidation(name, validFunc)
	}
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil
	}
	return &ValidConfig{
		Locales:       trans,
		ValidInstance: validate,
	}
}

func (v *ValidConfig) CreateValidatorByVar(dest any, key string, message string, tag ...string) (bool, string) {
	err := v.ValidInstance.Var(dest, strings.Join(tag, ","))
	str := []string{}
	if err != nil {
		for range err.(validator.ValidationErrors) {
			str = append(str, message)
		}
		return false, key + ":" + strings.Join(str, ",")
	}
	return true, ""
}

func (v *ValidConfig) CreateValidatorByStruct(dest any) (bool, string) {
	err := v.ValidInstance.Struct(dest)
	status, errors := v.checkResult(dest, err)
	var errs = []string{}
	for _, message := range errors {
		errs = append(errs, message)
	}
	return status, strings.Join(errs, ",")
}

// 检测validator的验证结果
func (v *ValidConfig) checkResult(dest any, err error) (bool, map[string]string) {
	var errors = map[string]string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			field, ok := reflect.TypeOf(dest).FieldByName(fieldName)
			key := utils.ToSnakeString(fieldName)
			if ok {
				err_info := field.Tag.Get("err_info")
				if err_info == "" {
					errors[key] = replaceString(fieldName, err.Translate(v.Locales), key+":")
				} else {
					errors[key] = key + ":" + err_info
				}
			} else {
				errors[key] = replaceString(fieldName, err.Translate(v.Locales), key+":")
			}

		}
		return false, errors
	}
	return true, errors
}

// 检测目标参数是否为map或者struct(可以为切片)的指针
func (v *ValidConfig) CheckDestPtr(dest any, typ DestEnum) bool {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		logger.Infow("请传入一个指针, 而非引用类型")
		return false
	}
	switch typ {
	case row:
		if t.Elem().Kind() == reflect.Struct || t.Elem().Kind() == reflect.Map {
			return true
		} else {
			logger.Infow("期望入参为struct或者Map类型的指针")
		}
	case rows:
		if t.Elem().Kind() == reflect.Slice {
			if t.Elem().Elem().Kind() == reflect.Struct || t.Elem().Elem().Kind() == reflect.Map {
				return true
			}
		}
		logger.Infow("期望入参为struct或者Map类型的切片指针")
	default:
		return false
	}
	return false
}

/* 检测是否为struct类型的指针*/
func (v *ValidConfig) IsStructPtr(dest any) bool {
	t := reflect.TypeOf(dest)
	if v.IsPtr(dest) {
		return t.Elem().Kind() == reflect.Struct
	} else {
		logger.Infow("期望入参为struct类型的指针")
		return false
	}
}

/* 检测参数是否为指针 */
func (v *ValidConfig) IsPtr(dest any) bool {
	t := reflect.TypeOf(dest)
	return t.Kind() == reflect.Ptr
}

func replaceString(key string, dest string, repl string) string {
	reg, _ := regexp.Compile(key)
	return reg.ReplaceAllString(dest, repl)
}
