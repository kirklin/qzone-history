package entity

import "time"

// User 表示QQ用户
type User struct {
	QQ              string
	Nickname        string
	AvatarURL       string
	Cookies         map[string]string
	LoginStatus     LoginStatus
	LoginExpireTime time.Time
}
