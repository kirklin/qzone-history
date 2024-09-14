package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// BoardMessageUseCase 定义了留言板消息相关的用例接口
type BoardMessageUseCase interface {
	// CreateBoardMessage 创建新的留言板消息
	CreateBoardMessage(ctx context.Context, message *entity.BoardMessage) error

	// GetBoardMessagesByUserQQ 获取用户的留言板消息列表
	GetBoardMessagesByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.BoardMessage, error)

	// GetBoardMessageByID 根据ID获取留言板消息
	GetBoardMessageByID(ctx context.Context, messageID string) (*entity.BoardMessage, error)
}
