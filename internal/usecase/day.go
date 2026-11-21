package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gateway/internal/domain"
	"slices"
	"strings"
	"time"
)

const (
	maxSlot      = 32
	maxGroupName = 32
)

type Day struct {
	repo DayRepository
}

func NewDay(repo DayRepository) *Day {
	return &Day{repo: repo}
}

func (u *Day) Upsert(ctx context.Context, d *domain.Day) (domain.UpsertResult, error) {
	if err := normalizeDay(d); err != nil {
		return domain.UpsertResult{}, err
	}

	d.SourceHash = sourceHash(d)
	d.FetchedAt = time.Now()

	return u.repo.Upsert(ctx, d)
}

func (u *Day) Get(ctx context.Context, date time.Time) (*domain.Day, error) {
	return u.repo.Get(ctx, date)
}

func (u *Day) Groups(ctx context.Context) ([]string, error) {
	return u.repo.Groups(ctx)
}

func (u *Day) GroupDay(ctx context.Context, group string, date time.Time) (*domain.GroupDay, error) {
	return u.repo.GroupDay(ctx, group, date)
}

func (u *Day) Today() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func normalizeDay(d *domain.Day) error {
	d.Title = strings.TrimSpace(d.Title)
	if d.Title == "" {
		return fmt.Errorf("%w: title is required", domain.ErrInvalidDay)
	}

	slotSeen := make(map[int]bool, len(d.Slots))
	for _, s := range d.Slots {
		if s.Slot < 1 || s.Slot > maxSlot {
			return fmt.Errorf("%w: slot %d out of range", domain.ErrInvalidDay, s.Slot)
		}
		if slotSeen[s.Slot] {
			return fmt.Errorf("%w: duplicate slot %d", domain.ErrInvalidDay, s.Slot)
		}
		if !validClock(s.StartsAt) || !validClock(s.EndsAt) {
			return fmt.Errorf("%w: slot %d has invalid time", domain.ErrInvalidDay, s.Slot)
		}
		slotSeen[s.Slot] = true
	}

	groupSeen := make(map[string]bool, len(d.Groups))
	groups := make([]string, 0, len(d.Groups))
	for _, g := range d.Groups {
		g = strings.TrimSpace(g)
		if g == "" || len([]rune(g)) > maxGroupName {
			return fmt.Errorf("%w: invalid group name %q", domain.ErrInvalidDay, g)
		}
		if !groupSeen[g] {
			groupSeen[g] = true
			groups = append(groups, g)
		}
	}

	lessonSeen := make(map[string]bool, len(d.Lessons))
	for i := range d.Lessons {
		l := &d.Lessons[i]
		l.Group = strings.TrimSpace(l.Group)
		l.Subject = strings.TrimSpace(l.Subject)
		l.Room = strings.TrimSpace(l.Room)
		if !groupSeen[l.Group] {
			return fmt.Errorf("%w: unknown group %q", domain.ErrInvalidDay, l.Group)
		}
		if l.Slot < 1 || l.Slot > maxSlot {
			return fmt.Errorf("%w: slot %d out of range", domain.ErrInvalidDay, l.Slot)
		}
		key := fmt.Sprintf("%s/%d", l.Group, l.Slot)
		if lessonSeen[key] {
			return fmt.Errorf("%w: duplicate lesson %s", domain.ErrInvalidDay, key)
		}
		lessonSeen[key] = true
	}

	d.Lessons = slices.DeleteFunc(d.Lessons, func(l domain.Lesson) bool {
		return l.Subject == "" && l.Room == ""
	})

	slices.Sort(groups)
	d.Groups = groups
	slices.SortFunc(d.Slots, func(a, b domain.BellSlot) int {
		return a.Slot - b.Slot
	})
	slices.SortFunc(d.Lessons, func(a, b domain.Lesson) int {
		if c := strings.Compare(a.Group, b.Group); c != 0 {
			return c
		}
		return a.Slot - b.Slot
	})

	return nil
}

func validClock(v string) bool {
	if v == "" {
		return true
	}
	_, err := time.Parse("15:04", v)
	return err == nil
}

func sourceHash(d *domain.Day) string {
	payload, _ := json.Marshal(struct {
		Title   string
		Groups  []string
		Slots   []domain.BellSlot
		Lessons []domain.Lesson
	}{d.Title, d.Groups, d.Slots, d.Lessons})

	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
