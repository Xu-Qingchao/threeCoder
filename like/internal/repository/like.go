package repository

import (
	"like/internal/service"
	"like/pkg/util"
)

type Like struct {
	Id      uint `gorm:"primarykey"`
	UserId  uint `gorm:"index"`
	VideoId uint `gorm:"index"`
}

func (*Like) CreateLike(req *service.LikeRequest) error {
	like := Like{
		UserId:  uint(req.UserId),
		VideoId: uint(req.VideoId),
	}
	if err := DB.Create(&like).Error; err != nil {
		util.LogrusObj.Error("Insert Video Error:" + err.Error())
		return err
	}
	return nil
}

func (*Like) DeleteLike(req *service.LikeRequest) error {
	err := DB.Where("user_id=? and video_id=?", req.UserId, req.VideoId).Delete(Like{}).Error
	return err
}

func (*Like) SelectVideoIds(req *service.LikeRequest) (likeList []Like, err error) {
	err = DB.Model(Like{}).Where("user_id=?", req.UserId).Find(&likeList).Error
	if err != nil {
		return likeList, err
	}
	return likeList, nil
}

func (*Like) SelectCountByVideo(req *service.LikeRequest) (count int64, err error) {
	DB.Model(Like{}).Where("video_id=?", req.VideoId).Count(&count)
	return count, nil
}

func (like *Like) IsLike(req *service.LikeRequest) (isLike bool, err error) {
	err = DB.Where("user_id=? and video_id=?", req.UserId, req.VideoId).Find(&like).Error
	if err != nil {
		return false, err
	}
	if like.Id == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func BuildLikes(item []Like) (vList []*service.LikeModel) {
	for _, v := range item {
		f := BuildLike(v)
		vList = append(vList, f)
	}
	return vList
}

func BuildLike(item Like) *service.LikeModel {
	return &service.LikeModel{
		Id:      uint32(item.Id),
		UserId:  uint32(item.UserId),
		VideoId: uint32(item.VideoId),
	}
}
