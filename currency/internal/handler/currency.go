package handler

import (
	"context"
	"fmt"

	"github.com/Sevn9/currency-screener/pkg/currency"
)

// grpc proto server function realization
func (c *CurrencyServer) SayHello(ctx context.Context, request *currency.HelloRequest) (*currency.HelloReply, error) {
	message := fmt.Sprintf("Привет, %s!", request.Name)
	return &currency.HelloReply{Message: message}, nil
}
