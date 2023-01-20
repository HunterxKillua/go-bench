package response

import (
	models "ginBlog/model"
	utils "ginBlog/util"
	cValidator "ginBlog/valid"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type InterceptConfig struct {
	Errs  []string
	Query map[string]string
	Form  map[string]string
	*cValidator.ValidConfig
}

var intercepts = map[string]*InterceptConfig{}
var PathMap = map[string]string{}

func init() {
	for key, model := range models.RequestConfig {
		intercepts[key] = NewIntercept(model)
	}
}

func NewIntercept(dest any) *InterceptConfig {
	return &InterceptConfig{
		Errs:        []string{},
		Query:       utils.GetQueryByModel(dest),
		Form:        utils.GetFormByModel(dest),
		ValidConfig: cValidator.Init(),
	}
}

func GetIntercept(path string) *InterceptConfig {
	name, ok := models.Paths[path]
	if ok {
		instance, oks := intercepts[name]
		if oks {
			return instance
		}
	}
	return nil
}

/* 批量设置query默认值 */
func setDefault(ctx *gin.Context, dest map[string]string) map[string]string {
	var defaultValue = map[string]string{}
	method := ctx.Request.Method
	for key, value := range dest {
		var val string
		if method == "GET" {
			val = ctx.DefaultQuery(key, value)
		} else {
			val = ctx.DefaultPostForm(key, value)
		}
		defaultValue[key] = val
	}
	return defaultValue
}

func InterceptErrors(opts ...string) (bool, string) {
	var str = []string{}
	for _, opt := range opts {
		if opt != "" {
			str = append(str, opt)
		}
	}
	if len(str) > 0 {
		return false, strings.Join(str, ",")
	} else {
		return true, ""
	}
}

func (instance *InterceptConfig) HttpRequestIntercept(ctx *gin.Context) (any, map[string]any, []string) {
	var config map[string]string
	switch ctx.Request.Method {
	case "GET":
		config = setDefault(ctx, instance.Query)
	case "POST":
		config = setDefault(ctx, instance.Form)
	}
	name := models.Paths[ctx.Request.URL.Path]
	dest := models.RequestConfig[name]
	result, errs := instance.ParseConfig(dest, config)
	return dest, result, errs
}

func (instance *InterceptConfig) ParseConfig(dest any, config map[string]string) (map[string]any, []string) {
	var errs = []string{}
	var result = map[string]any{}
	t := reflect.TypeOf(dest).Elem()
	v := reflect.ValueOf(dest).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		vf := v.Field(i)
		key := utils.ToSnakeString(field.Name)
		value := config[key]
		switch field.Type.Kind() {
		case reflect.String:
			validator, ok := field.Tag.Lookup("validate")
			message, ok := field.Tag.Lookup("errr_info")
			if ok {
				status, message := instance.ValidConfig.CreateValidatorByVar(value, key, message, validator)
				if !status {
					errs = append(errs, message)
				}
			}
			if vf.CanSet() {
				result[key] = value
				vf.SetString(value)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			regInt, err := strconv.Atoi(value)
			if IsError(err) {
				errs = append(errs, key+":"+getErrorByTag(field.Tag))
			} else {
				result[key] = regInt
				if vf.CanSet() {
					vf.SetInt(int64(regInt))
				}
			}
		case reflect.Bool:
			regBool, err := strconv.ParseBool(value)
			if IsError(err) {
				errs = append(errs, key+":"+getErrorByTag(field.Tag))
			} else {
				result[key] = regBool
				if vf.CanSet() {
					vf.SetBool(regBool)
				}
			}
		}
	}
	return result, errs
}

func getErrorByTag(tag reflect.StructTag) string {
	err_info, ok := tag.Lookup("err_info")
	if ok {
		return err_info
	} else {
		return "非法的参数类型"
	}
}
