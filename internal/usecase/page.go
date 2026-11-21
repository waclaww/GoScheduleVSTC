package usecase

import (
	"context"
	"errors"
	"fmt"
	"gateway/internal/domain"
	"slices"
	"time"
)

var pageNames = []string{"today", "tomorrow"}

type Page struct {
	repo PageRepository
	days DayRepository
}

func NewPage(repo PageRepository, days DayRepository) *Page {
	return &Page{repo: repo, days: days}
}

func (u *Page) Set(ctx context.Context, name string, date time.Time) (domain.PageResult, error) {
	if !slices.Contains(pageNames, name) {
		return domain.PageResult{}, fmt.Errorf("%w: unknown page %q", domain.ErrInvalidPage, name)
	}

	return u.repo.Set(ctx, name, date)
}

func (u *Page) Get(ctx context.Context, name string) (*domain.Page, error) {
	if !slices.Contains(pageNames, name) {
		return nil, domain.ErrPageNotFound
	}

	return u.repo.Get(ctx, name)
}

func (u *Page) GroupSchedule(ctx context.Context, group string) (*domain.GroupPages, error) {
	ok, err := u.days.HasGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrGroupNotFound
	}

	res := &domain.GroupPages{Group: group}
	if res.Today, err = u.groupPage(ctx, "today", group); err != nil {
		return nil, err
	}
	if res.Tomorrow, err = u.groupPage(ctx, "tomorrow", group); err != nil {
		return nil, err
	}

	return res, nil
}

func (u *Page) groupPage(ctx context.Context, name, group string) (*domain.GroupDay, error) {
	page, err := u.repo.Get(ctx, name)
	if errors.Is(err, domain.ErrPageNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return u.days.GroupDay(ctx, group, page.Date)
}
