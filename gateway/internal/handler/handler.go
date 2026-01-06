package handler

import (
	"net/http"

	"github.com/Sevn9/currency-screener/gateway/internal/middleware"
	"github.com/Sevn9/currency-screener/gateway/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	authService     services.AuthService
	currencyService services.CurrencyService
	router          *gin.Engine
	logger          *zap.Logger
}

func RegisterRoutes(
	router *gin.Engine,
	logger *zap.Logger,
	authMiddleware middleware.Authorization,
	authService services.AuthService,
	currencyService services.CurrencyService) {

	cntrl := Controller{
		router:          router,
		logger:          logger,
		authService:     authService,
		currencyService: currencyService,
	}

	// Global health-check
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	//public routes
	public := router.Group("/api/v1")
	{
		public.POST("/login", cntrl.Login)
		public.POST("/register", cntrl.Register)
	}

	//private routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.Authorize())
	{
		protected.POST("/logout", cntrl.Logout)

		// GET/api/v1/rate
		protected.GET("/rate", cntrl.GetCurrencyRates)
	}
}
