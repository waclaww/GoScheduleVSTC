package postgres

import (
	"gateway/internal/domain"
	"time"
)

type GroupModel struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(32);not null;uniqueIndex"`
}

func (GroupModel) TableName() string {
	return "groups"
}

type DayModel struct {
	ID         int64           `gorm:"primaryKey"`
	Date       time.Time       `gorm:"type:date;not null;uniqueIndex"`
	Title      string          `gorm:"type:text;not null"`
	SourceHash string          `gorm:"type:char(64);not null"`
	FetchedAt  time.Time       `gorm:"type:timestamptz;not null"`
	Slots      []BellSlotModel `gorm:"foreignKey:DayID;constraint:OnDelete:CASCADE"`
	Lessons    []LessonModel   `gorm:"foreignKey:DayID;constraint:OnDelete:CASCADE"`
}

func (DayModel) TableName() string {
	return "days"
}

type BellSlotModel struct {
	ID       int64   `gorm:"primaryKey"`
	DayID    int64   `gorm:"not null;uniqueIndex:ux_bell_slots_day_slot"`
	Slot     int16   `gorm:"not null;uniqueIndex:ux_bell_slots_day_slot"`
	StartsAt *string `gorm:"type:time without time zone"`
	EndsAt   *string `gorm:"type:time without time zone"`
}

func (BellSlotModel) TableName() string {
	return "bell_slots"
}

type LessonModel struct {
	ID      int64      `gorm:"primaryKey"`
	DayID   int64      `gorm:"not null;uniqueIndex:ux_lessons_day_group_slot"`
	GroupID int64      `gorm:"not null;uniqueIndex:ux_lessons_day_group_slot;index"`
	Group   GroupModel `gorm:"constraint:OnDelete:RESTRICT"`
	Slot    int16      `gorm:"not null;uniqueIndex:ux_lessons_day_group_slot"`
	Subject string     `gorm:"type:text;not null"`
	Room    string     `gorm:"type:varchar(64);not null"`
}

func (LessonModel) TableName() string {
	return "lessons"
}

type groupLessonRow struct {
	Slot     int16
	Subject  string
	Room     string
	StartsAt *string
	EndsAt   *string
}

func (m DayModel) toDomain() *domain.Day {
	d := &domain.Day{
		Id:         m.ID,
		Date:       m.Date,
		Title:      m.Title,
		SourceHash: m.SourceHash,
		FetchedAt:  m.FetchedAt,
		Slots:      make([]domain.BellSlot, 0, len(m.Slots)),
		Lessons:    make([]domain.Lesson, 0, len(m.Lessons)),
	}

	for _, s := range m.Slots {
		d.Slots = append(d.Slots, domain.BellSlot{Slot: int(s.Slot), StartsAt: clock(s.StartsAt), EndsAt: clock(s.EndsAt)})
	}

	groups := make(map[string]bool)
	for _, l := range m.Lessons {
		d.Lessons = append(d.Lessons, domain.Lesson{Group: l.Group.Name, Slot: int(l.Slot), Subject: l.Subject, Room: l.Room})
		if !groups[l.Group.Name] {
			groups[l.Group.Name] = true
			d.Groups = append(d.Groups, l.Group.Name)
		}
	}

	return d
}

func (r groupLessonRow) toDomain() domain.GroupLesson {
	return domain.GroupLesson{
		Slot:     int(r.Slot),
		StartsAt: clock(r.StartsAt),
		EndsAt:   clock(r.EndsAt),
		Subject:  r.Subject,
		Room:     r.Room,
	}
}

func DayModelFromDomain(d *domain.Day) DayModel {
	return DayModel{
		ID:         d.Id,
		Date:       d.Date,
		Title:      d.Title,
		SourceHash: d.SourceHash,
		FetchedAt:  d.FetchedAt,
	}
}

func BellSlotModelsFromDomain(dayID int64, slots []domain.BellSlot) []BellSlotModel {
	res := make([]BellSlotModel, 0, len(slots))
	for _, s := range slots {
		res = append(res, BellSlotModel{DayID: dayID, Slot: int16(s.Slot), StartsAt: nullableClock(s.StartsAt), EndsAt: nullableClock(s.EndsAt)})
	}

	return res
}

func LessonModelsFromDomain(dayID int64, groupIDs map[string]int64, lessons []domain.Lesson) []LessonModel {
	res := make([]LessonModel, 0, len(lessons))
	for _, l := range lessons {
		res = append(res, LessonModel{DayID: dayID, GroupID: groupIDs[l.Group], Slot: int16(l.Slot), Subject: l.Subject, Room: l.Room})
	}

	return res
}

func clock(v *string) string {
	if v == nil || len(*v) < 5 {
		return ""
	}

	return (*v)[:5]
}

func nullableClock(v string) *string {
	if v == "" {
		return nil
	}

	return &v
}
