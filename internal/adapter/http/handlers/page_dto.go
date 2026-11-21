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
