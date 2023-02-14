package handler

import (
	"comment/internal/repository"
	"comment/internal/service"
	"comment/pkg/e"
	"context"
)

type CommentService struct {
	service.UnimplementedCommentServiceServer //必须嵌入 UnimplementedUserServiceServer 以具有向前兼容的实现。
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

// CreateComment 新增评论
func (*CommentService) CreateComment(ctx context.Context, req *service.CommentRequest) (resp *service.CommentCommonResponse, err error) {
	var comment repository.Comment
	resp = new(service.CommentCommonResponse)
	id, err := comment.CreateComment(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(uint(resp.Code))
	resp.Id = id
	return resp, nil
}

// DeleteComment 删除评论
func (*CommentService) DeleteComment(ctx context.Context, req *service.CommentRequest) (resp *service.CommentCommonResponse, err error) {
	var comment repository.Comment
	resp = new(service.CommentCommonResponse)
	err = comment.DeleteComment(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

// FindCommentByVideo 查询所有评论
func (*CommentService) FindCommentByVideo(ctx context.Context, req *service.CommentRequest) (resp *service.CommentDetailResponse, err error) {
	var comment repository.Comment
	resp = new(service.CommentDetailResponse)
	commentResp, err := comment.SelectCommentsByVideo(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.CommentDetail = repository.BuildComments(commentResp)
	return resp, nil
}
