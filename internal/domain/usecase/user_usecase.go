package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// UserUseCase 定义了用户相关的用例接口
type UserUseCase interface {
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, userQQ string) (*entity.User, error)

	// UpdateUserInfo 更新用户信息
	UpdateUserInfo(ctx context.Context, user *entity.User) error

	// SaveUser 保存用户信息
	SaveUser(ctx context.Context, user *entity.User) error
}
