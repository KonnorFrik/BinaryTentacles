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

var (
	fakeDB = db.New()

	spotInstrumentAddr = "spot_instrument:9999"
	spotInstrument     client.SpotInstrumentServiceClient
	logger             = logging.Default()
)

var (
	ErrDoesNotExist  = errors.New("object does not exist")
	ErrInvalidMarket = errors.New("market is unavailable")
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

		return
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
			slog.Uint64("Requested market id", req.GetMarketId()),
			slog.String("Status", "Not found"),
		)
		return nil, ErrInvalidMarket
	}

	var order = new(order.Order)
	order.FromGrpcCreateRequest(req)
	order.Status = pb.OrderStatus_ORDER_STATUS_CREATED
	order.Id = fakeDB.Create(ctx, order)
	return order, nil
}

func OrderStatus(
	ctx context.Context,
	req *pb.OrderStatusRequest,
) (
	*order.Order,
	error,
) {
	order, err := OrderById(ctx, req.GetOrderId())
	return order, err
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
