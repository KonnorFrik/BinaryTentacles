package spot_instrument_v1_test

import (
	"context"
	"testing"

	client "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	marketIdValid      uint64 = 0
	marketIdInvalidMax uint64 = 4

	orderServiceAddr = "0.0.0.0:9999"
)

var (
	service client.SpotInstrumentServiceClient
	baseCtx = context.Background()
	orderId uint64
)

func init() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(orderServiceAddr, opts...)

	if err != nil {
		panic(err)
	}

	service = client.NewSpotInstrumentServiceClient(conn)
}

func TestViewMarkets(t *testing.T) {
	req := client.ViewMarketsRequest{
		UserRole: client.UserRole_USER_ROLE_CUSTOMER,
	}
	resp, err := service.ViewMarkets(baseCtx, &req)

	if err != nil {
		t.Fatalf("Got = %q\n", err)
	}

	if len(resp.Market) == 0 {
		t.Fatalf("Got no markets")
	}
}
