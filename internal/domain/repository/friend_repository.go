package repository

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// FriendRepository 定义了好友关系存储的接口
type FriendRepository interface {
	BatchImport(ctx context.Context, friends []entity.Friend) error
	Insert(ctx context.Context, friend entity.Friend) error
	FindFriendsByUserQQ(ctx context.Context, userQQ string) ([]entity.Friend, error)
	FindFriend(ctx context.Context, userQQ, friendQQ string) (*entity.Friend, error)
	IsFriend(ctx context.Context, userQQ, friendQQ string) (bool, error)
}
