package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type authClient interface {
	ValidateToken(ctx context.Context, token string) error
}

type Authorization struct {
	authClient authClient
	logger     *zap.Logger
}

func NewAuthorization(authClient authClient, logger *zap.Logger) Authorization {
	return Authorization{
		authClient: authClient,
		logger:     logger,
	}
}

func (auth *Authorization) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")

		if len(authHeaderParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if err := auth.authClient.ValidateToken(c.Request.Context(), authHeaderParts[1]); err != nil {
			auth.logger.Error(
				"Invalid token",
				zap.String("token", authHeaderParts[1]),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.GetHeader("User-Agent")),
				zap.Error(err),
			)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
