package handler

import (
	"api-gateway/internal/service"
	"api-gateway/middleware/ffmpeg"
	"api-gateway/pkg/e"
	"api-gateway/pkg/res"
	"api-gateway/pkg/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
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

func PublishList(ginCtx *gin.Context) {
	//token := ginCtx.Query("token")    // 当前用户
	//// token用来鉴权
	//claim, _ := util.ParseToken(token)
	//cUserId := claim.UserID

	// 目标用户id
	userId := ginCtx.Query("user_id")
	tUserId, _ := strconv.ParseUint(userId, 10, 32)

	videoReq := service.VideoRequest{AuthorId: uint32(tUserId)}
	// gin.Key 中取出服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	videoService := ginCtx.Keys["video"].(service.VideoServiceClient)
	videoResp, err := videoService.FindVideosByUser(context.Background(), &videoReq)
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	// 遍历结果封装
	resList := videoResp.VideoDetail
	length := len(resList)
	var videoList []Video
	// 获取用户信息
	userReq := service.UserRequest{Id: uint32(tUserId)}
	userResp, _ := userService.FindUser(context.Background(), &userReq)
	for i := 0; i < length; i++ {
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

func PublishVideo(ginCtx *gin.Context) {
	token := ginCtx.PostForm("token") // 当前用户
	// token用来鉴权
	claim, _ := util.ParseToken(token)
	cUserId := claim.UserID

	// 获取上传的数据
	title := ginCtx.PostForm("title")
	data, err := ginCtx.FormFile("data")
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 存放视频，并生成路径
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", cUserId, filename)
	saveFile := filepath.Join("./static/", finalName)
	if err := ginCtx.SaveUploadedFile(data, saveFile); err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//http://127.0.0.1:4000/douyin/static/bear.mp4
	playUrl := "http://" + viper.GetString("server.host") + viper.GetString("server.port") + "/douyin/static/" + finalName
	// 截取封面并生成路径
	resp := ffmpeg.GetIpcScreenShot(
		viper.GetString("ffmpeg"),
		viper.GetString("staticUrl")+finalName,
		viper.GetString("staticUrl")+finalName+".jpg")
	print(resp)
	coverUrl := "http://" + viper.GetString("server.host") + viper.GetString("server.port") + "/douyin/static/" + finalName + ".jpg"
	videoReq := service.VideoRequest{
		AuthorId:   uint32(cUserId),
		Title:      title,
		PlayUrl:    playUrl,
		CoverUrl:   coverUrl,
		CreateTime: uint32(time.Now().Unix()),
	}
	videoService := ginCtx.Keys["video"].(service.VideoServiceClient)
	videoResp, err := videoService.CreateVideo(context.Background(), &videoReq)
	if err := ginCtx.SaveUploadedFile(data, saveFile); err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if videoResp.Code == e.ERROR {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1,
			StatusMsg: "发布失败,请联系管理员",
		})
		return
	}
	ginCtx.JSON(http.StatusOK, res.Response{
		StatusCode: 0,
		StatusMsg:  "发布成功666",
	})
}
