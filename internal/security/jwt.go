package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type JWTManager struct {
	secretKey []byte
	expiresIn time.Duration
}

func NewJWTManager(secret string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secret),
		expiresIn: expiresIn,
	}
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (j *JWTManager) GenerateToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWTManager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func (j *JWTManager) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Нет токена — просто идём дальше, без ошибки
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.secretKey, nil
		})

		if err != nil || !token.Valid {
			// Токен есть, но не валидный — тоже пропускаем без ошибки,
			// чтобы гость мог смотреть публичные страницы
			c.Next()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.Next()
			return
		}

		// Сохраняем user_id в контекст для последующего использования
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
