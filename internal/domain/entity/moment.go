package entity

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/gorm"
	"time"
)

// Moment 表示QQ空间说说
type Moment struct {
	ID              string    `json:"id" gorm:"primaryKey"`                // 说说ID
	UserQQ          string    `json:"userQQ" gorm:"index"`                 // 发布者QQ
	Content         string    `json:"content"`                             // 说说内容
	Timestamp       time.Time `json:"timestamp"`                           // 发布时间戳
	TimeText        string    `json:"timeText"`                            // 发布时间文本
	ImageURLs       []string  `json:"imageURLs" gorm:"type:text"`          // 图片URL列表
	Likes           int       `json:"likes"`                               // 点赞数
	Comments        []Comment `json:"comments" gorm:"foreignKey:MomentID"` // 评论列表
	Views           int       `json:"views"`                               // 浏览次数
	IsDeleted       bool      `json:"isDeleted"`                           // 是否已删除
	IsReconstructed bool      `json:"isReconstructed"`                     // 是否通过活动记录重建
}

// BeforeCreate 钩子，在创建记录之前自动生成UUID（仅当ID为空时）
func (moment *Moment) BeforeCreate(tx *gorm.DB) (err error) {
	if moment.ID == "" {
		// 使用内容和收到的人的QQ生成唯一键
		key := moment.Content + moment.UserQQ
		hash := md5.Sum([]byte(key))
		moment.ID = hex.EncodeToString(hash[:])
	}
	return
}
