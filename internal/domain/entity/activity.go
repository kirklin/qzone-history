package entity

import "time"

// Activity 表示用户活动记录
type Activity struct {
	ID         string
	SenderQQ   string
	SenderName string
	SenderLink string
	ReceiverQQ string // 当前登录的用户
	Content    string
	Timestamp  time.Time
	TimeText   string
	ImageURLs  []string
	Type       ActivityType
	//RelatedID  string // 关联的Moment、Comment或BoardMessage的ID
}

// ActivityType 表示活动类型
type ActivityType int

const (
	TypeStatus ActivityType = iota
	TypeForward
	TypeLike
	TypeComment
	TypeBoardMessage
	TypeOther
)
