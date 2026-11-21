package usecase

import (
	"context"
	"gateway/internal/domain"
)

type Archive struct {
	repo ArchiveRepository
}

func NewArchive(repo ArchiveRepository) *Archive {
	return &Archive{repo: repo}
}

func (u *Archive) List(ctx context.Context) ([]*domain.Archive, error) {
	return u.repo.List(ctx)
}

func (u *Archive) ListByDate(ctx context.Context, date string) ([]*domain.Archive, error) {
	return u.repo.ListByDate(ctx, date)
}

func (u *Archive) Create(ctx context.Context, a *domain.Archive) error {
	return u.repo.Create(ctx, a)
}
