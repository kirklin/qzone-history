package repository

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// BoardMessageRepository 定义了留言板消息存储的接口
type BoardMessageRepository interface {
	BatchImport(ctx context.Context, messages []entity.BoardMessage) error
	Insert(ctx context.Context, message entity.BoardMessage) error
	FindByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.BoardMessage, error)
	FindByID(ctx context.Context, messageID string) (*entity.BoardMessage, error)
}
