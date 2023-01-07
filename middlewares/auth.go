package middleware

import (
	"errors"
	"reflect"
	"time"

	res "ginBlog/response"
	utils "ginBlog/util"
	valid "ginBlog/valid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

var v = valid.Init()
var key = []byte("ginBlog")
var pages = map[string]bool{
	"user/list": true,
}

func SetToken(user map[string]any) string {
	claims := User{
		res.AnyToInt(user["id"]),
		res.AnyToString(user["name"]),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err == nil {
		return ss
	} else {
		return ""
	}
}

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString != "" {
			utils.Try(func() {
				token, err := jwt.ParseWithClaims(tokenString, &User{}, func(token *jwt.Token) (interface{}, error) {
					return key, nil
				})
				if claims, ok := token.Claims.(*User); ok && token.Valid {
					ctx.Set("self_id", claims.ID)
					ctx.Set("self_name", claims.Name)
					ctx.Next()
				} else {
					if errors.Is(err, jwt.ErrTokenExpired) {
						res.SetExpired(ctx, "token已过期, 请重新登录")
						ctx.Abort()
					}
					if errors.Is(err, jwt.ErrTokenNotValidYet) {
						res.SetExpired(ctx, "非法token")
						ctx.Abort()
					}
				}
			}, func(err any) {
				res.SetExpired(ctx, "非法token")
				ctx.Abort()
			})
		} else {
			res.SetFail(ctx, "请先登录")
			ctx.Abort()
		}
	}
}

func InterceptRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Params) > 0 { /* 动态路径交由下一个handler处理 */
			ctx.Next()
		}
		path := ctx.Request.URL.Path
		intercept := res.GetIntercept(path)
		dest, mapConfig, errs := intercept.HttpRequestIntercept(ctx)
		var message string
		var status bool
		if len(errs) > 0 {
			status, message = res.InterceptErrors(errs...)
		} else {
			vf := reflect.ValueOf(dest)
			status, message = v.CreateValidatorByStruct(vf.Elem().Interface())
		}
		if status {
			ctx.Set("requestModel", dest)
			ctx.Set("requestConfig", mapConfig)
			ctx.Next()
		} else {
			setParamsFail(ctx, message)
		}
	}
}

func setParamsFail(ctx *gin.Context, message string) {
	if _, ok := pages[ctx.Request.URL.Path]; ok {
		res.SetParamsError(ctx, map[string]any{}, message)
	} else {
		res.SetParamsError(ctx, nil, message)
	}
	ctx.Abort()
}
