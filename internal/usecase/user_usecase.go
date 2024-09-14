package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) usecase.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (u *userUseCase) GetUserInfo(ctx context.Context, userQQ string) (*entity.User, error) {
	return u.userRepo.FindByQQ(ctx, userQQ)
}

func (u *userUseCase) UpdateUserInfo(ctx context.Context, user *entity.User) error {
	return u.userRepo.Update(ctx, *user)
}

func (u *userUseCase) SaveUser(ctx context.Context, user *entity.User) error {
	return u.userRepo.Save(ctx, *user)
}
