package entity

import "time"

// User 表示QQ用户
type User struct {
	QQ              string            `json:"qq" gorm:"primaryKey"`           // QQ号
	Nickname        string            `json:"nickname"`                       // 昵称
	AvatarURL       string            `json:"avatarURL"`                      // 头像URL
	Cookies         map[string]string `json:"cookies" gorm:"serializer:json"` // 登录Cookie
	LoginStatus     LoginStatus       `json:"loginStatus"`                    // 登录状态
	LoginExpireTime time.Time         `json:"loginExpireTime"`                // 登录过期时间
}
