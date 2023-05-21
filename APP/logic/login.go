package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myVote/APP/model"
	"myVote/APP/tools"
)

type User struct {
	Name string `json:"name" form:"name" binding:"required"`
	Pwd  string `json:"pwd" form:"pwd" binding:"required"`
}
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login godoc
//
//	@Summary		用户登录
//	@Description	会执行用户登录操作
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	formData	string	true	"用户名"
//	@Param			pwd		formData	string	true	"密码"
//	@response		200,500	{object}	tools.HttpCode{data=Token}
//	@Router			/login [POST]
func Login(c *gin.Context) {
	var user User
	if c.ShouldBind(&user) != nil {
		c.JSON(200, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "用户信息错误",
			Data:    struct{}{},
		})
	}
	dbUser := model.GetUser(user.Name, user.Pwd)
	if dbUser.Id <= 0 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "用户信息错误",
			Data:    struct{}{},
		})
	}
	a, r, err := tools.Token.GetToken(dbUser.Id, dbUser.Name)
	fmt.Printf("aToken:%s\n", a)
	fmt.Printf("rToken%s\n", r)
	if err != nil {
		c.JSON(200, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "token生成失败,错误信息:" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.UserInfoErr,
		Message: "生成成功,正在跳转",
		Data: Token{
			AccessToken:  a,
			RefreshToken: r,
		},
	})
}

// Logout godoc
//
//	@Summary		用户退出
//	@Description	会执行用户退出操作
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@response		500,401	{object}	tools.HttpCode
//	@Router			/logout [get]
func Logout(c *gin.Context) {
	_ = model.FlushSession(c)
	c.Redirect(301, "/login")
	return
}
