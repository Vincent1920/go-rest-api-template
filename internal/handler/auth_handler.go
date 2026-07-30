package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"todo-api/internal/dto"
	"todo-api/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register User
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			dto.ErrorResponse("Validation failed", err.Error()))
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated,
		dto.SuccessResponse("User registered successfully", user))
}

// Login godoc
// @Summary Login
// @Description Login user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			dto.ErrorResponse("Validation failed", err.Error()))
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			dto.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK,
		dto.SuccessResponse("Login successful", result))
}

// GetProfile godoc
// @Summary Get Profile
// @Description Get current user profile
// @Tags Authentication
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized,
			dto.ErrorResponse("Unauthorized", nil))
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			dto.ErrorResponse("Invalid user", nil))
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK,
		dto.SuccessResponse("Success", user))
}
