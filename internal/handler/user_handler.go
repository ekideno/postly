package handler

import (
	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	token, err := h.UserService.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.UserService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) UserProfileByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UserService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Username not found"})
		return
	}

	public := domain.PublicUserDTO{
		ID:       user.ID,
		Username: user.Username,
	}

	c.JSON(http.StatusOK, public)
}

func (h *UserHandler) OwnProfile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.UserService.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
func (h *UserHandler) UserProfileByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.UserService.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return // важно не забыть return, чтобы не пойти дальше
	}

	currentUserID, exists := c.Get("user_id")

	if exists && currentUserID == user.ID {
		private := domain.PrivateUserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		c.JSON(http.StatusOK, private)
		return
	}

	public := domain.PublicUserDTO{
		ID:       user.ID,
		Username: user.Username,
	}
	c.JSON(http.StatusOK, public)
}
