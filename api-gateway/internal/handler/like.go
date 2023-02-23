package handler

import (
	"api-gateway/internal/service"
	"api-gateway/pkg/e"
	"api-gateway/pkg/res"
	"api-gateway/pkg/util"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FavoriteAction(ginCtx *gin.Context) {
	token := ginCtx.Query("token") // 当前用户
	// token用来鉴权
	claim, _ := util.ParseToken(token)
	cUserId := claim.UserID

	likeService := ginCtx.Keys["like"].(service.LikeServiceClient)

	// 获取上传的数据
	actionType := ginCtx.Query("action_type")
	videoId := ginCtx.Query("video_id")
	vId, _ := strconv.ParseUint(videoId, 10, 32)
	likeReq := service.LikeRequest{
		VideoId: uint32(vId),
		UserId:  uint32(cUserId),
	}
	if actionType == "1" { // 点赞
		likeResp, err := likeService.CreateLike(context.Background(), &likeReq)
		if err != nil || likeResp.Code == e.ERROR {
			ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: "点赞失败"})
			return
		}
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 0, StatusMsg: "点赞成功！"})
	} else if actionType == "2" { // 取消点赞
		likeResp, err := likeService.DeleteLike(context.Background(), &likeReq)
		if err != nil || likeResp.Code == e.ERROR {
			ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: "异常"})
			return
		}
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 0, StatusMsg: "取消成功！"})
	}
}

func FavoriteList(ginCtx *gin.Context) {
	//token := ginCtx.Query("token") // 当前用户
	//// token用来鉴权
	//claim, _ := util.ParseToken(token)
	//cUserId := claim.UserID

	likeService := ginCtx.Keys["like"].(service.LikeServiceClient)
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	videoService := ginCtx.Keys["video"].(service.VideoServiceClient)
	commentService := ginCtx.Keys["comment"].(service.CommentServiceClient)

	// 获取上传的数据
	userId := ginCtx.Query("user_id")
	tuserId, _ := strconv.ParseUint(userId, 10, 32)
	likeReq := service.LikeRequest{
		UserId: uint32(tuserId),
	}
	// 查询所有的videoId
	likeResp, err := likeService.FindVideoIds(context.Background(), &likeReq)
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	likeResList := likeResp.LikeDetail
	likeLength := len(likeResList)
	var videoIds []uint32
	for i := 0; i < likeLength; i++ {
		videoIds = append(videoIds, likeResList[i].VideoId)
	}
	videoReq := service.VideoRequest{Id: videoIds}
	// 查询video
	videoResp, err := videoService.FindVideosByIds(context.Background(), &videoReq)
	if err != nil {
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	// 遍历结果封装
	resList := videoResp.VideoDetail
	length := len(resList)
	var videoList []Video
	for i := 0; i < length; i++ {
		// 获取用户信息
		userReq := service.UserRequest{Id: resList[i].AuthorId}
		userResp, _ := userService.FindUser(context.Background(), &userReq)

		// 获取点赞数量和评论数量
		likeReq2 := service.LikeRequest{VideoId: resList[i].Id}
		likeResp2, _ := likeService.FindCountByVideo(context.Background(), &likeReq2)
		// 评论数量
		commentReq := service.CommentRequest{VideoId: resList[i].Id}
		commentResp, _ := commentService.FindCommentByVideo(context.Background(), &commentReq)
		commentCount := len(commentResp.CommentDetail)
		// 看是否登录
		token := ginCtx.Query("token") // 当前用户
		// token用来鉴权
		isLike := false
		if token != "" {
			claim, _ := util.ParseToken(token)
			cUserId := claim.UserID
			likeReq.UserId = uint32(cUserId)
			// 是否点赞
			islikeResp, _ := likeService.IsLike(context.Background(), &likeReq2)
			isLike = islikeResp.IsLike
		}

		videoList = append(videoList, Video{
			Id: int64(resList[i].Id),
			Author: User{
				UserId:   int64(userResp.UserDetail.Id),
				UserName: userResp.UserDetail.Username,
			},
			PlayUrl:       resList[i].PlayUrl,
			CoverUrl:      resList[i].CoverUrl,
			Title:         resList[i].Title,
			IsFavorite:    isLike,
			CommentCount:  int64(commentCount),
			FavoriteCount: likeResp2.Count,
		})
	}
	ginCtx.JSON(http.StatusOK, VideoInfoResponse{
		Response:  res.Response{StatusCode: 0},
		VideoList: videoList,
	})
}
