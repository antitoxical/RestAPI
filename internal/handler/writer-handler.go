package handler

import (
	"net/http"
	"strconv"
	"strings"

	"RESTAPI/internal/dto"
	"RESTAPI/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type WriterHandler struct {
	service *service.WriterService
}

func NewWriterHandler(service *service.WriterService) *WriterHandler {
	return &WriterHandler{service: service}
}

func (h *WriterHandler) Create(c echo.Context) error {
	var req dto.WriterRequestTo
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := validator.New().Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.Create(req)
	if err != nil {
		if err.Error() == "login_already_exists" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Login already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *WriterHandler) GetById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	writer, err := h.service.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Writer not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, writer)
}

func (h *WriterHandler) Update(c echo.Context) error {
	/*id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}*/
	var req dto.WriterUpdateRequestTo
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
		if err.Error() == "writer not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Writer not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *WriterHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	err = h.service.Delete(id)
	if err != nil {
		if err.Error() == "writer not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Writer not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *WriterHandler) GetAll(c echo.Context) error {
	writers, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, writers)
}
