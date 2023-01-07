package models

var RequestConfig = map[string]any{
	"register": &Register{},
	"login":    &Login{},
	"query":    &QueryList{},
}

type Register struct {
	Password string `json:"password" validate:"VerifyPassword" err_info:"密码必须是长度6-12只能包含数字或者字母" form:"password"`
	Name     string `json:"name" validate:"VerifyName" err_info:"账号名必须是有效的手机号码或者邮箱地址" form:"name"`
}

type Login struct {
	Password string `json:"password" validate:"VerifyPassword" err_info:"密码必须是长度6-12只能包含数字或者字母" form:"password"`
	Name     string `json:"name" validate:"VerifyName" err_info:"账号名必须是有效的手机号码或者邮箱地址" form:"name"`
}

type QueryList struct {
	Page        int    `json:"page" validate:"omitempty,gte=1" err_info:"页数必须大于等于1" uri:"page" default:"1"`
	Size        int    `json:"size" validate:"omitempty,gte=1" err_info:"页码必须大于等于1" uri:"size" default:"10"`
	Name        string `json:"name" uri:"name"`
	Deleted     bool   `json:"deleted" validate:"omitempty,boolean" uri:"deleted" default:"false"`
	CreatedTime string `json:"created_time" validate:"omitempty,oneof=DESC ASC" uri:"created_time"`
	UpdatedTime string `json:"updated_time" validate:"omitempty,oneof=DESC ASC" uri:"updated_time"`
}
