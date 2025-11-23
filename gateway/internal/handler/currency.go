package handler

import (
	"net/http"
	"time"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type currencyRequest struct {
	Currency string `form:"currency" binding:"required"`
	DateFrom string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo   string `form:"date_to" binding:"required,datetime=2006-01-02"`
}

func (s *Controller) GetCurrencyRates(c *gin.Context) {
	var req currencyRequest
	err := c.BindQuery(&req)
	if err != nil {
		s.logger.Error("Error binding request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for date_from, expected YYYY-MM-DD"})
		return
	}

	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for date_to, expected YYYY-MM-DD"})
		return
	}

	parsedCurrencyRequest := dto.ParsedCurrencyRequest{
		Currency: req.Currency,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	data, err := s.currencyService.GetCurrencyRates(c.Request.Context(), parsedCurrencyRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
