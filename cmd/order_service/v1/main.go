/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	interceptor "github.com/KonnorFrik/BinaryTentacles/pkg/interceptor"
	loggingWrap "github.com/KonnorFrik/BinaryTentacles/pkg/logging"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedOrderServiceServer
}

var (
	logger = loggingWrap.Default()
)

const (
	laddr = ":8888"
)

func main() {
	listener, err := net.Listen("tcp", laddr)

	if err != nil {
		logger.Error("[Server/Listen]", "error", err)
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[Server/Listen]",
			slog.String("error", err.Error()),
		)
		return
	}

	userServer := &server{}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			logging.UnaryServerInterceptor(
				InterceptorLogger(logger.Logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
				logging.WithCodes(ErrorToCode),
			),
			interceptor.UnaryServerXRequestId,
			// From doc - "use those as "last" interceptor, so panic does not skip other interceptors"
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(RecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			logging.StreamServerInterceptor(
				InterceptorLogger(logger.Logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
				logging.WithCodes(ErrorToCode),
			),
			interceptor.StreamServerXRequestId,
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(RecoveryHandler)),
		),
	)
	pb.RegisterOrderServiceServer(grpcServer, userServer)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[Server/Listen]",
		slog.String("Local address", laddr),
	)

	err = grpcServer.Serve(listener)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[Server/Serve]",
			slog.String("error", err.Error()),
		)
		return
	}
}

func (s *server) Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*pb.CreateResponse,
	error,
) {
	order, err := usecase.Create(ctx, req)

	if err != nil {
		return nil, WrapError(err)
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
	order, err := usecase.OrderStatus(ctx, req)

	if err != nil {
		return nil, WrapError(err)
	}

	var response pb.OrderStatusResponse
	response.Status = order.GetStatus()
	return &response, status.Error(codes.OK, "ok")
}

func (s *server) OrderUpdates(
	req *pb.OrderUpdatesRequest,
	stream grpc.ServerStreamingServer[pb.OrderUpdatesResponse],
) error {

	order, err := usecase.OrderById(stream.Context(), req.GetOrderId())

	if err != nil {
		return WrapError(err)
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
