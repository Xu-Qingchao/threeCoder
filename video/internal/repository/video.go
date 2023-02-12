package repository

import (
	"video/internal/service"
	"video/pkg/util"
)

type Video struct {
	Id         uint `gorm:"primarykey"`
	AuthorId   uint `gorm:"index"`
	Title      string
	PlayUrl    string
	CoverUrl   string
	CreateTime int64
}

func (*Video) CreateVideo(req *service.VideoRequest) error {
	video := Video{
		AuthorId:   uint(req.AuthorId),
		Title:      req.Title,
		PlayUrl:    req.PlayUrl,
		CoverUrl:   req.CoverUrl,
		CreateTime: int64(req.CreateTime),
	}
	if err := DB.Create(&video).Error; err != nil {
		util.LogrusObj.Error("Insert Video Error:" + err.Error())
		return err
	}
	return nil
}

func (*Video) SelectVideosByTime(req *service.VideoRequest) (videoList []Video, err error) {
	err = DB.Model(Video{}).Where("create_time<=?", req.CreateTime).Limit(30).Order("create_time desc").Find(&videoList).Error
	if err != nil {
		return videoList, err
	}
	return videoList, nil
}

func (*Video) SelectVideosByUser(req *service.VideoRequest) (videoList []Video, err error) {
	err = DB.Model(Video{}).Where("author_id=?", req.AuthorId).Find(&videoList).Error
	if err != nil {
		return videoList, err
	}
	return videoList, nil
}

func (*Video) SelectVideosByIds(req *service.VideoRequest) (videoList []Video, err error) {
	err = DB.Model(Video{}).Where("id=?", req.Id).Find(&videoList).Error
	if err != nil {
		return videoList, err
	}
	return videoList, nil
}

func BuildVideos(item []Video) (vList []*service.VideoModel) {
	for _, v := range item {
		f := BuildVideo(v)
		vList = append(vList, f)
	}
	return vList
}

func BuildVideo(item Video) *service.VideoModel {
	return &service.VideoModel{
		Id:         uint32(item.Id),
		AuthorId:   uint32(item.AuthorId),
		Title:      item.Title,
		PlayUrl:    item.PlayUrl,
		CoverUrl:   item.CoverUrl,
		CreateTime: uint32(item.CreateTime),
	}
}
