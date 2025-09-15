package v1

import (
	"errors"
	"net/http"
	"test_go/internal/controller/http/v1/request"
	"test_go/internal/entity"
	"test_go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService usecase.UserService
}

func NewAuthHandler(userService usecase.UserService) *AuthHandler {
	return &AuthHandler{authService: userService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user entity.UserInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.AuthenticateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, entity.ErrAccessDenied):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest

	userIdAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userId, ok := userIdAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userId,
		req.OldPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пароль успешно изменен"})
}

func (h *AuthHandler) SetProfilePhoto(c *gin.Context) {
	userIdAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userId, ok := userIdAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// Получаем файл из формы
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	err = h.authService.SetProfilePhoto(c.Request.Context(), userId, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile photo updated successfully"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
