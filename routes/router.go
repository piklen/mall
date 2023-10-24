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
		v.POST("user/batchRegister", api.BotchUserRegister)
		v.POST("user/login", api.UserLogin)

		//轮播图
		v.GET("carousels", api.ListCarousels)
		v.GET("GetProductView", api.GetProductView)
		v.POST("AddProductView", api.AddProductView)

		// 商品操作
		v.GET("products", api.ListProducts)
		v.POST("products/search", api.SearchProducts)
		v.GET("product/show/", api.ShowProduct)         //展示商品信息
		v.GET("products/imgs/", api.ListProductImg)     // 商品图片
		v.GET("product/categories", api.ListCategories) //商品分类
		authed := v.Group("/")                          //需要登录保护
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
			//收藏夹操作
			authed.POST("favorites", api.CreateFavorite)
			authed.GET("favorites/show", api.ShowFavorites)
			authed.POST("favorites/", api.DeleteFavorite)

			// 收获地址操作
			authed.POST("addresses/create", api.CreateAddress)
			authed.GET("addresses/get", api.GetAddress)
			authed.GET("addresses/list", api.ListAddress) //展示全部地址
			authed.POST("addresses/update", api.UpdateAddress)
			authed.POST("addresses/del", api.DeleteAddress)

			// 购物车
			authed.POST("carts/create", api.CreateCart)
			authed.GET("carts/show", api.ShowCarts)
			authed.POST("carts/update", api.UpdateCart) //修改的主要是数量
			authed.POST("carts/del", api.DeleteCart)

		}
	}
	return r
}
