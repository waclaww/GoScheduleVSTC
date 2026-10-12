package usecase

import (
	"context"
	"gateway/internal/domain"
)

type ScheduleRepository interface {
	Get(ctx context.Context, u *domain.Get)
}
