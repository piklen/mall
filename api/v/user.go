package v

import (
	"fmt"
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
		c.JSON(http.StatusBadRequest, ErrorResponse(err)) //绑定不成功返回错误
		util.LogrusObj.Infoln(err)
	}
}

func BotchUserRegister(c *gin.Context) {
	var userRegister service.BatchUsersService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.BatchRegister(c.Request.Context())
		c.JSON(http.StatusOK, res)
		//
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err)) //绑定不成功返回错误
		util.LogrusObj.Infoln(err)
	}
}

// UserLogin 用户登陆接口
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
func UserUpdate(c *gin.Context) {
	var userUpdateService service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdateService); err == nil {
		res := userUpdateService.Update(c.Request.Context(), claims.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")
	fileSize := fileHeader.Size
	fmt.Println("测试.........")
	uploadAvatarService := service.UserService{}
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatarService); err == nil {
		res := uploadAvatarService.Post(c.Request.Context(), chaim.ID, file, fileSize)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

// SendEmail 发送邮件
func SendEmail(c *gin.Context) {
	var sendEmailService service.SendEmailService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendEmailService); err == nil {
		res := sendEmailService.Send(c.Request.Context(), chaim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
func ValidEmail(c *gin.Context) {
	var validEmailService service.ValidEmailService
	if err := c.ShouldBind(validEmailService); err == nil {
		res := validEmailService.Valid(c.Request.Context(), c.GetHeader("Authorization"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
func ShowMoney(c *gin.Context) {
	showMoneyService := service.ShowMoneyService{}
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&showMoneyService); err == nil {
		res := showMoneyService.Show(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
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
