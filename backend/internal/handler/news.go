package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

type NewsHandler struct {
	service domain.NewsService
}

func NewNewsHandler(service domain.NewsService) *NewsHandler {
	return &NewsHandler{service: service}
}

func (h *NewsHandler) Create(c *fiber.Ctx) error {
	var news domain.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.Create(&news); err != nil {
		if err == domain.ErrNewsTitleRequired || err == domain.ErrNewsContentRequired || err == domain.ErrNewsAuthorRequired {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(news)
}

func (h *NewsHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	news, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News not found",
		})
	}

	return c.JSON(news)
}

func (h *NewsHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news domain.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	news.ID = id
	if err := h.service.Update(&news); err != nil {
		if err == domain.ErrNewsTitleRequired || err == domain.ErrNewsContentRequired || err == domain.ErrNewsAuthorRequired {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(news)
}

func (h *NewsHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NewsHandler) List(c *fiber.Ctx) error {
	news, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(news)
}
