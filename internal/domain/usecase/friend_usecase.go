package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// FriendUseCase 定义了好友相关的用例接口
type FriendUseCase interface {
	// AddFriend 添加好友
	AddFriend(ctx context.Context, friend *entity.Friend) error

	// GetFriendsByUserQQ 获取用户的好友列表
	GetFriendsByUserQQ(ctx context.Context, userQQ string) ([]entity.Friend, error)

	// GetFriend 获取特定好友信息
	GetFriend(ctx context.Context, userQQ, friendQQ string) (*entity.Friend, error)

	// CheckFriendship 检查两个用户是否为好友关系
	CheckFriendship(ctx context.Context, userQQ, friendQQ string) (bool, error)
}
