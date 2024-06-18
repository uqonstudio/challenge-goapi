package middleware

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(c *gin.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path
	header := c.Request.Header
	// log the request details
	fmt.Printf("Req : %s %s\nHeader : %v", method, path, header)
	c.Next()
}

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}
	tokenString := authHeader[7:] //remove "Bearer" prefix from token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("@Enigma2024"), nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return
	}

	// cek user permission
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if claims["permission"] != "admin" {
		c.JSON(401, gin.H{"error": "Permission denied"})
		c.Abort()
		return
	}

	c.Next()
}
