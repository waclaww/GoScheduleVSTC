package postgres

import (
	"gateway/internal/domain"

	"gorm.io/datatypes"
)

type ScheduleModel struct {
	ID      int64          `gorm:"primaryKey"`
	Type    string         `gorm:"type:string"`
	Date    string         `gorm:"type:date"`
	Content datatypes.JSON `gorm:"type:jsonb"`
}

func (ScheduleModel) TableName() string {
	return "schedule"
}

func (m ScheduleModel) toDomain() *domain.Get {
	return &domain.Get{Id: m.ID, Type: m.Type, Date: m.Date, Content: m.Content.GormDataType()}
}

func ScheduleModelFromDomain(s *domain.Get) ScheduleModel {
	return ScheduleModel{ID: s.Id, Type: s.Type, Date: s.Date, Content: datatypes.JSON(s.Content)}
}
