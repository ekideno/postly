package main

import (
	"time"

	"github.com/ekideno/postly/internal/config"
	"github.com/ekideno/postly/internal/handler"
	"github.com/ekideno/postly/internal/repository"
	"github.com/ekideno/postly/internal/security"
	"github.com/ekideno/postly/internal/service"
	"github.com/ekideno/postly/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()

	r := setupRouter(cfg)

	utils.InitSnowflake(1)
	r.Run(":8080")
}

func setupRouter(cfg *config.Config) *gin.Engine {
	// Initialize components
	userRepository, err := repository.NewUserRepository(cfg)
	if err != nil {
		panic(err)
	}
	jwtManager := security.NewJWTManager("sda*3oj9(FD4)324%34fk#1", time.Hour*24)
	userService := service.NewUserService(userRepository, jwtManager)
	userHandler := handler.NewUserHandler(userService)

	postRepository, err := repository.NewPostRepository()
	if err != nil {
		panic(err)
	}
	postService := service.NewPostService(postRepository)
	postHandler := handler.NewPostHandler(postService)

	// Initialize router
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "Postly!"})
	})

	// Routes
	api := r.Group("/api")
	auth := api.Group("/auth")
	protected := api.Group("/")

	api.GET("/user/profile", jwtManager.AuthMiddleware(), userHandler.OwnProfile)
	api.GET("/user/:id/profile", userHandler.UserProfile)
	api.GET("/user/:id/posts", postHandler.GetPostsByUser)

	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)

	protected.Use(jwtManager.AuthMiddleware())
	protected.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "111!"})
	})
	protected.POST("/post", postHandler.Create)

	return r
}
