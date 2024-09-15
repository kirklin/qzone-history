package persistence

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/pkg/database"
)

type friendRepository struct {
	db database.Database
}

func NewFriendRepository(db database.Database) repository.FriendRepository {
	return &friendRepository{db: db}
}

func (r *friendRepository) BatchImport(ctx context.Context, friends []entity.Friend) error {
	return r.db.DB().WithContext(ctx).Save(&friends).Error
}

func (r *friendRepository) Insert(ctx context.Context, friend entity.Friend) error {
	return r.db.DB().WithContext(ctx).Save(&friend).Error
}

func (r *friendRepository) FindFriendsByUserQQ(ctx context.Context, userQQ string) ([]entity.Friend, error) {
	var friends []entity.Friend
	err := r.db.DB().WithContext(ctx).Where("user_qq = ?", userQQ).Find(&friends).Error
	return friends, err
}

func (r *friendRepository) FindFriend(ctx context.Context, userQQ, friendQQ string) (*entity.Friend, error) {
	var friend entity.Friend
	err := r.db.DB().WithContext(ctx).Where("user_qq = ? AND friend_qq = ?", userQQ, friendQQ).First(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

func (r *friendRepository) IsFriend(ctx context.Context, userQQ, friendQQ string) (bool, error) {
	var count int64
	err := r.db.DB().WithContext(ctx).Model(&entity.Friend{}).Where("user_qq = ? AND friend_qq = ?", userQQ, friendQQ).Count(&count).Error
	return count > 0, err
}
