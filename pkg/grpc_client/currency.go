package grpc_client

import (
	"fmt"

	"github.com/Sevn9/currency-screener/pkg/currency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCurrencyServiceClient(address string) (currency.CurrencyServiceClient, *grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(address, dialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	client := currency.NewCurrencyServiceClient(conn)
	return client, conn, nil
}
