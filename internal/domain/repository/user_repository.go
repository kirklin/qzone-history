package repository

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// UserRepository 定义了用户存储的接口，支持保存、通过 QQ 获取用户，以及获取所有用户
// UserRepository defines the interface for user storage, supporting saving, getting by QQ, and retrieving all users.
type UserRepository interface {
	// Save 保存单个用户信息
	Save(ctx context.Context, user entity.User) error
	// FindByQQ 通过QQ号查找用户
	FindByQQ(ctx context.Context, qq string) (*entity.User, error)
	// GetLastLoginUser 获取最后登录的用户
	GetLastLoginUser(ctx context.Context) (*entity.User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user entity.User) error
	// UpdateLoginStatus 更新用户的登录状态
	UpdateLoginStatus(ctx context.Context, qq string, status entity.LoginStatus) error
}
