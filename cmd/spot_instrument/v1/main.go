/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase"
	"github.com/KonnorFrik/BinaryTentacles/pkg/call_chain"
	loggingWrap "github.com/KonnorFrik/BinaryTentacles/pkg/logging"

	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	interceptor "github.com/KonnorFrik/BinaryTentacles/pkg/interceptor"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// TODO: implement gracefull shutdown
var (
	logger = loggingWrap.Default()
)

const (
	laddr = ":9999"
)

func main() {
	osSignalChan := make(chan os.Signal, 3)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	listener, err := net.Listen("tcp", laddr)

	if err != nil {
		logger.Error("[Server/Listen]", "error", err)
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
	pb.RegisterSpotInstrumentServiceServer(grpcServer, userServer)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"Listen at",
		slog.String("local address", laddr),
	)

	gracefullShutdownChain := callChain.New(
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			grpcServer.GracefulStop()
			return nil
		},
		usecase.ShutdownMarketCache,
	)

	var chainGroup sync.WaitGroup
	chainGroup.Add(1)

	go func() {
		defer chainGroup.Done()
		<-osSignalChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		ind, err := gracefullShutdownChain.Call(ctx)

		if err != nil {
			logger.LogAttrs(
				nil,
				slog.LevelError,
				"GracefullShutdownChain",
				slog.Int("stopped at", ind),
				slog.String("with error", err.Error()),
			)
			return
		}

		logger.LogAttrs(
			nil,
			slog.LevelInfo,
			"GracefullShutdownChain",
			slog.String("status", "successfull"),
		)
	}()

	err = grpcServer.Serve(listener)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"Serve",
			slog.String("error", err.Error()),
		)
		return
	}

	chainGroup.Wait()
}
