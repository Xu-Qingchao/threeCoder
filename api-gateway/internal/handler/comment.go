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
	"time"
)

type Comment struct {
	Id         int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type CommentInfoResponse struct {
	res.Response
	CommentList []Comment `json:"comment_list"`
}

type CommentActionResponse struct {
	res.Response
	Comment Comment `json:"comment"`
}

func CommentAction(ginCtx *gin.Context) {
	token := ginCtx.Query("token") // 当前用户
	// token用来鉴权
	claim, _ := util.ParseToken(token)
	cUserId := claim.UserID

	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	commentService := ginCtx.Keys["comment"].(service.CommentServiceClient)

	// 获取上传的数据
	actionType := ginCtx.Query("action_type")
	if actionType == "1" {
		content := ginCtx.Query("comment_text")
		videoId := ginCtx.Query("video_id")
		vId, _ := strconv.ParseUint(videoId, 10, 32)
		commentReq := service.CommentRequest{
			VideoId:    uint32(vId),
			UserId:     uint32(cUserId),
			Content:    content,
			CreateTime: uint32(time.Now().Unix()),
		}
		commentResp, err := commentService.CreateComment(context.Background(), &commentReq)
		if err != nil || commentResp.Code == e.ERROR {
			ginCtx.JSON(http.StatusOK, CommentActionResponse{
				Response: res.Response{StatusCode: 1, StatusMsg: "发布失败"},
			})
			return
		}
		// 查出用户信息
		userReq := service.UserRequest{Id: uint32(cUserId)}
		userResp, _ := userService.FindUser(context.Background(), &userReq)
		ginCtx.JSON(http.StatusOK, CommentActionResponse{
			Response: res.Response{StatusCode: 0},
			Comment: Comment{
				Id:         int64(commentResp.Id),
				User:       User{UserId: int64(userResp.UserDetail.Id), UserName: userResp.UserDetail.Username},
				Content:    content,
				CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			},
		})
	} else if actionType == "2" {
		commentId := ginCtx.Query("comment_id")
		cId, _ := strconv.ParseUint(commentId, 10, 32)
		commentReq := service.CommentRequest{Id: uint32(cId)}
		commentResp, err := commentService.DeleteComment(context.Background(), &commentReq)
		if err != nil || commentResp.Code == e.ERROR {
			ginCtx.JSON(http.StatusOK, CommentActionResponse{
				Response: res.Response{StatusCode: 1, StatusMsg: "删除失败"},
			})
			return
		}
		ginCtx.JSON(http.StatusOK, res.Response{StatusCode: 0})
	}
}

func CommentList(ginCtx *gin.Context) {
	//token := ginCtx.Query("token") // 当前用户
	//// token用来鉴权
	//claim, _ := util.ParseToken(token)
	//cUserId := claim.UserID

	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	commentService := ginCtx.Keys["comment"].(service.CommentServiceClient)

	// 获取上传的数据
	videoId := ginCtx.Query("video_id")
	vId, _ := strconv.ParseUint(videoId, 10, 32)
	commentReq := service.CommentRequest{VideoId: uint32(vId)}
	commentResp, err := commentService.FindCommentByVideo(context.Background(), &commentReq)
	if err != nil {
		ginCtx.JSON(http.StatusOK, CommentActionResponse{
			Response: res.Response{StatusCode: 1, StatusMsg: "获取评论失败"},
		})
		return
	}
	// 遍历结果封装
	resList := commentResp.CommentDetail
	length := len(resList)
	var commentList []Comment
	// 获取用户信息
	for i := 0; i < length; i++ {
		userReq := service.UserRequest{Id: resList[i].UserId}
		userResp, _ := userService.FindUser(context.Background(), &userReq)
		commentList = append(commentList, Comment{
			Id: int64(resList[i].Id),
			User: User{
				UserId:   int64(userResp.UserDetail.Id),
				UserName: userResp.UserDetail.Username,
			},
			Content:    resList[i].Content,
			CreateDate: time.Unix(int64(resList[i].CreateTime), 0).Format("2006-01-02 15:04:05"),
		})
	}
	ginCtx.JSON(http.StatusOK, CommentInfoResponse{
		Response:    res.Response{StatusCode: 0},
		CommentList: commentList,
	})
}
