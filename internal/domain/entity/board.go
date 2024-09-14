package entity

import "time"

// BoardMessage 表示留言板消息
type BoardMessage struct {
	ID        string
	UserQQ    string // 留言板所有者的QQ
	SenderQQ  string
	Content   string
	Timestamp time.Time
	TimeText  string
}
