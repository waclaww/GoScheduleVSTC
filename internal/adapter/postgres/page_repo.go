package postgres

import (
	"context"
	"errors"
	"gateway/internal/domain"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageRepo struct {
	db *gorm.DB
}

func NewPageRepo(db *gorm.DB) *PageRepo {
	return &PageRepo{db: db}
}

func (r *PageRepo) Set(ctx context.Context, name string, date time.Time) (domain.PageResult, error) {
	res := domain.PageResult{Name: name, Date: date}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var day DayModel
		err := tx.Where("date = ?", date).Take(&day).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrDayNotFound
		}
		if err != nil {
			return err
		}

		var page PageModel
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Day").Where("name = ?", name).Take(&page).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			page = PageModel{Name: name}
		case err != nil:
			return err
		default:
			previous := page.Day.Date
			res.PreviousDate = &previous
			res.Changed = !previous.Equal(day.Date)
		}

		page.DayID = day.ID
		page.Day = DayModel{}
		page.UpdatedAt = time.Now()

		return tx.Omit("Day").Save(&page).Error
	})

	return res, err
}

func (r *PageRepo) Get(ctx context.Context, name string) (*domain.Page, error) {
	var m PageModel
	err := r.db.WithContext(ctx).Preload("Day").Where("name = ?", name).Take(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPageNotFound
	}
	if err != nil {
		return nil, err
	}

	return m.toDomain(), nil
}
