package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type MetricsMiddleware struct {
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	appUptime       prometheus.Gauge
}

func NewMetricsMiddleware(
	requestCount *prometheus.CounterVec,
	requestDuration *prometheus.HistogramVec,
	appUptime prometheus.Gauge,
) *MetricsMiddleware {
	return &MetricsMiddleware{
		requestCount:    requestCount,
		requestDuration: requestDuration,
		appUptime:       appUptime,
	}
}

func (m *MetricsMiddleware) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		method := extractMethodName(info.FullMethod)

		// Увеличиваем счетчик запросов
		m.requestCount.WithLabelValues(method).Inc()

		// Замеряем время выполнения
		start := time.Now()

		// Вызываем следующий обработчик
		resp, err := handler(ctx, req)

		// Записываем время выполнения
		duration := time.Since(start).Seconds()
		m.requestDuration.WithLabelValues(method).Observe(duration)

		return resp, err
	}
}

// SetUptime устанавливает значение uptime метрики
func (m *MetricsMiddleware) SetUptime(uptime float64) {
	m.appUptime.Set(uptime)
}

// extractMethodName извлекает имя метода из полного пути gRPC
func extractMethodName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullMethod
}
