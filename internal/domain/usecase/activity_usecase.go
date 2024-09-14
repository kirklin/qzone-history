package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// ActivityUseCase 定义了用户活动相关的用例接口
type ActivityUseCase interface {
	// GetActivities 获取指定数量的用户活动
	GetActivities(ctx context.Context, userQQ string, offset, count int) ([]entity.Activity, error)

	// GetAllActivities 获取用户的所有活动
	GetAllActivities(ctx context.Context, userQQ string) ([]entity.Activity, error)

	// SaveActivity 保存单个活动
	SaveActivity(ctx context.Context, activity entity.Activity) error

	// GetActivityCount 获取用户活动总数
	GetActivityCount(ctx context.Context, userQQ string) (int, error)

	// GetActivitiesByType 根据活动类型获取活动
	GetActivitiesByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]entity.Activity, error)

	// FetchActivities 从网页或其他接口获取全部活动
	FetchActivities(ctx context.Context, user entity.User) ([]entity.Activity, error)

	// FetchActivity 从网页或其他接口获取一条活动
	FetchActivity(ctx context.Context, user entity.User, offset int) (entity.Activity, error)
}
