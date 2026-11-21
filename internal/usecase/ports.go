package usecase

import (
	"context"
	"gateway/internal/domain"
	"time"
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

type DayRepository interface {
	Upsert(ctx context.Context, d *domain.Day) (domain.UpsertResult, error)
	Get(ctx context.Context, date time.Time) (*domain.Day, error)
	Groups(ctx context.Context) ([]string, error)
	GroupDay(ctx context.Context, group string, date time.Time) (*domain.GroupDay, error)
}

type PageRepository interface {
	Set(ctx context.Context, name string, date time.Time) (domain.PageResult, error)
	Get(ctx context.Context, name string) (*domain.Page, error)
}
