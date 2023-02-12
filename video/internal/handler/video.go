package handler

import (
	"context"
	"video/internal/repository"
	"video/internal/service"
	"video/pkg/e"
)

type VideoService struct {
	service.UnimplementedVideoServiceServer //必须嵌入 UnimplementedUserServiceServer 以具有向前兼容的实现。
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

// FindVideosByTime 按照时间顺序查询多个视频
func (*VideoService) FindVideosByTime(ctx context.Context, req *service.VideoRequest) (resp *service.VideoDetailResponse, err error) {
	var video repository.Video
	resp = new(service.VideoDetailResponse)
	videoResp, err := video.SelectVideosByTime(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.VideoDetail = repository.BuildVideos(videoResp)
	return resp, nil
}

// FindVideosByUser 按照用户查询所有视频
func (*VideoService) FindVideosByUser(ctx context.Context, req *service.VideoRequest) (resp *service.VideoDetailResponse, err error) {
	var video repository.Video
	resp = new(service.VideoDetailResponse)
	videoResp, err := video.SelectVideosByUser(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.VideoDetail = repository.BuildVideos(videoResp)
	return resp, nil
}

// FindVideosByIds 按照ID查询所有视频
func (*VideoService) FindVideosByIds(ctx context.Context, req *service.VideoRequest) (resp *service.VideoDetailResponse, err error) {
	var video repository.Video
	resp = new(service.VideoDetailResponse)
	videoResp, err := video.SelectVideosByIds(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	resp.Code = e.SUCCESS
	resp.VideoDetail = repository.BuildVideos(videoResp)
	return resp, nil
}

// CreateVideo 创建视频
func (*VideoService) CreateVideo(ctx context.Context, req *service.VideoRequest) (resp *service.CommonResponse, err error) {
	var video repository.Video
	resp = new(service.CommonResponse)
	err = video.CreateVideo(req)
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
