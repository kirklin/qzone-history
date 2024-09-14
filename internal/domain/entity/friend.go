package entity

import "time"

// Friend 表示好友关系
type Friend struct {
	UserQQ    string    `json:"userQQ" gorm:"primaryKey"`   // 用户QQ
	FriendQQ  string    `json:"friendQQ" gorm:"primaryKey"` // 好友QQ
	Name      string    `json:"name"`                       // 好友名称
	AvatarURL string    `json:"avatarURL"`                  // 头像URL
	AddedTime time.Time `json:"addedTime"`                  // 添加好友的时间
}
