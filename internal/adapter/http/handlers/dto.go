package handlers

import (
	"gateway/internal/domain"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

var groupNamePattern = regexp.MustCompile(`^[А-Я]{1,3}-\d{2,3}$`)

var dateTimeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04",
	time.DateOnly,
}

type statusResponse struct {
	Status string `json:"status"`
}

var done = statusResponse{Status: "done"}

type scheduleInput struct {
	Date     *string `json:"date"`
	Schedule *string `json:"schedule"`
}

type archiveInput struct {
	Date     *string `json:"date"`
	Schedule *string `json:"schedule"`
}

type groupNameRequest struct {
	GroupName *string `json:"group_name"`
}

type scheduleResponse struct {
	Id       int64  `json:"id"`
	Schedule string `json:"schedule"`
	Date     string `json:"date"`
}

type archiveResponse struct {
	Id       int64  `json:"id"`
	Schedule string `json:"schedule"`
	Date     string `json:"date"`
}

type archiveByDateResponse struct {
	Status   string            `json:"status"`
	Schedule []archiveResponse `json:"schedule"`
}

type groupScheduleResponse struct {
	Today    string `json:"today"`
	Tomorrow string `json:"tomorrow"`
}

func unprocessable(msg string) error {
	return fiber.NewError(fiber.StatusUnprocessableEntity, msg)
}

func parseDateTime(value string) (time.Time, error) {
	var err error
	for _, layout := range dateTimeLayouts {
		var t time.Time
		if t, err = time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, err
}

func toScheduleResponse(s []*domain.Schedule) []scheduleResponse {
	res := make([]scheduleResponse, 0, len(s))
	for _, v := range s {
		res = append(res, scheduleResponse{Id: v.Id, Schedule: v.Schedule, Date: v.Date})
	}

	return res
}

func toArchiveResponse(a []*domain.Archive) []archiveResponse {
	res := make([]archiveResponse, 0, len(a))
	for _, v := range a {
		res = append(res, archiveResponse{Id: v.Id, Schedule: v.Schedule, Date: v.Date})
	}

	return res
}
