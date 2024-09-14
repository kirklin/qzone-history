package persistence

import (
	"context"
	"gorm.io/gorm"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/pkg/database"
)

type momentRepository struct {
	db database.Database
}

func NewMomentRepository(db database.Database) repository.MomentRepository {
	return &momentRepository{db: db}
}

func (r *momentRepository) BatchImport(ctx context.Context, moments []entity.Moment) error {
	return r.db.DB().WithContext(ctx).Create(&moments).Error
}

func (r *momentRepository) Insert(ctx context.Context, moment entity.Moment) error {
	return r.db.DB().WithContext(ctx).Create(&moment).Error
}

func (r *momentRepository) FindByUserQQ(ctx context.Context, userQQ string, limit, offset int) ([]entity.Moment, error) {
	var moments []entity.Moment
	err := r.db.DB().WithContext(ctx).Where("user_qq = ?", userQQ).Limit(limit).Offset(offset).Find(&moments).Error
	return moments, err
}

func (r *momentRepository) AddLike(ctx context.Context, momentID string) error {
	return r.db.DB().WithContext(ctx).Model(&entity.Moment{}).Where("id = ?", momentID).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}

func (r *momentRepository) AddComment(ctx context.Context, comment entity.Comment) error {
	return r.db.DB().WithContext(ctx).Create(&comment).Error
}

func (r *momentRepository) IncrementViews(ctx context.Context, momentID string) error {
	return r.db.DB().WithContext(ctx).Model(&entity.Moment{}).Where("id = ?", momentID).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func (r *momentRepository) MarkAsDeleted(ctx context.Context, momentID string) error {
	return r.db.DB().WithContext(ctx).Model(&entity.Moment{}).Where("id = ?", momentID).Update("is_deleted", true).Error
}

func (r *momentRepository) MarkAsReconstructed(ctx context.Context, momentID string) error {
	return r.db.DB().WithContext(ctx).Model(&entity.Moment{}).Where("id = ?", momentID).Update("is_reconstructed", true).Error
}

func (r *momentRepository) FindByID(ctx context.Context, momentID string) (*entity.Moment, error) {
	var moment entity.Moment
	err := r.db.DB().WithContext(ctx).Where("id = ?", momentID).First(&moment).Error
	if err != nil {
		return nil, err
	}
	return &moment, nil
}
