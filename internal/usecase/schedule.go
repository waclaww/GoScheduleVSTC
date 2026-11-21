package usecase

import (
	"context"
	"encoding/json"
	"gateway/internal/domain"
	"gateway/pkg/pyrepr"
	"strings"
)

type Schedule struct {
	repo ScheduleRepository
}

func NewSchedule(repo ScheduleRepository) *Schedule {
	return &Schedule{repo: repo}
}

func (u *Schedule) List(ctx context.Context) ([]*domain.Schedule, error) {
	return u.repo.List(ctx)
}

func (u *Schedule) Create(ctx context.Context, s *domain.Schedule) error {
	return u.repo.Create(ctx, s)
}

func (u *Schedule) UpdateByDate(ctx context.Context, s *domain.Schedule) error {
	return u.repo.UpdateByDate(ctx, s)
}

func (u *Schedule) ByGroup(ctx context.Context, group string) (*domain.GroupSchedule, error) {
	schedules, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(schedules) < 2 {
		return nil, domain.ErrScheduleNotEnough
	}

	today, err := groupEntry(schedules[0].Schedule, group)
	if err != nil {
		return nil, err
	}

	tomorrow, err := groupEntry(schedules[1].Schedule, group)
	if err != nil {
		return nil, err
	}

	return &domain.GroupSchedule{Today: today, Tomorrow: tomorrow}, nil
}

func groupEntry(raw, group string) (string, error) {
	var page map[string]any
	if err := json.Unmarshal([]byte(raw), &page); err != nil {
		return "", domain.ErrScheduleMalformed
	}

	value, ok := page[group]
	if !ok {
		return "", domain.ErrGroupNotFound
	}

	lessons, ok := value.(string)
	if !ok {
		return "", domain.ErrScheduleMalformed
	}

	date, ok := page["date"].(string)
	if !ok {
		return "", domain.ErrScheduleMalformed
	}

	lessons = strings.TrimSpace(lessons)
	if !strings.HasPrefix(lessons, "[") || !strings.HasSuffix(lessons, "]") {
		return "", domain.ErrScheduleMalformed
	}

	body := strings.TrimSpace(lessons[1 : len(lessons)-1])
	if body == "" {
		return "[" + pyrepr.String(date) + "]", nil
	}

	return "[" + body + ", " + pyrepr.String(date) + "]", nil
}
