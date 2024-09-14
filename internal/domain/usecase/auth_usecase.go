package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// AuthUseCase 定义了身份认证相关的用例接口
type AuthUseCase interface {
	// CheckLocalLoginStatus 检查本地存储的登录状态
	// 如果本地有有效的登录信息，返回用户信息和 true
	// 如果本地没有有效的登录信息，返回 nil 和 false
	CheckLocalLoginStatus(ctx context.Context) (*entity.User, bool, error)

	// GetLoginQRCode 获取登录二维码
	// 返回二维码图片数据和用于后续状态检查的 qrsig
	GetLoginQRCode(ctx context.Context) ([]byte, string, error)

	// CheckQRCodeLoginStatus 检查二维码登录状态
	// 返回登录状态和可能的错误
	CheckQRCodeLoginStatus(ctx context.Context, qrsig string) (entity.LoginStatus, string, error)

	// CompleteLogin 完成登录过程
	// 参数 loginResponse 是从 QQ 服务器返回的登录成功响应
	// 返回登录成功的用户信息
	CompleteLogin(ctx context.Context, loginResponse string) (*entity.User, error)

	// RefreshLogin 刷新登录状态
	// 使用已保存的用户信息尝试刷新登录状态
	// 如果刷新成功，返回更新后的用户信息；如果失败，返回错误
	RefreshLogin(ctx context.Context, user *entity.User) (*entity.User, error)

	// Logout 登出当前用户
	// 清除本地存储的登录信息
	Logout(ctx context.Context) error
}
