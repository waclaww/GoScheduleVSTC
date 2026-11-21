package usecase

import (
	"context"
	"fmt"
	"gateway/internal/domain"
	"slices"
	"time"
)

var pageNames = []string{"today", "tomorrow"}

type Page struct {
	repo PageRepository
}

func NewPage(repo PageRepository) *Page {
	return &Page{repo: repo}
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
