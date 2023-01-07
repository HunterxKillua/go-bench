package api

import (
	"crypto/md5"
	"fmt"
	middleware "ginBlog/middlewares"
	models "ginBlog/model"

	// res "ginBlog/response"
	"github.com/gin-gonic/gin"
)

func init() {
	Routes = append(Routes, Login, Register)
	models.Paths["/login"] = "login"
	models.Paths["/register"] = "register"
}

func Login(r *gin.Engine) {
	r.POST("/login", func(ctx *gin.Context) {
		str := "phone"
		data := []byte(str)
		has := md5.Sum(data)
		md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
		fmt.Println(md5str1)
		// var user map[string]any
		// err := db.GetDB().Model(&models.User{}).
		// 	Select("user_id", dict).
		// 	Where("user_id = ?", 4).
		// 	Find(&user)

		token := middleware.SetToken(map[string]any{
			"name": "hah",
			"id":   1,
		})
		ctx.JSON(200, gin.H{
			"token": token,
		})
		// res.SetCallBack(err.Error, ctx, token, res.WithInstallRecord(true))
	})
}

func Register(r *gin.Engine) {
	r.POST("/register", func(ctx *gin.Context) {
		requestModel, _ := ctx.Get("requestModel")
		requestConfig, _ := ctx.Get("requestConfig")
		config, _ := requestConfig.(map[string]any)
		fmt.Println(requestModel)
		ctx.JSON(200, gin.H{
			"message": nil,
			"config":  config,
		})
		// forms, _ := ctx.Get("interceptForm")
		// config, _ := forms.(map[string]any)
		// fmt.Println(config)
	})
}
