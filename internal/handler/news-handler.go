package handler

import (
	"net/http"
	"strconv"

	"RESTAPI/internal/dto"
	"RESTAPI/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type NewsHandler struct {
	service *service.NewsService
}

func NewNewsHandler(service *service.NewsService) *NewsHandler {
	return &NewsHandler{service: service}
}

func (h *NewsHandler) Create(c echo.Context) error {
	var req dto.NewsRequestTo
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if err := validator.New().Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.service.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *NewsHandler) GetById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	resp, err := h.service.GetById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *NewsHandler) Update(c echo.Context) error {
	/*id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}*/

	var req dto.NewsUpdateRequestTo
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	var Validate = validator.New()

	// Валидация входных данных
	if err := Validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.Update(req)
	if err != nil {
		if err.Error() == "news not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "News not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *NewsHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	err = h.service.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "News not found"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *NewsHandler) GetAll(c echo.Context) error {
	newsList, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, newsList)
}
