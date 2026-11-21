package handlers

import (
	"gateway/internal/domain"
	"time"
)

type pageInput struct {
	Date *string `json:"date"`
}

type pageResultResponse struct {
	Name         string  `json:"name"`
	Date         string  `json:"date"`
	PreviousDate *string `json:"previous_date"`
	Changed      bool    `json:"changed"`
}

type pageResponse struct {
	Name      string    `json:"name"`
	Date      string    `json:"date"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

type pageLessonsResponse struct {
	Date    string           `json:"date"`
	Title   string           `json:"title"`
	Lessons []groupLessonDTO `json:"lessons"`
}

type groupPagesResponse struct {
	Group    string               `json:"group"`
	Today    *pageLessonsResponse `json:"today"`
	Tomorrow *pageLessonsResponse `json:"tomorrow"`
}

func toGroupPagesResponse(g *domain.GroupPages) groupPagesResponse {
	return groupPagesResponse{Group: g.Group, Today: toPageLessonsResponse(g.Today), Tomorrow: toPageLessonsResponse(g.Tomorrow)}
}

func toPageLessonsResponse(d *domain.GroupDay) *pageLessonsResponse {
	if d == nil {
		return nil
	}

	day := toGroupDayResponse(d)
	return &pageLessonsResponse{Date: day.Date, Title: day.Title, Lessons: day.Lessons}
}

func toPageResultResponse(r domain.PageResult) pageResultResponse {
	res := pageResultResponse{Name: r.Name, Date: r.Date.Format(time.DateOnly), Changed: r.Changed}
	if r.PreviousDate != nil {
		previous := r.PreviousDate.Format(time.DateOnly)
		res.PreviousDate = &previous
	}

	return res
}

func toPageResponse(p *domain.Page) pageResponse {
	return pageResponse{Name: p.Name, Date: p.Date.Format(time.DateOnly), Title: p.Title, UpdatedAt: p.UpdatedAt}
}
