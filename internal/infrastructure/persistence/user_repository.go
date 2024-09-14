package persistence

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/pkg/database"
)

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user entity.User) error {
	return r.db.DB().WithContext(ctx).Create(&user).Error
}

func (r *userRepository) FindByQQ(ctx context.Context, qq string) (*entity.User, error) {
	var user entity.User
	err := r.db.DB().WithContext(ctx).Where("qq = ?", qq).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user entity.User) error {
	return r.db.DB().WithContext(ctx).Save(&user).Error
}

func (r *userRepository) UpdateLoginStatus(ctx context.Context, qq string, status entity.LoginStatus) error {
	return r.db.DB().WithContext(ctx).Model(&entity.User{}).Where("qq = ?", qq).Update("login_status", status).Error
}

func (r *userRepository) GetLastLoginUser(ctx context.Context) (*entity.User, error) {
	var user entity.User
	err := r.db.DB().WithContext(ctx).
		Where("login_status = ?", entity.LoginStatusSuccess).
		Order("login_expire_time DESC").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
