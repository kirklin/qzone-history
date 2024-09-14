package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type momentUseCase struct {
	momentRepo repository.MomentRepository
}

func NewMomentUseCase(repo repository.MomentRepository) usecase.MomentUseCase {
	return &momentUseCase{
		momentRepo: repo,
	}
}

func (u *momentUseCase) CreateMoment(ctx context.Context, moment *entity.Moment) error {
	return u.momentRepo.Insert(ctx, *moment)
}

func (u *momentUseCase) GetMomentsByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.Moment, error) {
	return u.momentRepo.FindByUserQQ(ctx, userQQ, limit, offset)
}

func (u *momentUseCase) AddLikeToMoment(ctx context.Context, momentID string) error {
	return u.momentRepo.AddLike(ctx, momentID)
}

func (u *momentUseCase) AddCommentToMoment(ctx context.Context, comment *entity.Comment) error {
	return u.momentRepo.AddComment(ctx, *comment)
}

func (u *momentUseCase) IncrementMomentViews(ctx context.Context, momentID string) error {
	return u.momentRepo.IncrementViews(ctx, momentID)
}

func (u *momentUseCase) MarkMomentAsDeleted(ctx context.Context, momentID string) error {
	return u.momentRepo.MarkAsDeleted(ctx, momentID)
}

func (u *momentUseCase) MarkMomentAsReconstructed(ctx context.Context, momentID string) error {
	return u.momentRepo.MarkAsReconstructed(ctx, momentID)
}

func (u *momentUseCase) GetMomentByID(ctx context.Context, momentID string) (*entity.Moment, error) {
	return u.momentRepo.FindByID(ctx, momentID)
}
