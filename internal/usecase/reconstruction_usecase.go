package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
	activities, err := u.activityRepo.FindByUserQQ(ctx, userQQ, -1, 0)
	if err != nil {
		return err
	}

	momentMap := make(map[string]*entity.Moment)

	for _, activity := range activities {
		momentKey := generateMomentKey(activity)
		moment, ok := momentMap[momentKey]
		if !ok {
			moment = &entity.Moment{
				ID:              momentKey,
				UserQQ:          activity.SenderQQ,
				Content:         activity.Content,
				Timestamp:       activity.Timestamp,
				TimeText:        activity.TimeText,
				ImageURLs:       activity.ImageURLs,
				IsReconstructed: true,
			}
			momentMap[momentKey] = moment
		}

		switch activity.Type {
		case entity.TypeMoment:
			updateMomentFromActivity(moment, activity)
		case entity.TypeLike:
			moment.Likes++
		case entity.TypeView:
			moment.Views++
		case entity.TypeComment:
			comment := reconstructCommentFromActivity(activity)
			moment.Comments = append(moment.Comments, comment)
		}
	}

	// 将重建的Moment插入或更新到数据库
	for _, moment := range momentMap {
		if err := u.momentRepo.UpsertMoment(ctx, *moment); err != nil {
			return err
		}
	}

	return nil
}

func generateMomentKey(activity entity.Activity) string {
	// 使用内容和收到的人的QQ生成唯一键
	key := activity.Content + activity.ReceiverQQ
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// 更新Moment信息
func updateMomentFromActivity(moment *entity.Moment, activity entity.Activity) {
	// 更新逻辑，可以选择保留更详细的信息
	if len(activity.Content) > len(moment.Content) {
		moment.Content = activity.Content
	}
	if activity.Timestamp.Before(moment.Timestamp) {
		moment.Timestamp = activity.Timestamp
		moment.TimeText = activity.TimeText
	}
	if len(activity.ImageURLs) > len(moment.ImageURLs) {
		moment.ImageURLs = activity.ImageURLs
	}
}

// 从活动重建评论
func reconstructCommentFromActivity(activity entity.Activity) entity.Comment {
	return entity.Comment{
		UserQQ:    activity.SenderQQ,
		Content:   activity.Content,
		Timestamp: activity.Timestamp,
		TimeText:  activity.TimeText,
	}
}

func (u *reconstructionUseCase) ReconstructBoardMessagesFromActivities(ctx context.Context, userQQ string) error {
	activities, err := u.activityRepo.FindByUserQQ(ctx, userQQ, -1, 0)
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
