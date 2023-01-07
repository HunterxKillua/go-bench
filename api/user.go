package api

import (
	"fmt"
	models "ginBlog/model"
	res "ginBlog/response"

	md "ginBlog/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dict = []string{"sex", "user_roll", "user_name", "created_time", "updated_time"}

/* 初始化路由 */
func init() {
	Groups = append(Groups, func(e *gin.Engine) {
		group := e.Group("user")
		group.Use(md.Authorization())
		users := []GroupRoute{
			GetAllUser,
			GetUser,
			AddUser,
		}
		for _, route := range users {
			route(group)
		}
	})
	models.Paths["/user/list"] = "query"
}

/* 查询所有用户 分页带查询 */
func GetAllUser(g *gin.RouterGroup) {
	g.GET("/list", func(ctx *gin.Context) {
		var current int
		var sort = "user_id ASC"
		var user = []map[string]any{}
		params, _ := ctx.Get("requestConfig")
		fmt.Println(params)
		config, _ := params.(map[string]any)
		page := res.AnyToInt(config["page"])
		size := res.AnyToInt(config["size"])
		name := res.AnyToString(config["name"])
		created := res.AnyToString(config["created_time"])
		updated := res.AnyToString(config["updated_time"])
		deleted := res.AnyToBool(config["deleted"])
		if created != "" {
			sort = created
		}
		if updated != "" {
			sort = updated
		}
		current = (page - 1) * size
		var model *gorm.DB
		if deleted {
			model = DB.GetDB().Unscoped().Model(&models.User{}).Where("deleted")
		} else {
			model = DB.GetDB().Model(&models.User{})
		}
		var count int64
		err := model.
			Where("user_name like ?", "%"+name+"%").
			Or("user_id = ?", name).
			Order(sort)
		err.Select("user_id", dict).Limit(size).Offset(current).Find(&user)
		err.Count(&count)
		res.SetCallBack(err.Error, ctx, user,
			res.WithInstallRecord(true),
			res.WithInstallCount(int(count)),
			res.WithInstallPage(page),
			res.WithInstallSize(size),
		)
	})
}

/* 获取单个用户 */
func GetUser(g *gin.RouterGroup) {
	g.GET("/:id", func(ctx *gin.Context) {
		var user map[string]any
		id := res.StringToInt(ctx.Param("id"))
		if id != 0 {
			err := DB.GetDB().Model(&models.User{}).
				Select("user_id", dict).
				Where("user_id = ?", id).
				Find(&user)
			res.SetCallBack(err.Error, ctx, user, res.WithInstallRecord(true))
		} else {
			res.SetParamsError(ctx, user, "用户id不存在")
		}
	})
}

/* 新增用户 */
func AddUser(g *gin.RouterGroup) {
	g.POST("/new", func(ctx *gin.Context) {
		forms, _ := ctx.Get("interceptForm")
		config, _ := forms.(map[string]any)
		user := models.User{}
		fmt.Println(config)
		resp := DB.GetDB().Create(&user)
		res.SetCallBack(resp.Error, ctx, user, res.WithInstallRecord(true))

	})
}
