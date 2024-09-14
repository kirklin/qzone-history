package entity

import "time"

// Activity 表示用户活动记录
type Activity struct {
	ID         string       `json:"id" gorm:"primaryKey"`       // 活动ID
	SenderQQ   string       `json:"senderQQ" gorm:"index"`      // 发送者QQ
	SenderName string       `json:"senderName"`                 // 发送者名称
	SenderLink string       `json:"senderLink"`                 // 发送者链接
	ReceiverQQ string       `json:"receiverQQ" gorm:"index"`    // 接收者QQ
	Content    string       `json:"content"`                    // 活动内容
	Timestamp  time.Time    `json:"timestamp"`                  // 时间戳
	TimeText   string       `json:"timeText"`                   // 时间文本
	ImageURLs  []string     `json:"imageURLs" gorm:"type:text"` // 图片URL列表
	Type       ActivityType `json:"type"`                       // 活动类型
}

// ActivityType 表示活动类型
type ActivityType int

const (
	TypeMoment       ActivityType = iota // 说说
	TypeForward                          // 转发
	TypeLike                             // 点赞
	TypeComment                          // 评论
	TypeBoardMessage                     // 留言
	TypeOther                            // 其他
)
