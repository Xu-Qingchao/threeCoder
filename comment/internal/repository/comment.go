package repository

import (
	"comment/internal/service"
	"comment/pkg/util"
)

type Comment struct {
	Id         uint   `gorm:"primarykey"`
	UserId     uint   `gorm:"index"`
	VideoId    uint   `gorm:"index"`
	Content    string `gorm:"type:longtext"`
	CreateTime int64
	Status     int `gorm:"default:0"` // 0正常 1删除
}

func (*Comment) CreateComment(req *service.CommentRequest) (id uint32, err error) {
	comment := Comment{
		UserId:     uint(req.UserId),
		VideoId:    uint(req.VideoId),
		Content:    req.Content,
		CreateTime: int64(req.CreateTime),
	}
	if err = DB.Create(&comment).Error; err != nil {
		util.LogrusObj.Error("Insert Video Error:" + err.Error())
		return 0, err
	}
	return uint32(comment.Id), nil
}

func (*Comment) DeleteComment(req *service.CommentRequest) error {
	c := Comment{}
	err := DB.Where("id=?", req.Id).First(&c).Error
	if err != nil {
		return err
	}
	c.Status = 1
	err = DB.Save(&c).Error
	return err
}

func (*Comment) SelectCommentsByVideo(req *service.CommentRequest) (commentList []Comment, err error) {
	err = DB.Model(Comment{}).Where("video_id=? and status = 0", req.VideoId).Find(&commentList).Error
	if err != nil {
		return commentList, err
	}
	return commentList, nil
}

func BuildComments(item []Comment) (vList []*service.CommentModel) {
	for _, v := range item {
		f := BuildComment(v)
		vList = append(vList, f)
	}
	return vList
}

func BuildComment(item Comment) *service.CommentModel {
	return &service.CommentModel{
		Id:         uint32(item.Id),
		UserId:     uint32(item.UserId),
		VideoId:    uint32(item.VideoId),
		Content:    item.Content,
		CreateTime: uint32(item.CreateTime),
	}
}
