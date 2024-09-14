package entity

import "time"

// BoardMessage 表示留言板消息
type BoardMessage struct {
	ID        string    `json:"id" gorm:"primaryKey"`  // 留言ID
	UserQQ    string    `json:"userQQ" gorm:"index"`   // 留言板所有者QQ
	SenderQQ  string    `json:"senderQQ" gorm:"index"` // 发送者QQ
	Content   string    `json:"content"`               // 留言内容
	Timestamp time.Time `json:"timestamp"`             // 时间戳
	TimeText  string    `json:"timeText"`              // 时间文本
}
