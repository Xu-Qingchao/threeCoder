package handler

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user/internal/repository"
	"user/internal/service"
	"user/pkg/e"
)

type UserService struct {
	service.UnimplementedUserServiceServer //必须嵌入 UnimplementedUserServiceServer 以具有向前兼容的实现。
}

func NewUserService() *UserService {
	return &UserService{}
}

// UserLogin 用户登录 返回token（应该在网关层实现）
func (*UserService) UserLogin(ctx context.Context, req *service.UserRequest) (resp *service.UserDetailResponse, err error) {
	var user repository.User
	resp = new(service.UserDetailResponse)
	resp.Code = e.Success
	exit := user.CheckUserExit(req)
	if !exit {
		resp.Code = e.Error
		err = errors.New("UserName Not Exit")
		return resp, err
	}
	// 检查密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return resp, errors.New("密码错误")
	}
	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}

// UserRegister 用户注册
func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest) (resp *service.UserDetailResponse, err error) {
	var user repository.User
	resp = new(service.UserDetailResponse)
	resp.Code = e.Success
	user, err = user.UserCreate(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}

// FindUser 获取用户信息
func (*UserService) FindUser(ctx context.Context, req *service.UserRequest) (resp *service.UserDetailResponse, err error) {
	var user repository.User
	resp = new(service.UserDetailResponse)
	resp.Code = e.Success
	err = user.ShowUserInfo(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}
