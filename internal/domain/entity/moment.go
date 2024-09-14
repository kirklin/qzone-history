package entity

import "time"

// Moment 表示QQ空间说说
type Moment struct {
	ID              string
	UserQQ          string
	Content         string
	Timestamp       time.Time
	TimeText        string
	ImageURLs       []string
	Likes           int
	Comments        []Comment
	Views           int
	IsDeleted       bool // 标记说说是否被删除
	IsReconstructed bool // 标记说说是否通过活动记录重建
}
