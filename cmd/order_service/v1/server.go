package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	"github.com/KonnorFrik/BinaryTentacles/pkg/interceptor"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Option - option for customize server at creation.
type Option func(*server) error

type server struct {
	pb.UnimplementedOrderServiceServer
	logger *slog.Logger
	tracer *tracesdk.TracerProvider
}

// NewServer - create a new server with options.
func NewServer(opts ...Option) (*server, error) {
	var srv server

	for _, opt := range opts {
		if err := opt(&srv); err != nil {
			return nil, err
		}
	}

	return &srv, nil
}

// Create - create a new order.
func (s *server) Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*pb.CreateResponse,
	error,
) {
	const method = "Create"
	defer s.startTraceMetdod(ctx, method)()
	order, err := usecase.Create(ctx, req)

	if err != nil {
		return nil, s.wrapError(err, method)
	}

	var response pb.CreateResponse
	order.ToGrpcCreateResponse(&response)
	return &response, status.Error(codes.OK, "ok")
}

// OrderStatus - get a order status.
func (s *server) OrderStatus(
	ctx context.Context,
	req *pb.OrderStatusRequest,
) (
	*pb.OrderStatusResponse,
	error,
) {
	const method = "OrderStatus"
	defer s.startTraceMetdod(ctx, method)()
	order, err := usecase.OrderStatus(ctx, req)

	if err != nil {
		return nil, s.wrapError(err, method)
	}

	var response pb.OrderStatusResponse
	response.Status = order.GetStatus()
	return &response, status.Error(codes.OK, "ok")
}

// OrderUpdates - get order's status update in realtime.
func (s *server) OrderUpdates(
	req *pb.OrderUpdatesRequest,
	stream grpc.ServerStreamingServer[pb.OrderUpdatesResponse],
) error {
	const method = "OrderUpdates"
	defer s.startTraceMetdod(stream.Context(), method)()
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

// startTraceMetdod - start tracing.
// Returns function for end tracing.
func (s *server) startTraceMetdod(ctx context.Context, method string) func() {
	var span trace.Span

	if s.tracer != nil {
		ctx, span = s.tracer.Tracer("Server").Start(ctx, method)
		requestID, ok := ContextValueAs[string](ctx, interceptor.RequestIDHeader)

		if !ok {
			requestID = "-"

			if s.logger != nil {
				s.logger.LogAttrs(
					ctx,
					slog.LevelWarn,
					"order_service/"+method,
					slog.String(interceptor.RequestIDHeader, "not found"),
				)
			}
		}

		span.SetAttributes(
			attribute.String(interceptor.RequestIDHeader, requestID),
		)
		defer span.End()
	}

	return func() { span.End() }
}

// wrapError - log error and call wrapError function.
func (s *server) wrapError(err error, method string) error {
	if err == nil {
		return nil
	}

	if s.logger != nil {
		s.logger.LogAttrs(
			nil,
			slog.LevelError,
			"server/"+method,
			slog.String("error", err.Error()),
		)
	}

	return wrapError(err)
}

// WithSlog - create a server with slog.
func WithSlog(l *slog.Logger) Option {
	return func(s *server) error {
		if l == nil {
			return errors.New("'WithSlog: logger can't be nil")
		}

		s.logger = l
		return nil
	}
}

func WithOtelTracerProvider(tr *tracesdk.TracerProvider) Option {
	return func(s *server) error {
		if tr == nil {
			return errors.New("'WithSlog: tracer can't be nil")
		}

		s.tracer = tr
		return nil
	}
}
