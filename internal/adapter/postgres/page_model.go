package postgres

import (
	"gateway/internal/domain"
	"time"
)

type PageModel struct {
	Name      string    `gorm:"type:varchar(16);primaryKey"`
	DayID     int64     `gorm:"not null"`
	Day       DayModel  `gorm:"constraint:OnDelete:RESTRICT"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"`
}

func (PageModel) TableName() string {
	return "pages"
}

func (m PageModel) toDomain() *domain.Page {
	return &domain.Page{Name: m.Name, Date: m.Day.Date, Title: m.Day.Title, UpdatedAt: m.UpdatedAt}
}
