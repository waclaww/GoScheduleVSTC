package handlers

import (
	"errors"
	"gateway/internal/domain"
	"gateway/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Page struct {
	uc *usecase.Page
}

func NewPage(uc *usecase.Page) *Page {
	return &Page{uc: uc}
}

func (h *Page) Set(c *fiber.Ctx) error {
	var req pageInput
	if err := c.BodyParser(&req); err != nil {
		return unprocessable(err.Error())
	}
	if req.Date == nil {
		return unprocessable("date is required")
	}

	date, err := parseDate(*req.Date)
	if err != nil {
		return unprocessable("date: expected YYYY-MM-DD")
	}

	res, err := h.uc.Set(c.UserContext(), c.Params("name"), date)
	if errors.Is(err, domain.ErrInvalidPage) || errors.Is(err, domain.ErrDayNotFound) {
		return unprocessable(err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(toPageResultResponse(res))
}

func (h *Page) Get(c *fiber.Ctx) error {
	page, err := h.uc.Get(c.UserContext(), c.Params("name"))
	if errors.Is(err, domain.ErrPageNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(toPageResponse(page))
}

func (h *Page) GroupSchedule(c *fiber.Ctx) error {
	group := strings.TrimSpace(c.Query("group"))
	if group == "" {
		return unprocessable("group is required")
	}

	schedule, err := h.uc.GroupSchedule(c.UserContext(), group)
	if errors.Is(err, domain.ErrGroupNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(toGroupPagesResponse(schedule))
}
