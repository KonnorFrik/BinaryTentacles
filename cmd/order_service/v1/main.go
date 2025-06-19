/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"log/slog"
	"net"

	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	interceptor "github.com/KonnorFrik/BinaryTentacles/pkg/interceptor"
	loggingWrap "github.com/KonnorFrik/BinaryTentacles/pkg/logging"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

var (
// logger = loggingWrap.Default()
)

const (
	laddr = ":8888"
)

func main() {
	logger := loggingWrap.Default()
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

	userServer, err := NewServer(
		WithSlog(logger.Logger),
	)
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
