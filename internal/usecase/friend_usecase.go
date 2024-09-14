package usecase

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type friendUseCase struct {
	friendRepo repository.FriendRepository
}

func NewFriendUseCase(repo repository.FriendRepository) usecase.FriendUseCase {
	return &friendUseCase{
		friendRepo: repo,
	}
}

func (u *friendUseCase) AddFriend(ctx context.Context, friend *entity.Friend) error {
	return u.friendRepo.Insert(ctx, *friend)
}

func (u *friendUseCase) GetFriendsByUserQQ(ctx context.Context, userQQ string) ([]entity.Friend, error) {
	return u.friendRepo.FindFriendsByUserQQ(ctx, userQQ)
}

func (u *friendUseCase) GetFriend(ctx context.Context, userQQ, friendQQ string) (*entity.Friend, error) {
	return u.friendRepo.FindFriend(ctx, userQQ, friendQQ)
}

func (u *friendUseCase) CheckFriendship(ctx context.Context, userQQ, friendQQ string) (bool, error) {
	return u.friendRepo.IsFriend(ctx, userQQ, friendQQ)
}
