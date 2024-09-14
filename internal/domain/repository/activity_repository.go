package repository

import (
	"context"
	"qzone-history/internal/domain/entity"
)

// ActivityRepository 定义了活动记录存储的接口
type ActivityRepository interface {
	BatchImport(ctx context.Context, activities []entity.Activity) error
	Insert(ctx context.Context, activity entity.Activity) error
	FindByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.Activity, error)
	FindByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]entity.Activity, error)
}
