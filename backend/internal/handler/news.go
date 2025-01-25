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

// Create godoc
// @Summary Create a new news article
// @Description Create a new news article with the provided details
// @Tags news
// @Accept json
// @Produce json
// @Param news body domain.News true "News article details"
// @Success 201 {object} domain.News
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/news [post]
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

// GetByID godoc
// @Summary Get a news article by ID
// @Description Get a news article's details by its ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} domain.News
// @Failure 404 {object} map[string]string
// @Router /api/news/{id} [get]
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

// Update godoc
// @Summary Update a news article
// @Description Update an existing news article's details
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Param news body domain.News true "Updated news article details"
// @Success 200 {object} domain.News
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/news/{id} [put]
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

// Delete godoc
// @Summary Delete a news article
// @Description Delete a news article by its ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Router /api/news/{id} [delete]
func (h *NewsHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List godoc
// @Summary List all news articles
// @Description Get a list of all news articles
// @Tags news
// @Accept json
// @Produce json
// @Success 200 {array} domain.News
// @Failure 500 {object} map[string]string
// @Router /api/news [get]
func (h *NewsHandler) List(c *fiber.Ctx) error {
	news, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(news)
}
