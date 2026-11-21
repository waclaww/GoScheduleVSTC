package handlers

import (
	"errors"
	"gateway/internal/domain"
	"gateway/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type Schedule struct {
	uc *usecase.Schedule
}

func NewSchedule(uc *usecase.Schedule) *Schedule {
	return &Schedule{uc: uc}
}

func (h *Schedule) List(c *fiber.Ctx) error {
	schedule, err := h.uc.List(c.UserContext())
	if err != nil {
		return err
	}

	return c.JSON(toScheduleResponse(schedule))
}

func (h *Schedule) Create(c *fiber.Ctx) error {
	s, err := parseScheduleInput(c)
	if err != nil {
		return err
	}

	if err := h.uc.Create(c.UserContext(), s); err != nil {
		return err
	}

	return c.JSON(done)
}

func (h *Schedule) Update(c *fiber.Ctx) error {
	s, err := parseScheduleInput(c)
	if err != nil {
		return err
	}

	if err := h.uc.UpdateByDate(c.UserContext(), s); err != nil {
		return err
	}

	return c.JSON(done)
}

func (h *Schedule) ByGroup(c *fiber.Ctx) error {
	var req groupNameRequest
	if err := c.BodyParser(&req); err != nil {
		return unprocessable(err.Error())
	}
	if req.GroupName == nil || !groupNamePattern.MatchString(*req.GroupName) {
		return unprocessable("group_name: invalid format")
	}

	schedule, err := h.uc.ByGroup(c.UserContext(), *req.GroupName)
	if errors.Is(err, domain.ErrGroupNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(groupScheduleResponse{Today: schedule.Today, Tomorrow: schedule.Tomorrow})
}

func parseScheduleInput(c *fiber.Ctx) (*domain.Schedule, error) {
	var req scheduleInput
	if err := c.BodyParser(&req); err != nil {
		return nil, unprocessable(err.Error())
	}
	if req.Date == nil || req.Schedule == nil {
		return nil, unprocessable("date and schedule are required")
	}

	return &domain.Schedule{Date: *req.Date, Schedule: *req.Schedule}, nil
}
