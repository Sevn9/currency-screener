package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/Sevn9/currency-screener/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

func (s *Controller) handleError(c *gin.Context, err error) {

	var nfErr apperrors.NotFoundError
	if errors.As(err, &nfErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": nfErr.Error()})
		return
	}

	switch {
	case errors.Is(err, apperrors.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

	case errors.Is(err, apperrors.ErrUserAlreadyExist):
		c.JSON(http.StatusConflict, gin.H{"error": "User already exist"})

	case errors.Is(err, apperrors.ErrUnexpectedStatusCode):
		log.Printf("unexpected status code error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})

	case errors.Is(err, apperrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

	case errors.Is(err, apperrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})

	case errors.Is(err, apperrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})

	case errors.Is(err, apperrors.ErrAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "Resource already exists"})

	case errors.Is(err, apperrors.ErrTokenGeneration):
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)

	case errors.Is(err, apperrors.ErrTokenNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})

	case errors.Is(err, apperrors.ErrInvalidOrExpiredToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})

	default:
		log.Printf("internal error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
