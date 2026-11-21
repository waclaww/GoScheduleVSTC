package postgres

import (
	"context"
	"errors"
	"gateway/internal/domain"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const lessonBatchSize = 500

type DayRepo struct {
	db *gorm.DB
}

func NewDayRepo(db *gorm.DB) *DayRepo {
	return &DayRepo{db: db}
}

func (r *DayRepo) Upsert(ctx context.Context, d *domain.Day) (domain.UpsertResult, error) {
	var res domain.UpsertResult

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing DayModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("date = ?", d.Date).Take(&existing).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			res.Created = true
		case err != nil:
			return err
		case existing.SourceHash == d.SourceHash:
			d.Id = existing.ID
			return tx.Model(&existing).Update("fetched_at", d.FetchedAt).Error
		default:
			res.Changed = true
			d.Id = existing.ID
		}

		m := DayModelFromDomain(d)
		if err := tx.Save(&m).Error; err != nil {
			return err
		}
		d.Id = m.ID

		if err := tx.Where("day_id = ?", m.ID).Delete(&BellSlotModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("day_id = ?", m.ID).Delete(&LessonModel{}).Error; err != nil {
			return err
		}

		groupIDs, err := upsertGroups(tx, d.Groups)
		if err != nil {
			return err
		}

		if slots := BellSlotModelsFromDomain(m.ID, d.Slots); len(slots) > 0 {
			if err := tx.Create(&slots).Error; err != nil {
				return err
			}
		}

		if lessons := LessonModelsFromDomain(m.ID, groupIDs, d.Lessons); len(lessons) > 0 {
			if err := tx.Omit("Group").CreateInBatches(&lessons, lessonBatchSize).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return res, err
}

func (r *DayRepo) Get(ctx context.Context, date time.Time) (*domain.Day, error) {
	var m DayModel
	err := r.db.WithContext(ctx).
		Preload("Slots", func(db *gorm.DB) *gorm.DB { return db.Order("slot") }).
		Preload("Lessons", func(db *gorm.DB) *gorm.DB { return db.Order("group_id, slot") }).
		Preload("Lessons.Group").
		Where("date = ?", date).
		Take(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDayNotFound
	}
	if err != nil {
		return nil, err
	}

	return m.toDomain(), nil
}

func (r *DayRepo) Groups(ctx context.Context) ([]string, error) {
	var names []string
	if err := r.db.WithContext(ctx).Model(&GroupModel{}).Order("name").Pluck("name", &names).Error; err != nil {
		return nil, err
	}

	return names, nil
}

func (r *DayRepo) GroupDay(ctx context.Context, group string, date time.Time) (*domain.GroupDay, error) {
	db := r.db.WithContext(ctx)

	var g GroupModel
	err := db.Where("name = ?", group).Take(&g).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}

	var d DayModel
	err = db.Where("date = ?", date).Take(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDayNotFound
	}
	if err != nil {
		return nil, err
	}

	var rows []groupLessonRow
	err = db.Table("lessons AS l").
		Select("l.slot, l.subject, l.room, b.starts_at::text AS starts_at, b.ends_at::text AS ends_at").
		Joins("LEFT JOIN bell_slots AS b ON b.day_id = l.day_id AND b.slot = l.slot").
		Where("l.day_id = ? AND l.group_id = ?", d.ID, g.ID).
		Order("l.slot").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	lessons := make([]domain.GroupLesson, 0, len(rows))
	for _, row := range rows {
		lessons = append(lessons, row.toDomain())
	}

	return &domain.GroupDay{Date: d.Date, Title: d.Title, Group: g.Name, Lessons: lessons}, nil
}

func (r *DayRepo) HasGroup(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func upsertGroups(tx *gorm.DB, names []string) (map[string]int64, error) {
	ids := make(map[string]int64, len(names))
	if len(names) == 0 {
		return ids, nil
	}

	groups := make([]GroupModel, 0, len(names))
	for _, n := range names {
		groups = append(groups, GroupModel{Name: n})
	}

	err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).Create(&groups).Error
	if err != nil {
		return nil, err
	}

	var existing []GroupModel
	if err := tx.Where("name IN ?", names).Find(&existing).Error; err != nil {
		return nil, err
	}
	for _, g := range existing {
		ids[g.Name] = g.ID
	}

	return ids, nil
}
