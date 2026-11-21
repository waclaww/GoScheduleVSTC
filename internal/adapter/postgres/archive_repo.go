package postgres

import (
	"context"
	"gateway/internal/domain"

	"gorm.io/gorm"
)

type ArchiveRepo struct {
	db *gorm.DB
}

func NewArchiveRepo(db *gorm.DB) *ArchiveRepo {
	return &ArchiveRepo{db: db}
}

func (r *ArchiveRepo) List(ctx context.Context) ([]*domain.Archive, error) {
	return r.find(r.db.WithContext(ctx))
}

func (r *ArchiveRepo) ListByDate(ctx context.Context, date string) ([]*domain.Archive, error) {
	return r.find(r.db.WithContext(ctx).Where("date = ?", date))
}

func (r *ArchiveRepo) Create(ctx context.Context, a *domain.Archive) error {
	m := ArchiveModelFromDomain(a)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	a.Id = m.ID

	return nil
}

func (r *ArchiveRepo) find(q *gorm.DB) ([]*domain.Archive, error) {
	var m []ArchiveModel
	if err := q.Order("id").Find(&m).Error; err != nil {
		return nil, err
	}

	archive := make([]*domain.Archive, 0, len(m))
	for _, a := range m {
		archive = append(archive, a.toDomain())
	}

	return archive, nil
}
