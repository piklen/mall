package service

import (
	"context"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/util"
	"mall/serializer"
	"mime/multipart"
)

type UserService struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行验证
}

func (service UserService) Register(ctx context.Context) serializer.Response {
	var user model.User
	code := e.Success
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密钥长度不足",
		}
	}

	//10000  ----->密文存储,对称加密操作
	util.Encrypt.SetKey(service.Key)
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user = model.User{
		UserName: service.UserName,
		NickName: service.NickName,
		Status:   model.Active,
		Avatar:   "avatar.jpeg",
		Money:    util.Encrypt.AesEncoding("10000"), // 初始金额
	}
	// 加密密码
	//前端传入的是明文
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 创建用户
	err = userDao.CreateUser(&user) //传入指针,执行效率更高
	if err != nil {
		code = e.Error
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Login 用户登陆函数
func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	//判断用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if !exist || err != nil { // 如果查询不到，返回相应的错误
		code = e.ErrorExistUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在,请先注册!!!",
		}
	}
	//校验密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误,请重新输入密码!!!",
		}
	}
	//http 无状态(服务器需要token来认证)
	//Token签发
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "Token验证失败!!",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    e.GetMsg(code),
	}
}

// Update 用户修改信息
func (service UserService) Update(ctx context.Context, uId uint) serializer.Response {
	var user *model.User
	var err error
	code := e.Success
	// 找到用户
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	//修改昵称Nickname
	if service.NickName != "" {
		user.NickName = service.NickName
	}

	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Post 头像更新
func (service *UserService) Post(ctx context.Context, uId uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.Success
	var user *model.User
	var err error

	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//保存图片到本地
	path, err := UploadAvatarToLocalStatic(file, uId, user.UserName)

	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	user.Avatar = path
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

//var UserSrvIns *UserSrv
//var UserSrvOnce sync.Once
//
//type UserSrv struct {
//}
//
//func GetUserSrv() *UserSrv {
//	UserSrvOnce.Do(func() {
//		UserSrvIns = &UserSrv{}
//	})
//	return UserSrvIns
//}
