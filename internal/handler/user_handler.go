package handler

import (
	"fmt"
	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/service"
	"github.com/ekideno/postly/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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

	user := domain.FromRegisterRequest(&req)

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

	public := domain.ToPublicUserDTO(user)

	c.JSON(http.StatusOK, public)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := c.Get("user_id")
	fmt.Println(userID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.UserService.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	private := domain.ToPrivateUserDTO(user)

	c.JSON(http.StatusOK, private)
}
func (h *UserHandler) UserProfileByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.UserService.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	public := domain.ToPublicUserDTO(user)

	c.JSON(http.StatusOK, public)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userIDRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	var dto domain.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	updatedUser, err := h.UserService.UpdateUserProfile(userID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := domain.ToPrivateUserDTO(updatedUser)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userIDRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDRaw.(string)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	uploadPath := "./uploads/avatars"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create upload dir"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", utils.GenerateUUID(), ext)
	fullPath := filepath.Join(uploadPath, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	avatarPath := "/uploads/avatars/" + filename
	err = h.UserService.UpdateAvatar(userID, avatarPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update avatar in DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarPath})
}

func (h *UserHandler) UploadBanner(c *gin.Context) {
	userIDRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	userID, _ := userIDRaw.(string)
	file, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	uploadPath := "./uploads/banners"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create upload dir"})
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", utils.GenerateUUID(), ext)
	fullPath := filepath.Join(uploadPath, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	bannerPath := "/uploads/banners/" + filename
	err = h.UserService.UpdateBanner(userID, bannerPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update avatar in DB"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"banner_url": bannerPath})
}
