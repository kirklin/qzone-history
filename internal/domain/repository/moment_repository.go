package repository

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// MomentRepository 定义了说说存储的接口
type MomentRepository interface {
	BatchImport(ctx context.Context, moments []entity.Moment) error
	Insert(ctx context.Context, moment entity.Moment) error
	FindByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.Moment, error)
	AddLike(ctx context.Context, momentID string) error
	AddComment(ctx context.Context, comment entity.Comment) error
	IncrementViews(ctx context.Context, momentID string) error
	MarkAsDeleted(ctx context.Context, momentID string) error
	MarkAsReconstructed(ctx context.Context, momentID string) error
	FindByID(ctx context.Context, momentID string) (*entity.Moment, error)
}
