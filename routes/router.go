package routes

import (
	"github.com/gin-gonic/gin"
	api "mall/api/v"
	"mall/middleware"
	"net/http"
)

// 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	r.StaticFS("/static", http.Dir("./static"))

	v := r.Group("api/v")
	{
		v.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(200, "success")
		})
		// 用户操作
		v.POST("user/register", api.UserRegister)
		v.POST("user/login", api.UserLogin)
		authed := v.Group("/") //需要登录保护
		authed.Use(middleware.JWT())
		{
			// 用户操作
			authed.PUT("user", api.UserUpdate)
			//authed.POST("user/sending-email", api.SendEmail)
			//authed.POST("user/valid-email", api.ValidEmail)
			//authed.POST("avatar", api.UploadAvatar) // 上传头像
		}
		//v.POST("user/register", api.UserRegisterHandler())
		//v.POST("user/login", api.UserLogin)

	}
	return r
}
