package middleware

import (
    "errors"
    "strings"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "backend/internal/config"
)

var (
    ErrMissingToken = errors.New("missing authentication token")
    ErrInvalidToken = errors.New("invalid authentication token")
)


func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // Skip auth for unauthenticated apis
        if isPublicPath(ctx.Request.URL.Path) {
            ctx.Next()
            return
        }

        userID, err := validateToken(ctx, config.GetConfig().JWTSecret)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            ctx.Abort()
            return
        }

        // Add user ID to context
        ctx.Set("userID", userID)
        ctx.Next()
    }
}

func isPublicPath(path string) bool {
    publicPaths := map[string]bool{
        "/api/login":    true,
        "/api/register": true,
    }
    return publicPaths[path]
}

func validateToken(ctx *gin.Context, secret string) (int, error) {
    tokenString := extractToken(ctx)
    if tokenString == "" {
        return 0, ErrMissingToken
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil || !token.Valid {
        return 0, ErrInvalidToken
    }

    claims, ok := token.Claims.(jwt.MapClaims)

    if !ok {
        return 0, ErrInvalidToken
    }

    userIDFloat, ok := claims["user_id"].(float64)
    if !ok {
        return 0, ErrInvalidToken
    }

    return int(userIDFloat), nil
}

func extractToken(c *gin.Context) string {
    bearerToken := c.Request.Header.Get("Authorization")
    if len(strings.Split(bearerToken, " ")) == 2 {
        return strings.Split(bearerToken, " ")[1]
    }
    return ""
}
