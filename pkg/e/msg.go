package e

var MsgFlags = map[int]string{
	Success:                    "ok",
	Error:                      "fail",
	InvalidParams:              "请求参数错误",
	ErrorExistUser:             "用户名已经存在",
	ErrorFailEncryption:        "密码加密失败",
	ErrorExistUserNotFound:     "用户不存在",
	ErrorNotCompare:            "密码错误",
	ErrorAuthToken:             "Token验证失败",
	ErrorAuthCheckTokenTimeout: "Token过期",
	ErrorUploadFile:            "文件上传失败!!!",
	ErrorSendEmail:             "邮件发送失败!!!!",
	ErrorProductImgUpload:      "图片上传错误！！！",
	ErrorDatabase:              "数据库错误",
	ErrorExistFavorite:         "已收藏该商品",
	ErrorProductExistCart:      "商品已经在购物车了，数量+1",
	ErrorProductMoreCart:       "超过最大上限",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
