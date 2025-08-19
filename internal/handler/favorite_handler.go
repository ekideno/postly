package handler

import (
	"net/http"
	"strconv"

	"github.com/ekideno/postly/internal/service"
	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	FavoriteService *service.FavoriteService
}

func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{FavoriteService: favoriteService}
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	imageID := c.Param("imageID")

	if err := h.FavoriteService.AddFavorite(userID.(string), imageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to favorites"})
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	imageID := c.Param("imageID")

	if err := h.FavoriteService.RemoveFavorite(userID.(string), imageID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to remove favorite"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "removed from favorites"})
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	images, err := h.FavoriteService.GetFavorites(userID.(string), limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch favorites"})
		return
	}

	c.JSON(200, images)
}
