package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase/order"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	client "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	db "github.com/KonnorFrik/BinaryTentacles/pkg/fake_db"
	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: create market client and call it

var (
	fakeDB = db.New()

	spotInstrumentAddr = "0.0.0.0:9999"
	spotInstrument     client.SpotInstrumentServiceClient
	logger             = logging.Default()
)

var (
	ErrDoesNotExist = errors.New("object does not exist")
)

func init() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(spotInstrumentAddr, opts...)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[OrderService/init]",
			slog.String("Connection to", "SpotInstrumentService"),
			slog.String("Error", err.Error()),
		)
	}

	spotInstrument = client.NewSpotInstrumentServiceClient(conn)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[OrderService/init]",
		slog.String("Connection to", "SpotInstrumentService"),
		slog.String("Status", "Successfull"),
	)
}

func Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*order.Order,
	error,
) {
	clientReq := client.ViewMarketsRequest{
		UserRole: client.UserRole_USER_ROLE_CUSTOMER,
	}
	marketsResponse, err := spotInstrument.ViewMarkets(ctx, &clientReq)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[OrderService/Create]",
			slog.String("Call", "SpotInstrumentService.ViewMarkets"),
			slog.String("Error", err.Error()),
		)
		return nil, err
	}

	var (
		marketIdCount   int
		marketIdRequest = req.GetMarketId()
	)

	for _, m := range marketsResponse.Market {
		if marketIdRequest == m.GetId() {
			marketIdCount++
		}
	}

	if marketIdCount == 0 {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[OrderService/Create]",
			slog.String("Markets status", "No available"),
		)
		return nil, errors.New("No available markets")
	}

	var order = new(order.Order)
	order.FromGrpcCreateRequest(req)
	id := fakeDB.Create(ctx, order)
	order.Id = id
	return order, nil
}

func OrderStatus(
	ctx context.Context,
	req *pb.OrderStatusRequest,
) (
	*order.Order,
	error,
) {
	order, ok := db.As[*order.Order](ctx, fakeDB, req.GetOrderId())

	if !ok {
		return nil, ErrDoesNotExist
	}

	return order, nil
}

func OrderById(
	ctx context.Context,
	id uint64,
) (
	*order.Order,
	error,
) {
	order, ok := db.As[*order.Order](ctx, fakeDB, id)

	if !ok {
		return nil, ErrDoesNotExist
	}

	return order, nil
}
