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

		//轮播图
		v.GET("carousels", api.ListCarousels)
		v.GET("GetProductView", api.GetProductView)
		v.POST("AddProductView", api.AddProductView)

		// 商品操作
		v.GET("products", api.ListProducts)    //获取商品列表
		v.POST("products", api.SearchProducts) //搜索商品
		authed := v.Group("/")                 //需要登录保护
		authed.Use(middleware.JWT())
		{
			// 用户操作
			authed.PUT("user", api.UserUpdate)               //更新昵称
			authed.POST("avatar", api.UploadAvatar)          // 上传头像
			authed.POST("user/sending-email", api.SendEmail) //发送邮件
			authed.POST("user/valid-email", api.ValidEmail)  //邮箱变更修改绑定等
			// 显示金额
			authed.POST("money", api.ShowMoney)

			//商品操作
			authed.POST("create", api.CreateProduct)
			//fmt.Println("没有进入了Create Product...")
			//authed.PUT("product/:id", api.UpdateProduct)
			//authed.DELETE("product/:id", api.DeleteProduct)
		}
	}
	return r
}
