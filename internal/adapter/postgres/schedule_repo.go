package postgres

import (
	"context"
	"gateway/internal/domain"

	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) Get(ctx context.Context) (*domain.Get, error) {
	var m ScheduleModel
	err := r.db.WithContext(ctx).Find(&m).Error
	if err != nil {
		return nil, err
	}

	return m.toDomain(), nil
}
