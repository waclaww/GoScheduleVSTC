package handlers

import (
	"errors"
	"gateway/internal/domain"
	"gateway/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Day struct {
	uc *usecase.Day
}

func NewDay(uc *usecase.Day) *Day {
	return &Day{uc: uc}
}

func (h *Day) Upsert(c *fiber.Ctx) error {
	date, err := parseDate(c.Params("date"))
	if err != nil {
		return unprocessable("date: expected YYYY-MM-DD")
	}

	var req dayInput
	if err := c.BodyParser(&req); err != nil {
		return unprocessable(err.Error())
	}
	if req.Title == nil {
		return unprocessable("title is required")
	}

	res, err := h.uc.Upsert(c.UserContext(), req.toDomain(date))
	if errors.Is(err, domain.ErrInvalidDay) {
		return unprocessable(err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(upsertResponse{Date: date.Format(time.DateOnly), Created: res.Created, Changed: res.Changed})
}

func (h *Day) Get(c *fiber.Ctx) error {
	date, err := parseDate(c.Params("date"))
	if err != nil {
		return unprocessable("date: expected YYYY-MM-DD")
	}

	day, err := h.uc.Get(c.UserContext(), date)
	if errors.Is(err, domain.ErrDayNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(toDayResponse(day))
}

func (h *Day) Groups(c *fiber.Ctx) error {
	groups, err := h.uc.Groups(c.UserContext())
	if err != nil {
		return err
	}

	return c.JSON(groups)
}

func (h *Day) GroupSchedule(c *fiber.Ctx) error {
	group, err := groupParam(c.Params("name"))
	if err != nil || group == "" {
		return unprocessable("name: invalid group")
	}

	date := h.uc.Today()
	if raw := c.Query("date"); raw != "" {
		if date, err = parseDate(raw); err != nil {
			return unprocessable("date: expected YYYY-MM-DD")
		}
	}

	schedule, err := h.uc.GroupDay(c.UserContext(), group, date)
	if errors.Is(err, domain.ErrGroupNotFound) || errors.Is(err, domain.ErrDayNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(toGroupDayResponse(schedule))
}
