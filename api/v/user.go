package v

import (
	"github.com/gin-gonic/gin"
	"mall/pkg/util"
	"mall/service"
	"net/http"
)

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err) //绑定不成功返回错误
	}
}

// UserLogin 用户登陆接口
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}
func UserUpdate(c *gin.Context) {
	var userUpdateService service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdateService); err == nil {
		res := userUpdateService.Update(c.Request.Context(), claims.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

//
//func UserRegisterHandler() gin.HandlerFunc {
//
//	return func(ctx *gin.Context) {
//		userDao := dao.NewUserDao(ctx.Request.Context())
//		var req types.UserRegisterReq
//		if err := ctx.ShouldBind(&req); err != nil {
//			return
//		}
//		user := &model.User{
//			NickName: req.NickName,
//			UserName: req.UserName,
//			Status:   "active",
//			Money:    "10000",
//		}
//		userDao.CreateUser(user)
//	}
//}
