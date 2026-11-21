package postgres

import "gateway/internal/domain"

type ScheduleModel struct {
	ID       int64  `gorm:"primaryKey"`
	Schedule string `gorm:"type:text;not null"`
	Date     string `gorm:"type:varchar(50);not null"`
}

func (ScheduleModel) TableName() string {
	return "vstechnical_datas"
}

func (m ScheduleModel) toDomain() *domain.Schedule {
	return &domain.Schedule{Id: m.ID, Date: m.Date, Schedule: m.Schedule}
}

func ScheduleModelFromDomain(s *domain.Schedule) ScheduleModel {
	return ScheduleModel{ID: s.Id, Date: s.Date, Schedule: s.Schedule}
}
