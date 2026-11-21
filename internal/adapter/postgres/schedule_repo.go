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

func (r *ScheduleRepo) List(ctx context.Context) ([]*domain.Schedule, error) {
	var m []ScheduleModel
	if err := r.db.WithContext(ctx).Order("id").Find(&m).Error; err != nil {
		return nil, err
	}

	schedule := make([]*domain.Schedule, 0, len(m))
	for _, s := range m {
		schedule = append(schedule, s.toDomain())
	}

	return schedule, nil
}

func (r *ScheduleRepo) Create(ctx context.Context, s *domain.Schedule) error {
	m := ScheduleModelFromDomain(s)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	s.Id = m.ID

	return nil
}

func (r *ScheduleRepo) UpdateByDate(ctx context.Context, s *domain.Schedule) error {
	return r.db.WithContext(ctx).
		Model(&ScheduleModel{}).
		Where("date = ?", s.Date).
		Update("schedule", s.Schedule).Error
}
