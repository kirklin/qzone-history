package entity

import "time"

// Comment 表示说说的评论
type Comment struct {
	ID        string
	MomentID  string
	UserQQ    string
	Content   string
	Timestamp time.Time
	TimeText  string
}
