package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Option func(*server) error

type server struct {
	pb.UnimplementedOrderServiceServer
	logger *slog.Logger
}

func NewServer(opts ...Option) (*server, error) {
	var srv server

	for _, opt := range opts {
		if err := opt(&srv); err != nil {
			return nil, err
		}
	}

	return &srv, nil
}

func (s *server) Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*pb.CreateResponse,
	error,
) {
	const method = "Create"
	order, err := usecase.Create(ctx, req)

	if err != nil {
		return nil, s.wrapError(err, method)
	}

	var response pb.CreateResponse
	order.ToGrpcCreateResponse(&response)
	return &response, status.Error(codes.OK, "ok")
}

func (s *server) OrderStatus(
	ctx context.Context,
	req *pb.OrderStatusRequest,
) (
	*pb.OrderStatusResponse,
	error,
) {
	const method = "OrderStatus"
	order, err := usecase.OrderStatus(ctx, req)

	if err != nil {
		return nil, s.wrapError(err, method)
	}

	var response pb.OrderStatusResponse
	response.Status = order.GetStatus()
	return &response, status.Error(codes.OK, "ok")
}

func (s *server) OrderUpdates(
	req *pb.OrderUpdatesRequest,
	stream grpc.ServerStreamingServer[pb.OrderUpdatesResponse],
) error {
	const method = "OrderUpdates"
	order, err := usecase.OrderById(stream.Context(), req.GetOrderId())

	if err != nil {
		return s.wrapError(err, method)
	}

	statuses := order.UpdateStatus(stream.Context())

	for {
		var (
			resp     = new(pb.OrderUpdatesResponse)
			isClosed bool
		)

		select {
		case <-stream.Context().Done():
			return nil
		case resp.Status, isClosed = <-statuses:
		}

		if isClosed {
			return nil
		}

		if e := stream.Send(resp); e != nil {
			// TODO: catch a closed by a client connection
			return e
		}
	}
}

func (s *server) wrapError(err error, method string) error {
	if err == nil {
		return nil
	}

	s.logger.LogAttrs(
		nil,
		slog.LevelError,
		"server/"+method,
		slog.String("error", err.Error()),
	)

	return wrapError(err)
}

func WithSlog(l *slog.Logger) Option {
	return func(s *server) error {
		if l == nil {
			return errors.New("'WithSlog: logger can't be nil")
		}

		s.logger = l
		return nil
	}
}
