package main

import (
	"time"

	"github.com/ekideno/postly/internal/database"

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
	db, err := database.PostgreConnect(cfg)
	if err != nil {
		panic(err)
	}

	if err = database.Migrate(db.Conn); err != nil {
		panic(err)
	}

	userRepository, err := repository.NewUserRepository(db.Conn)
	if err != nil {
		panic(err)
	}

	jwtManager := security.NewJWTManager("sda*3oj9(FD4)324%34fk#1", time.Hour*24)
	userService := service.NewUserService(userRepository, jwtManager)
	userHandler := handler.NewUserHandler(userService)

	postRepository, err := repository.NewPostRepository(db.Conn)
	if err != nil {
		panic(err)
	}

	postService := service.NewPostService(postRepository)
	postHandler := handler.NewPostHandler(postService, userService)

	r := gin.Default()

	r.Static("/api/uploads", "./uploads")

	api := r.Group("/api")
	api.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	users := api.Group("/users")
	users.Use(jwtManager.OptionalAuthMiddleware())
	{
		users.GET("/:username/posts", postHandler.GetPostsByUser)
		users.GET("/:username", userHandler.UserProfileByUsername)

	}

	protectedUsers := api.Group("/users/@me")
	protectedUsers.Use(jwtManager.AuthMiddleware())
	{
		protectedUsers.GET("", userHandler.GetMe)
		protectedUsers.PATCH("", userHandler.UpdateMe)
		protectedUsers.POST("/avatar", userHandler.UploadAvatar)
		protectedUsers.POST("/banner", userHandler.UploadBanner)
		protectedUsers.GET("/posts", postHandler.PostsForMe)
		protectedUsers.GET("/following", userHandler.GetFollowing)
		protectedUsers.POST("/following/:target_id", userHandler.Follow)
		protectedUsers.DELETE("/following/:target_id", userHandler.Unfollow)

	}

	posts := api.Group("/posts")
	posts.Use(jwtManager.OptionalAuthMiddleware())
	{
		posts.GET("/feed", postHandler.GetFeed)
	}

	protected := api.Group("/")
	protected.Use(jwtManager.AuthMiddleware())
	{
		protected.POST("/posts", postHandler.Create)
	}

	return r
}
