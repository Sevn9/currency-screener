package handler

import (
	"context"
	"testing"

	mock_handler "github.com/Sevn9/currency-screener/currency/internal/handler/mocks"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrency_GetRates_Success(t *testing.T) {
	//arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_handler.NewMockCurrencyService(ctrl)
	// Создаем мок-объект сервиса
	mockService.EXPECT().GetCurrencyRateFromInterval(
		gomock.Any(), // Любой контекст
		&currency.GetRateRequest{ // Конкретный пустой запрос
			Currency: "",
			DateFrom: nil,
			DateTo:   nil,
		},
	).Return(&currency.GetRateResponse{}, nil) // Вернуть пустой ответ и нет ошибки

	// Создаем "пустой" логгер специально для тестов
	testLogger := zap.NewNop()

	//create server
	server := NewCurrencyServer(mockService, testLogger)

	expected := &currency.GetRateResponse{} // Ожидаем пустой ответ (т.к. мок его вернул)

	ctx := context.Background()

	req := &currency.GetRateRequest{
		Currency: "",
		DateFrom: nil,
		DateTo:   nil,
	}

	// Act
	fact, err := server.GetRates(ctx, req)

	//Assert
	require.NoError(t, err)         // Убедимся, что ошибки нет
	assert.Equal(t, expected, fact) // Убедимся, что вернулся тот же объект, что и от мока
}
