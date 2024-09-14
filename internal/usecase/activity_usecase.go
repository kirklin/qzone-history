package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type activityUseCase struct {
	activityRepo repository.ActivityRepository
}

func NewActivityUseCase(activityRepo repository.ActivityRepository) usecase.ActivityUseCase {
	return &activityUseCase{
		activityRepo: activityRepo,
	}
}

func (a *activityUseCase) GetActivities(ctx context.Context, userQQ string, limit, offset int) ([]entity.Activity, error) {
	return a.activityRepo.FindByUserQQ(ctx, userQQ, limit, offset)
}

func (a *activityUseCase) GetAllActivities(ctx context.Context, userQQ string) ([]entity.Activity, error) {
	return a.activityRepo.FindByUserQQ(ctx, userQQ, 0, 0) // 0, 0 表示获取所有活动
}

func (a *activityUseCase) SaveActivity(ctx context.Context, activity entity.Activity) error {
	return a.activityRepo.Insert(ctx, activity)
}

func (a *activityUseCase) GetActivityCount(ctx context.Context, userQQ string) (int, error) {
	activities, err := a.activityRepo.FindByUserQQ(ctx, userQQ, 0, 0)
	if err != nil {
		return 0, err
	}
	return len(activities), nil
}

func (a *activityUseCase) GetActivitiesByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]entity.Activity, error) {
	return a.activityRepo.FindByType(ctx, activityType, limit, offset)
}

func (a *activityUseCase) FetchActivities(ctx context.Context, user entity.User) ([]entity.Activity, error) {
	//TODO implement me
	panic("implement me")
}
