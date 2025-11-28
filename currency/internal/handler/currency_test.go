package handler

import (
	"context"
	"testing"

	"github.com/Sevn9/currency-screener/currency/internal/dto"
	mock_handler "github.com/Sevn9/currency-screener/currency/internal/handler/mocks"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetCurrencyRateFromInterval_Success(t *testing.T) {
	//arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_handler.NewMockCurrencyService(ctrl)

	// Создаем "пустой" логгер специально для тестов
	testLogger := zap.NewNop()

	req := &currency.GetRateRequest{
		Currency: "",
		DateFrom: nil,
		DateTo:   nil,
	}

	expectedDTO := dto.CurrencyRequestDTOFromProtobuf(req, dto.DefaultBaseCurrency)

	serviceResponse := []dto.CurrencyRateResponseDTO{}

	// Создаем мок-объект сервиса
	mockService.EXPECT().GetCurrencyRateFromInterval(
		gomock.Any(), // Любой контекст
		gomock.Eq(expectedDTO),
	).Return(serviceResponse, nil) // Вернуть пустой ответ и нет ошибки

	//create server
	server := NewCurrencyServer(mockService, testLogger)

	expectedProtoResponse := &currency.GetRateResponse{
		Currency: "",
		Rates:    []*currency.RateRecord{},
	}

	ctx := context.Background()

	// Act
	fact, err := server.GetRates(ctx, req)

	//Assert
	require.NoError(t, err)                                        // Убедимся, что ошибки нет
	assert.Equal(t, expectedProtoResponse.Currency, fact.Currency) // Убедимся, что вернулся тот же объект, что и от мока

	assert.Len(t, fact.Rates, 0)
	// Убедимся, что это пустой слайс, а не nil
	assert.NotNil(t, fact.Rates)
}
