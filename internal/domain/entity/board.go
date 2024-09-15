package entity

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/gorm"
	"time"
)

// BoardMessage 表示留言板消息
type BoardMessage struct {
	ID        string    `json:"id" gorm:"primaryKey"`  // 留言ID
	UserQQ    string    `json:"userQQ" gorm:"index"`   // 留言板所有者QQ
	SenderQQ  string    `json:"senderQQ" gorm:"index"` // 发送者QQ
	Content   string    `json:"content"`               // 留言内容
	Timestamp time.Time `json:"timestamp"`             // 时间戳
	TimeText  string    `json:"timeText"`              // 时间文本
}

// BeforeCreate 钩子，在创建记录之前自动生成ID
func (message *BoardMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if message.ID == "" {
		// 使用内容和收到的人的QQ生成唯一键
		key := message.Content + message.UserQQ + message.SenderQQ
		hash := md5.Sum([]byte(key))
		message.ID = hex.EncodeToString(hash[:])
	}
	return
}
