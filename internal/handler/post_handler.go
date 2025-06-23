package handler

import (
	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"strconv"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(string)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}
	files := form.File["images"]

	var imagePaths []string
	for _, file := range files {
		filename := uuid.NewString() + filepath.Ext(file.Filename)
		savePath := "./uploads/images/" + filename

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imagePaths = append(imagePaths, "/uploads/images/"+filename)
	}

	post, err := h.postService.Create(userID, &domain.CreatePostRequest{
		Title:   title,
		Content: content,
		Images:  imagePaths,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ToPostResponse(post)
	c.JSON(http.StatusCreated, gin.H{"post": response})
}

func (h *PostHandler) GetPostsByUser(c *gin.Context) {
	username := c.Param("username")

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	posts, err := h.postService.GetPostsByUser(username, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	responses := ToPostResponseList(posts)

	c.JSON(http.StatusOK, gin.H{
		"posts":  responses,
		"limit":  limit,
		"offset": offset,
	})
}

func ToPostResponseList(posts []domain.Post) []domain.PostResponse {
	result := make([]domain.PostResponse, 0, len(posts))
	for _, p := range posts {
		result = append(result, ToPostResponse(&p))
	}
	return result
}

func ToPostResponse(post *domain.Post) domain.PostResponse {
	var imageURLs []string
	for _, img := range post.Images {
		imageURLs = append(imageURLs, img.URL)
	}

	return domain.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Author: domain.PublicUserDTO{
			ID:        post.User.ID,
			Username:  post.User.Username,
			Bio:       post.User.Bio,
			AvatarURL: post.User.AvatarURL,
			BannerURL: post.User.BannerURL,
		},
		Images: imageURLs,
	}
}

func (h *PostHandler) GetFeed(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	posts, err := h.postService.GetFeed(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load feed"})
		return
	}

	postResponses := ToPostResponseList(posts)

	c.JSON(http.StatusOK, postResponses)
}

func (h *PostHandler) PostsForMe(c *gin.Context) {
	usernameRaw, ok := c.Get("username")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	username, _ := usernameRaw.(string)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	posts, err := h.postService.GetPostsByUser(username, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load posts"})
		return
	}

	postResponses := ToPostResponseList(posts)
	c.JSON(http.StatusOK, postResponses)

}
