package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// AuthUseCase 定义了身份认证相关的用例接口
type AuthUseCase interface {
	// GetLoginQRCode 获取登录二维码
	GetLoginQRCode(ctx context.Context) ([]byte, string, error)

	// CheckLoginStatus 检查登录状态
	CheckLoginStatus(ctx context.Context, qrsig string) (entity.LoginStatus, error)

	// CompleteLogin 完成登录过程
	CompleteLogin(ctx context.Context, responseText string) (*entity.User, error)
}
