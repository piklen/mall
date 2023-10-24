package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type BatchUser struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Avatar         string `gorm:"size:1000"`
	Money          string
}

const (
	BatchPassWordCost        = 12       //密码加密难度
	BatchActive       string = "active" //激活用户
)

// SetPassword 设置密码
func (user *BatchUser) BatchSetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BatchPassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *BatchUser) BatchCheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
