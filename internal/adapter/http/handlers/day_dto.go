package handlers

import (
	"gateway/internal/domain"
	"net/url"
	"strings"
	"time"
)

type slotDTO struct {
	Slot     int    `json:"slot"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type lessonDTO struct {
	Slot    int    `json:"slot"`
	Subject string `json:"subject"`
	Room    string `json:"room"`
}

type dayInput struct {
	Title  *string                `json:"title"`
	Slots  []slotDTO              `json:"slots"`
	Groups map[string][]lessonDTO `json:"groups"`
}

type dayResponse struct {
	Date      string                 `json:"date"`
	Title     string                 `json:"title"`
	FetchedAt time.Time              `json:"fetched_at"`
	Slots     []slotDTO              `json:"slots"`
	Groups    map[string][]lessonDTO `json:"groups"`
}

type upsertResponse struct {
	Date    string `json:"date"`
	Created bool   `json:"created"`
	Changed bool   `json:"changed"`
}

type groupLessonDTO struct {
	Slot     int    `json:"slot"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	Subject  string `json:"subject"`
	Room     string `json:"room"`
}

type groupDayResponse struct {
	Date    string           `json:"date"`
	Title   string           `json:"title"`
	Group   string           `json:"group"`
	Lessons []groupLessonDTO `json:"lessons"`
}

func parseDate(value string) (time.Time, error) {
	return time.Parse(time.DateOnly, value)
}

func groupParam(value string) (string, error) {
	name, err := url.PathUnescape(value)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(name), nil
}

func (in dayInput) toDomain(date time.Time) *domain.Day {
	d := &domain.Day{
		Date:   date,
		Title:  *in.Title,
		Groups: make([]string, 0, len(in.Groups)),
		Slots:  make([]domain.BellSlot, 0, len(in.Slots)),
	}

	for _, s := range in.Slots {
		d.Slots = append(d.Slots, domain.BellSlot{Slot: s.Slot, StartsAt: s.StartsAt, EndsAt: s.EndsAt})
	}

	for group, lessons := range in.Groups {
		d.Groups = append(d.Groups, group)
		for _, l := range lessons {
			d.Lessons = append(d.Lessons, domain.Lesson{Group: group, Slot: l.Slot, Subject: l.Subject, Room: l.Room})
		}
	}

	return d
}

func toDayResponse(d *domain.Day) dayResponse {
	res := dayResponse{
		Date:      d.Date.Format(time.DateOnly),
		Title:     d.Title,
		FetchedAt: d.FetchedAt,
		Slots:     make([]slotDTO, 0, len(d.Slots)),
		Groups:    make(map[string][]lessonDTO, len(d.Groups)),
	}

	for _, s := range d.Slots {
		res.Slots = append(res.Slots, slotDTO{Slot: s.Slot, StartsAt: s.StartsAt, EndsAt: s.EndsAt})
	}

	for _, l := range d.Lessons {
		res.Groups[l.Group] = append(res.Groups[l.Group], lessonDTO{Slot: l.Slot, Subject: l.Subject, Room: l.Room})
	}

	return res
}

func toGroupDayResponse(g *domain.GroupDay) groupDayResponse {
	res := groupDayResponse{
		Date:    g.Date.Format(time.DateOnly),
		Title:   g.Title,
		Group:   g.Group,
		Lessons: make([]groupLessonDTO, 0, len(g.Lessons)),
	}

	for _, l := range g.Lessons {
		res.Lessons = append(res.Lessons, groupLessonDTO{Slot: l.Slot, StartsAt: l.StartsAt, EndsAt: l.EndsAt, Subject: l.Subject, Room: l.Room})
	}

	return res
}
