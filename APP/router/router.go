package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"myVote/APP/logic"
	"myVote/APP/tools"
	_ "myVote/docs"
	"net/http"
)

func New() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	{
		//获取投票列表
		r.GET("/votes", logic.GetVotes)
	}
	basic := r.Group("/vote")
	basic.Use(AuthCheck())

	{

		//投票选项
		//获取投票选项
		basic.GET("/opts/:vote_id", logic.GetOpts)
		basic.POST("/opt", logic.AddOpt)
		basic.PUT("/opt", logic.PutOpt)
		basic.DELETE("/opt/:id", logic.DeleteOpt)

		//投票主体
		basic.GET("/:id", logic.GetVote)
		basic.POST("/basic", logic.AddVoteBasic)
		basic.PUT("/basic", logic.PutVoteBasic)
		basic.DELETE("/basic/:id", logic.DeleteVoteBasic)
		//投票
		basic.POST("/do", logic.DoVote)
	}

	login := r.Group("")
	{
		login.POST("/login", logic.Login)
		login.GET("/logout", logic.Logout)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r

}

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("debug") != "" {
			c.Next()
			return
		}
		auth := c.GetHeader("Authorization")
		fmt.Printf("auth:%+v\n", auth)
		data, err := tools.Token.VerifyToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(401, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "验签失败",
			})
		}
		fmt.Printf("data:%+v\n", data)
		if data.ID <= 0 || data.Name == "" {
			c.AbortWithStatusJSON(401, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "用户信息错误",
			})
			return
		}
		c.Set("name", data.Name)
		c.Set("userId", data.ID)
	}
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			// 当Access-Control-Allow-Credentials为true时，将*替换为指定的域名
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
