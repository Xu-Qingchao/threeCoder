package repository

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"user/internal/service"
)

type User struct {
	Id       uint   `gorm:"primarykey"`
	Username string `gorm:"unique"`
	Password string
}

const (
	PasswordCost = 12 // 密码加密难度
)

// CheckUserExit 检查用户是否存在
func (user *User) CheckUserExit(req *service.UserRequest) bool {
	if err := DB.Where("username=?", req.Username).First(&user).Error; err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

// ShowUserInfo 根据id获取用户信息
func (user *User) ShowUserInfo(req *service.UserRequest) error {
	err := DB.Where("id=?", req.Id).First(&user).Error
	return err
}

// UserLogin 用户登录

// UserCreate 创建用户
func (*User) UserCreate(req *service.UserRequest) (user User, err error) {
	var count int64
	DB.Where("username=?", req.Username).Count(&count)
	if count != 0 {
		return User{}, errors.New("UserName Exited")
	}
	user = User{
		Username: req.Username,
	}
	// 密码加密
	_ = user.SetPassword(req.Password)

	err = DB.Create(&user).Error
	return user, err
}

// SetPassword 密码加密
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// BuildUser 序列话user
func BuildUser(item User) *service.UserModel {
	userModel := service.UserModel{
		Id:       uint32(item.Id),
		Username: item.Username,
		Password: item.Password,
	}
	return &userModel
}
