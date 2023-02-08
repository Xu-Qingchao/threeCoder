package handler

import (
	"api-gateway/internal/service"
	"api-gateway/pkg/res"
	"api-gateway/pkg/util"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	res.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type User struct {
	UserId        int64  `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserInfoResponse struct {
	res.Response
	User User `json:"user"`
}

// Register 用户注册
func Register(ginCtx *gin.Context) {
	username := ginCtx.Query("username")
	password := ginCtx.Query("password")
	//var userReq service.UserRequest
	userReq := service.UserRequest{
		Username: username,
		Password: password,
	}
	//PanicIfUserError(ginCtx.Bind(&userReq))
	// gin.Key 中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserRegister(context.Background(), &userReq)
	//PanicIfUserError(err)
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		token, err := util.GenerateToken(uint(userResp.UserDetail.Id))
		if err != nil {
			ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 2, StatusMsg: "系统错误！"})
		} else {
			ginCtx.JSON(http.StatusOK, UserLoginResponse{Response: res.Response{
				StatusCode: 0},
				UserId: int64(userResp.UserDetail.Id),
				Token:  token,
			})
		}
	}
}

// Login 用户登录
func Login(ginCtx *gin.Context) {
	username := ginCtx.Query("username")
	password := ginCtx.Query("password")
	userReq := service.UserRequest{
		Username: username,
		Password: password,
	}
	//PanicIfUserError(ginCtx.Bind(&userReq))
	// gin.Key 中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserLogin(context.Background(), &userReq)
	//PanicIfUserError(err)
	// 其实返回一下err信息就好
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		token, err := util.GenerateToken(uint(userResp.UserDetail.Id))
		if err != nil {
			ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 2, StatusMsg: "系统错误"})
		} else {
			ginCtx.JSON(http.StatusOK, UserLoginResponse{
				Response: res.Response{StatusCode: 0},
				UserId:   int64(userResp.UserDetail.Id),
				Token:    token,
			})
		}
	}
}

// UserInfo 获取用户信息
func UserInfo(ginCtx *gin.Context) {
	userId := ginCtx.Query("user_id") // 目标用户
	//token := ginCtx.Query("token")    // 当前用户
	//// token用来鉴权
	//claim, _ := util.ParseToken(token)
	//cUserId := claim.UserID
	tUserId, _ := strconv.ParseUint(userId, 10, 32)
	userReq := service.UserRequest{Id: uint32(tUserId)}
	// 获取用户信息
	// gin.Key 中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.FindUser(context.Background(), &userReq)
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	// 获取目标用户关注信息
	ginCtx.JSON(http.StatusOK, UserInfoResponse{
		Response: res.Response{StatusCode: 0},
		User: User{
			UserId:   int64(userResp.UserDetail.Id),
			UserName: userResp.UserDetail.Username,
		},
	})

}
