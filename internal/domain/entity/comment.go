package entity

import "time"

// Comment 表示说说的评论
type Comment struct {
	ID        string    `json:"id" gorm:"primaryKey"`  // 评论ID
	MomentID  string    `json:"momentID" gorm:"index"` // 关联的说说ID
	UserQQ    string    `json:"userQQ" gorm:"index"`   // 评论者QQ
	Content   string    `json:"content"`               // 评论内容
	Timestamp time.Time `json:"timestamp"`             // 时间戳
	TimeText  string    `json:"timeText"`              // 时间文本
}
