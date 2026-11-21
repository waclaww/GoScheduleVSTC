package postgres

import "gateway/internal/domain"

type ArchiveModel struct {
	ID       int64  `gorm:"primaryKey"`
	Schedule string `gorm:"type:text;not null"`
	Date     string `gorm:"type:varchar(50);not null"`
}

func (ArchiveModel) TableName() string {
	return "vstechnical_archives"
}

func (m ArchiveModel) toDomain() *domain.Archive {
	return &domain.Archive{Id: m.ID, Date: m.Date, Schedule: m.Schedule}
}

func ArchiveModelFromDomain(a *domain.Archive) ArchiveModel {
	return ArchiveModel{ID: a.Id, Date: a.Date, Schedule: a.Schedule}
}
