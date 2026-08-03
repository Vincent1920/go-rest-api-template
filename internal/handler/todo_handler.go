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
	return &TodoHandler{
		todoService: todoService,
	}
}

// Create godoc
// @Summary Create Todo
// @Description Create a new Todo for the authenticated user
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTodoRequest true "Create Todo Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req dto.CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse("Validation failed", err.Error()),
		)
		return
	}

	result, err := h.todoService.Create(userID, req)
	if err != nil {
		handleTodoServiceError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		dto.SuccessResponse("Todo created successfully", result),
	)
}

// GetAll godoc
// @Summary Get Todos
// @Description Get Todo list for the authenticated user
// @Tags Todos
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search Todo title"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /todos [get]
func (h *TodoHandler) GetAll(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 10)
	search := c.Query("search")

	result, err := h.todoService.GetAll(
		userID,
		page,
		limit,
		search,
	)
	if err != nil {
		handleTodoServiceError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.SuccessResponse("Todos retrieved successfully", result),
	)
}

// GetByID godoc
// @Summary Get Todo Detail
// @Description Get a Todo by ID for the authenticated user
// @Tags Todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	todoID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse("Invalid Todo ID", nil),
		)
		return
	}

	result, err := h.todoService.GetByID(userID, todoID)
	if err != nil {
		handleTodoServiceError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.SuccessResponse("Todo retrieved successfully", result),
	)
}

// Update godoc
// @Summary Update Todo
// @Description Update a Todo owned by the authenticated user
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param request body dto.UpdateTodoRequest true "Update Todo Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	todoID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse("Invalid Todo ID", nil),
		)
		return
	}

	var req dto.UpdateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse("Validation failed", err.Error()),
		)
		return
	}

	result, err := h.todoService.Update(
		userID,
		todoID,
		req,
	)
	if err != nil {
		handleTodoServiceError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.SuccessResponse("Todo updated successfully", result),
	)
}

// Delete godoc
// @Summary Delete Todo
// @Description Delete a Todo owned by the authenticated user
// @Tags Todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	todoID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse("Invalid Todo ID", nil),
		)
		return
	}

	if err := h.todoService.Delete(userID, todoID); err != nil {
		handleTodoServiceError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.SuccessResponse("Todo deleted successfully", nil),
	)
}

func getUserIDFromContext(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			dto.ErrorResponse("Unauthorized", nil),
		)
		return 0, false
	}

	userID, ok := value.(uint)
	if !ok || userID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			dto.ErrorResponse("Invalid user", nil),
		)
		return 0, false
	}

	return userID, true
}

func parseIDParam(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid ID")
	}

	return uint(id), nil
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil || result < 1 {
		return fallback
	}

	return result
}

func handleTodoServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTodoNotFound):
		c.JSON(
			http.StatusNotFound,
			dto.ErrorResponse(err.Error(), nil),
		)

	case errors.Is(err, service.ErrInvalidTodoStatus):
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse(err.Error(), nil),
		)

	case errors.Is(err, service.ErrEmptyTodoTitle):
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse(err.Error(), nil),
		)

	default:
		c.JSON(
			http.StatusInternalServerError,
			dto.ErrorResponse("Internal server error", nil),
		)
	}
}
