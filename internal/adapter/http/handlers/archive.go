package handlers

import (
	"gateway/internal/domain"
	"gateway/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Archive struct {
	uc *usecase.Archive
}

func NewArchive(uc *usecase.Archive) *Archive {
	return &Archive{uc: uc}
}

func (h *Archive) List(c *fiber.Ctx) error {
	archive, err := h.uc.List(c.UserContext())
	if err != nil {
		return err
	}

	return c.JSON(toArchiveResponse(archive))
}

func (h *Archive) Create(c *fiber.Ctx) error {
	var req archiveInput
	if err := c.BodyParser(&req); err != nil {
		return unprocessable(err.Error())
	}
	if req.Date == nil || req.Schedule == nil {
		return unprocessable("date and schedule are required")
	}

	date, err := parseDateTime(*req.Date)
	if err != nil {
		return unprocessable("date: invalid datetime")
	}

	a := &domain.Archive{Date: date.Format(time.DateOnly), Schedule: *req.Schedule}
	if err := h.uc.Create(c.UserContext(), a); err != nil {
		return err
	}

	return c.JSON(done)
}

func (h *Archive) ByDate(c *fiber.Ctx) error {
	raw := c.Query("data")
	if raw == "" {
		return unprocessable("data: field required")
	}

	date, err := parseDateTime(raw)
	if err != nil {
		return unprocessable("data: invalid datetime")
	}

	archive, err := h.uc.ListByDate(c.UserContext(), date.Format(time.DateOnly))
	if err != nil {
		return err
	}

	return c.JSON(archiveByDateResponse{Status: done.Status, Schedule: toArchiveResponse(archive)})
}
