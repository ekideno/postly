package main

import (
	"github.com/ekideno/postly/internal/handler"
	"github.com/ekideno/postly/internal/repository"
	"github.com/ekideno/postly/internal/security"
	"github.com/ekideno/postly/internal/service"
	"github.com/gin-gonic/gin"
	"time"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "Postly!"})
	})

	userRepository, err := repository.NewUserRepository()
	if err != nil {
		panic(err)
	}
	jwtManager := security.NewJWTManager("sda*3oj9(FD4)324%34fk#1", time.Hour*24)

	userService := service.NewUserService(userRepository, jwtManager)
	userHandler := handler.NewUserHandler(userService)

	api := r.Group("/api")
	auth := api.Group("/auth")

	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)

	protected := api.Group("/")
	protected.Use(jwtManager.AuthMiddleware())

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
