package service

import (
	"io/ioutil"
	"mall/conf"
	"mime/multipart"
	"os"
	"strconv"
)

// UploadAvatarToLocalStatic 上传头像
func UploadAvatarToLocalStatic(file multipart.File, userId uint, userName string) (filePath string, err error) {
	bId := strconv.Itoa(int(userId)) //路径拼接
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userName + ".jpg" //todo:把 file 的后缀提取出来
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return "", err
	}
	return "user" + bId + "/" + userName + ".jpg", err
}

// DirExistOrNot 判断文件夹路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 7550)
	if err != nil {
		return false
	}
	return true
}
