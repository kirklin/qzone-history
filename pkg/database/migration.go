package database

import (
	"qzone-history/internal/domain/entity"
)

func AutoMigrate(db Database) error {
	return db.DB().AutoMigrate(
		&entity.User{},         // 用户实体
		&entity.Moment{},       // 说说实体
		&entity.BoardMessage{}, // 留言板消息实体
		&entity.Activity{},     // 活动记录实体
		&entity.Comment{},      // 评论实体
		&entity.Friend{},       // 好友实体
	)
}
