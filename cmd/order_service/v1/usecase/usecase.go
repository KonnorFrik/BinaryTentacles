package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase/order"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	client "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	redCache "github.com/KonnorFrik/BinaryTentacles/pkg/cache/redis"
	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	spotInstrumentAddr = "spot_instrument:9999"
)

var (
	spotInstrument client.SpotInstrumentServiceClient
	logger         = logging.Default()
)

var (
	// ErrDoesNotExist - data is not exist
	ErrDoesNotExist = errors.New("object does not exist")
	// ErrMarketUnavailable - market is unavailable for any reason
	ErrMarketUnavailable = errors.New("market is unavailable")
	// ErrUnknown - any undocumented error
	ErrUnknown = errors.New("unknown")
	// ErrInternal - indicate errors for any reason in OrderSevice/usecase logic
	ErrInternal = errors.New("internal")
)

// init a grpc connection with spot instrument service
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

// Create - create a order logic.
func Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*order.Order,
	error,
) {
	if e := uuid.Validate(req.MarketId); e != nil {
		return nil, fmt.Errorf("%w: requested market id is invalid", ErrMarketUnavailable)
	}

	clientReq := client.ViewMarketsRequest{
		UserRole: client.UserRole_USER_ROLE_CUSTOMER,
	}
	// TODO: cache this
	marketsResponse, err := spotInstrument.ViewMarkets(ctx, &clientReq)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[OrderService/Create]",
			slog.String("Call", "SpotInstrumentService.ViewMarkets"),
			slog.String("Error", err.Error()),
		)
		return nil, fmt.Errorf("%w: SpotInstrumentService: %w", ErrInternal, err)
	}

	var (
		marketIdCount   int
		marketIdRequest = req.GetMarketId()
	)

	// TODO: implement spot_instrument.IsAvailable(market_uuid) bool instead this.
	for _, m := range marketsResponse.Market {
		if marketIdRequest == m.GetId() {
			marketIdCount++
		}
	}

	if marketIdCount == 0 {
		return nil, fmt.Errorf("%w: found no markets", ErrMarketUnavailable)
	}

	var order = new(order.Order)
	order.FromGrpcCreateRequest(req)
	order.Status = pb.OrderStatus_ORDER_STATUS_CREATED
	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	orderJsonBytes, err := json.Marshal(order)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	order.Id = id.String()
	err = orderCache.Set(ctx, order.Id, string(orderJsonBytes), time.Hour)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return order, nil
}

// OrderStatus - return a order status logic.
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

// OrderById - get order from db by it id.
func OrderById(
	ctx context.Context,
	id string,
) (
	*order.Order,
	error,
) {
	orderJson, err := orderCache.Get(ctx, id)

	if err != nil {
		if err == redCache.ErrNil {
			return nil, ErrDoesNotExist
		}

		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	var ord order.Order
	err = json.Unmarshal([]byte(orderJson), &ord)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return &ord, nil
}
