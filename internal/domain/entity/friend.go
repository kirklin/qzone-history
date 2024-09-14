package entity

import "time"

// Friend 表示好友关系
type Friend struct {
	UserQQ    string
	FriendQQ  string
	Name      string
	AvatarURL string
	AddedTime time.Time
}
