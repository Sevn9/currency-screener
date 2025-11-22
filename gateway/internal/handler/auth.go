package handler

import (
	"net/http"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Controller) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}

func (s *Controller) Register(c *gin.Context) {
	var req registerRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.authService.Register(dto.RegisterRequest(req))
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Controller) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")

	err := s.authService.Logout(token)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
