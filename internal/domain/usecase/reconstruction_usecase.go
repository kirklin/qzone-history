package usecase

import (
	"context"
)

// ReconstructionUseCase 定义了数据重建相关的用例接口
type ReconstructionUseCase interface {
	// ReconstructMomentsFromActivities 从活动记录重建说说
	ReconstructMomentsFromActivities(ctx context.Context, userQQ string) error

	// ReconstructBoardMessagesFromActivities 从活动记录重建留言板消息
	ReconstructBoardMessagesFromActivities(ctx context.Context, userQQ string) error
}
