package handler

import (
	"net/http"

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

	// todo token logic
	token := "generated_jwt_token_here"

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}

func (s *Controller) Register(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func (s *Controller) Logout(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
