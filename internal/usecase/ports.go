package usecase

import (
	"context"
	"gateway/internal/domain"
)

type ScheduleRepository interface {
	List(ctx context.Context) ([]*domain.Schedule, error)
	Create(ctx context.Context, s *domain.Schedule) error
	UpdateByDate(ctx context.Context, s *domain.Schedule) error
}

type ArchiveRepository interface {
	List(ctx context.Context) ([]*domain.Archive, error)
	ListByDate(ctx context.Context, date string) ([]*domain.Archive, error)
	Create(ctx context.Context, a *domain.Archive) error
}
