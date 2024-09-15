package persistence

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/pkg/database"
)

type activityRepository struct {
	db database.Database
}

func NewActivityRepository(db database.Database) repository.ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) BatchImport(ctx context.Context, activities []entity.Activity) error {
	//return r.db.DB().WithContext(ctx).CreateInBatches(&activities, 100).Error
	return r.db.DB().WithContext(ctx).Save(&activities).Error
}

func (r *activityRepository) Insert(ctx context.Context, activity entity.Activity) error {
	return r.db.DB().WithContext(ctx).Save(&activity).Error
}

func (r *activityRepository) FindByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.Activity, error) {
	var activities []entity.Activity
	err := r.db.DB().WithContext(ctx).Where("receiver_qq = ?", userQQ).Limit(limit).Offset(offset).Find(&activities).Error
	return activities, err
}

func (r *activityRepository) FindByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]entity.Activity, error) {
	var activities []entity.Activity
	err := r.db.DB().WithContext(ctx).Where("type = ?", activityType).Limit(limit).Offset(offset).Find(&activities).Error
	return activities, err
}
