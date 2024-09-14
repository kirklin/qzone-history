package persistence

import (
	"context"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/pkg/database"
)

type boardMessageRepository struct {
	db database.Database
}

func NewBoardMessageRepository(db database.Database) repository.BoardMessageRepository {
	return &boardMessageRepository{db: db}
}

func (r *boardMessageRepository) BatchImport(ctx context.Context, messages []entity.BoardMessage) error {
	return r.db.DB().WithContext(ctx).Create(&messages).Error
}

func (r *boardMessageRepository) Insert(ctx context.Context, message entity.BoardMessage) error {
	return r.db.DB().WithContext(ctx).Create(&message).Error
}

func (r *boardMessageRepository) FindByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.BoardMessage, error) {
	var messages []entity.BoardMessage
	err := r.db.DB().WithContext(ctx).Where("user_qq = ?", userQQ).Limit(limit).Offset(offset).Find(&messages).Error
	return messages, err
}

func (r *boardMessageRepository) FindByID(ctx context.Context, messageID string) (*entity.BoardMessage, error) {
	var message entity.BoardMessage
	err := r.db.DB().WithContext(ctx).Where("id = ?", messageID).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}
