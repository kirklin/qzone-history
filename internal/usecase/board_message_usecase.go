package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type boardMessageUseCase struct {
	boardMessageRepo repository.BoardMessageRepository
}

func NewBoardMessageUseCase(repo repository.BoardMessageRepository) usecase.BoardMessageUseCase {
	return &boardMessageUseCase{
		boardMessageRepo: repo,
	}
}

func (u *boardMessageUseCase) CreateBoardMessage(ctx context.Context, message *entity.BoardMessage) error {
	return u.boardMessageRepo.Insert(ctx, *message)
}

func (u *boardMessageUseCase) GetBoardMessagesByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.BoardMessage, error) {
	return u.boardMessageRepo.FindByUserQQ(ctx, userQQ, limit, offset)
}

func (u *boardMessageUseCase) GetBoardMessageByID(ctx context.Context, messageID string) (*entity.BoardMessage, error) {
	return u.boardMessageRepo.FindByID(ctx, messageID)
}
