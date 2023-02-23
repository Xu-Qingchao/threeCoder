package handler

import (
	"context"
	"like/internal/repository"
	"like/internal/service"
	"like/pkg/e"
)

type LikeService struct {
	service.UnimplementedLikeServiceServer //必须嵌入 UnimplementedUserServiceServer 以具有向前兼容的实现。
}

func NewLikeService() *LikeService {
	return &LikeService{}
}

// CreateLike 新建喜欢
func (*LikeService) CreateLike(ctx context.Context, req *service.LikeRequest) (resp *service.LikeCommonResponse, err error) {
	var like repository.Like
	resp = new(service.LikeCommonResponse)
	err = like.CreateLike(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(uint(resp.Code))
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

// DeleteLike 取消点赞
func (*LikeService) DeleteLike(ctx context.Context, req *service.LikeRequest) (resp *service.LikeCommonResponse, err error) {
	var like repository.Like
	resp = new(service.LikeCommonResponse)
	err = like.DeleteLike(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(uint(resp.Code))
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

// FindVideoIds 根据用户id找到所有喜欢
func (*LikeService) FindVideoIds(ctx context.Context, req *service.LikeRequest) (resp *service.LikeDetailResponse, err error) {
	var like repository.Like
	resp = new(service.LikeDetailResponse)
	likeResp, err := like.SelectVideoIds(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.LikeDetail = repository.BuildLikes(likeResp)
	return resp, nil
}

// IsLike 判断用户id是否对video点赞
func (*LikeService) IsLike(ctx context.Context, req *service.LikeRequest) (resp *service.LikeCommonResponse, err error) {
	var like repository.Like
	resp = new(service.LikeCommonResponse)
	isLike, err := like.IsLike(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(uint(resp.Code))
		resp.IsLike = false
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(uint(resp.Code))
	resp.IsLike = isLike
	return resp, nil
}

// 统计视频点赞数量
func (*LikeService) FindCountByVideo(ctx context.Context, req *service.LikeRequest) (resp *service.LikeCommonResponse, err error) {
	var like repository.Like
	resp = new(service.LikeCommonResponse)
	count, err := like.SelectCountByVideo(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(uint(resp.Code))
		resp.Count = 0
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(uint(resp.Code))
	resp.Count = count
	return resp, nil
}
