package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type reconstructionUseCase struct {
	activityRepo     repository.ActivityRepository
	momentRepo       repository.MomentRepository
	boardMessageRepo repository.BoardMessageRepository
}

func NewReconstructionUseCase(
	activityRepo repository.ActivityRepository,
	momentRepo repository.MomentRepository,
	boardMessageRepo repository.BoardMessageRepository,
) usecase.ReconstructionUseCase {
	return &reconstructionUseCase{
		activityRepo:     activityRepo,
		momentRepo:       momentRepo,
		boardMessageRepo: boardMessageRepo,
	}
}

func (u *reconstructionUseCase) ReconstructMomentsFromActivities(ctx context.Context, userQQ string) error {
	activities, err := u.activityRepo.FindByUserQQ(ctx, userQQ, 1000, 0)
	if err != nil {
		return err
	}

	for _, activity := range activities {
		//TODO Check
		if activity.Type == entity.TypeMoment || activity.Type == entity.TypeLike || activity.Type == entity.TypeComment {
			// 根据活动记录重建Moment
			moment := reconstructMomentFromActivity(activity)
			if err := u.momentRepo.Insert(ctx, moment); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *reconstructionUseCase) ReconstructBoardMessagesFromActivities(ctx context.Context, userQQ string) error {
	activities, err := u.activityRepo.FindByUserQQ(ctx, userQQ, 1000, 0)
	if err != nil {
		return err
	}

	for _, activity := range activities {
		if activity.Type == entity.TypeBoardMessage {
			// 根据活动记录重建BoardMessage
			boardMessage := reconstructBoardMessageFromActivity(activity)
			if err := u.boardMessageRepo.Insert(ctx, boardMessage); err != nil {
				return err
			}
		}
	}

	return nil
}

// 辅助函数，根据活动记录重建Moment
func reconstructMomentFromActivity(activity entity.Activity) entity.Moment {
	// TODO 实现重建逻辑
	return entity.Moment{
		UserQQ:    activity.SenderQQ,
		Content:   activity.Content,
		Timestamp: activity.Timestamp,
		TimeText:  activity.TimeText,
		ImageURLs: activity.ImageURLs,
	}
}

// 辅助函数，根据活动记录重建BoardMessage
func reconstructBoardMessageFromActivity(activity entity.Activity) entity.BoardMessage {
	// TODO 实现重建逻辑
	return entity.BoardMessage{
		UserQQ:    activity.ReceiverQQ,
		SenderQQ:  activity.SenderQQ,
		Content:   activity.Content,
		Timestamp: activity.Timestamp,
		TimeText:  activity.TimeText,
	}
}
