package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// MomentUseCase 定义了说说相关的用例接口
type MomentUseCase interface {
	// CreateMoment 创建新的说说
	CreateMoment(ctx context.Context, moment *entity.Moment) error

	// GetMomentsByUserQQ 获取用户的说说列表
	GetMomentsByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.Moment, error)

	// AddLikeToMoment 为说说添加点赞
	AddLikeToMoment(ctx context.Context, momentID string) error

	// AddCommentToMoment 为说说添加评论
	AddCommentToMoment(ctx context.Context, comment *entity.Comment) error

	// IncrementMomentViews 增加说说的浏览次数
	IncrementMomentViews(ctx context.Context, momentID string) error

	// MarkMomentAsDeleted 标记说说为已删除
	MarkMomentAsDeleted(ctx context.Context, momentID string) error

	// MarkMomentAsReconstructed 标记说说为已重建
	MarkMomentAsReconstructed(ctx context.Context, momentID string) error

	// GetMomentByID 根据ID获取说说
	GetMomentByID(ctx context.Context, momentID string) (*entity.Moment, error)
}
