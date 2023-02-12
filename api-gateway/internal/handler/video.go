package handler

import (
	"api-gateway/internal/service"
	"api-gateway/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type FeedResponse struct {
	res.Response
	NextTime  int64   `json:"next_time"`
	VideoList []Video `json:"video_list"`
}

type VideoInfoResponse struct {
	res.Response
	VideoList []Video `json:"video_list"`
}

func Feed(ginCtx *gin.Context) {
	//token := ginCtx.Query("token")    // 当前用户
	//// token用来鉴权
	//claim, _ := util.ParseToken(token)
	//cUserId := claim.UserID

	// 未登录状态下的feed
	latestTime := ginCtx.Query("latest_time")

	createTime, _ := strconv.ParseUint(latestTime, 10, 32)
	videoReq := service.VideoRequest{CreateTime: uint32(createTime)}
	// gin.Key 中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	print(ginCtx.Keys["user"])
	print(ginCtx.Keys["video"])
	videoService := ginCtx.Keys["video"].(service.VideoServiceClient)
	videoResp, err := videoService.FindVideosByTime(context.Background(), &videoReq)
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	// 遍历结果封装
	resList := videoResp.VideoDetail
	length := len(resList)
	var videoList []Video
	// 获取用户信息
	for i := 0; i < length; i++ {
		userReq := service.UserRequest{Id: resList[i].AuthorId}
		userResp, _ := userService.FindUser(context.Background(), &userReq)
		videoList = append(videoList, Video{
			Id: int64(resList[i].Id),
			Author: User{
				UserId:   int64(userResp.UserDetail.Id),
				UserName: userResp.UserDetail.Username,
			},
			PlayUrl:  resList[i].PlayUrl,
			CoverUrl: resList[i].CoverUrl,
			Title:    resList[i].Title,
		})
	}
	ginCtx.JSON(http.StatusOK, FeedResponse{
		Response:  res.Response{StatusCode: 0},
		NextTime:  int64(resList[length-1].CreateTime),
		VideoList: videoList,
	})
}
