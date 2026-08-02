package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todo-api/internal/dto"
	"todo-api/internal/service"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", err.Error()))
		return
	}

	result, err := h.todoService.Create(userID, req)
	if err != nil {
		handleTodoError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Todo created successfully", result))
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	limit := parsePositiveInt(c.DefaultQuery("limit", "10"), 10)
	search := c.Query("search")

	result, err := h.todoService.GetAll(userID, page, limit, search)
	if err != nil {
		handleTodoError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Todos retrieved successfully", result))
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	todoID, ok := parseIDParam(c)
	if !ok {
		return
	}

	result, err := h.todoService.GetByID(userID, todoID)
	if err != nil {
		handleTodoError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Todo retrieved successfully", result))
}

func (h *TodoHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	todoID, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", err.Error()))
		return
	}

	result, err := h.todoService.Update(userID, todoID, req)
	if err != nil {
		handleTodoError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Todo updated successfully", result))
}

func (h *TodoHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	todoID, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := h.todoService.Delete(userID, todoID); err != nil {
		handleTodoError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Todo deleted successfully", nil))
}

func getUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Unauthorized", nil))
		return 0, false
	}

	userID, ok := value.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user", nil))
		return 0, false
	}
	return userID, true
}

func parseIDParam(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid todo ID", nil))
		return 0, false
	}
	return uint(value), true
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func handleTodoError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTodoNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error(), nil))
	case errors.Is(err, service.ErrTodoForbidden):
		c.JSON(http.StatusForbidden, dto.ErrorResponse(err.Error(), nil))
	case errors.Is(err, service.ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), nil))
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Internal server error", nil))
	}
}
