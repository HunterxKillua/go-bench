package response

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	success = "200" /* 请求成功 */
	expired = "401" /* token过期 */
	forbidden = "403" /* 没有权限 */
	query   = "460" /* 参数非法 */
	fail    = "999" /* 请求失败 */
)

const (
	statusOk   = "ok"
	statusFail = "fail"
)

// 检测error
func IsError(err error) bool {
	return err != nil
}

/* 请求回调 */
func SetCallBack(err error, ctx *gin.Context, record any, opt ...FuncResponseMap) { /* 新增后是否返回新增的结构数据 */
	if IsError(err) {
		SetFail(ctx, err.Error())
	} else {
		d := NewResponseDict(opt...)
		SetOk(ctx, record, d)
	}
}

/* 请求成功 */
func SetOk(ctx *gin.Context, record any, dict ResponseDict) {
	var result = map[string]any{
		"status":  statusOk,
		"code":    success,
		"message": "操作成功",
	}
	if dict.ReturnRecord {
		result["data"] = record
	} else {
		result["data"] = nil
	}
	t := reflect.TypeOf(record)
	if t.Kind() == reflect.Slice {
		result["count"] = dict.ReturnCount
	}
	if dict.ReturnPage > 0 {
		result["page"] = dict.ReturnPage
	}
	if dict.ReturnSize > 0 {
		result["size"] = dict.ReturnSize
	}
	if dict.ReturnId != "" {
		result["id"] = dict.ReturnId
	}
	ctx.JSON(200, result)
}

/* 请求错误或者sql操作失败 */
func SetFail(ctx *gin.Context, message ...string) {
	ctx.JSON(400, gin.H{
		"message": strings.Join(message, ","),
		"status":  statusFail,
		"code":    fail,
		"data":    nil,
	})
}

/* token过期 */
func SetExpired(ctx *gin.Context, message string) {
	ctx.JSON(401, gin.H{
		"message": message,
		"status": statusFail,
		"code": expired,
		"data": nil,
	})
}

/* 没有权限 */
func SetForbidden(ctx *gin.Context, message string) {
	ctx.JSON(403, gin.H{
		"message": message,
		"status": statusFail,
		"code": expired,
		"data": nil,
	})
}

/* 参数错误返回结果 */
func SetParamsError(ctx *gin.Context, record any, message string) {
	ctx.JSON(400, gin.H{
		"message": message,
		"status":  statusFail,
		"code":    query,
		"data":    record,
	})
}
