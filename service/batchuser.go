package service

import (
	"context"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	"mall/serializer"
)

type BatchUsersService struct {
	Users []BatchUserService `json:"users" binding:"required"`
}
type BatchUserService struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行验证
}

func (service BatchUsersService) BatchRegister(ctx context.Context) serializer.Response {
	var user []model.BatchUser
	var users = service.Users
	code := e.Success
	//先进行批处理用户名唯一性校验
	userNames := make([]string, len(users))
	for i, v := range users {
		userNames[i] = v.UserName
		//	//进行校验密码
		//	if service[i].Key == "" || len(service[i].Key) != 16 {
		//		code = e.Error
		//		return serializer.Response{
		//			Status: code,
		//			Msg:    e.GetMsg(code),
		//			Data:   "密钥长度不足",
		//		}
		//	}
		//
		//	//10000  ----->密文存储,对称加密操作
		//	util.Encrypt.SetKey(service[i].Key)
		//	user[i] = model.BatchUser{
		//		UserName: service[i].UserName,
		//		NickName: service[i].NickName,
		//		Status:   model.BatchActive,
		//		Avatar:   "avatar.jpeg",
		//		Money:    util.Encrypt.AesEncoding("10000"), // 初始金额
		//	}
		//
		//	// 加密密码
		//	//前端传入的是明文
		//	if err := user[i].BatchSetPassword(service[i].Password); err != nil {
		//		code = e.ErrorFailEncryption
		//		return serializer.Response{
		//			Status: code,
		//			Msg:    e.GetMsg(code),
		//		}
		//	}
	}
	userDao := dao.BatchNewUserDao(ctx)
	_, exist, err := userDao.BatchExistOrNotByUserNames(userNames)
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

	// 创建用户
	err = userDao.BatchCreateUsers(&user) //传入指针,执行效率更高
	if err != nil {
		code = e.Error
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
