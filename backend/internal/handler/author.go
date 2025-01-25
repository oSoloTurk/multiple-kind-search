package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

type AuthorHandler struct {
	service domain.AuthorService
}

func NewAuthorHandler(service domain.AuthorService) *AuthorHandler {
	return &AuthorHandler{service: service}
}

// Create godoc
// @Summary Create a new author
// @Description Create a new author with the provided details
// @Tags authors
// @Accept json
// @Produce json
// @Param author body domain.Author true "Author details"
// @Success 201 {object} domain.Author
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/authors [post]
func (h *AuthorHandler) Create(c *fiber.Ctx) error {
	var author domain.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.Create(&author); err != nil {
		if err == domain.ErrAuthorNameRequired {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(author)
}

// GetByID godoc
// @Summary Get an author by ID
// @Description Get an author's details by their ID
// @Tags authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} domain.Author
// @Failure 404 {object} map[string]string
// @Router /api/authors/{id} [get]
func (h *AuthorHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	author, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	return c.JSON(author)
}

// Update godoc
// @Summary Update an author
// @Description Update an existing author's details
// @Tags authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Param author body domain.Author true "Updated author details"
// @Success 200 {object} domain.Author
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/authors/{id} [put]
func (h *AuthorHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var author domain.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	author.ID = id
	if err := h.service.Update(&author); err != nil {
		if err == domain.ErrAuthorNameRequired {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(author)
}

// Delete godoc
// @Summary Delete an author
// @Description Delete an author by their ID
// @Tags authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Router /api/authors/{id} [delete]
func (h *AuthorHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List godoc
// @Summary List all authors
// @Description Get a list of all authors
// @Tags authors
// @Accept json
// @Produce json
// @Success 200 {array} domain.Author
// @Failure 500 {object} map[string]string
// @Router /api/authors [get]
func (h *AuthorHandler) List(c *fiber.Ctx) error {
	authors, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(authors)
}

// ... rest of the handler methods remain similar but use service instead of repo
