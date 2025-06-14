/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"context"
	"net"

	loggingWrap "github.com/KonnorFrik/BinaryTentacles/pkg/logging"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase"
	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	interceptor "github.com/KonnorFrik/BinaryTentacles/pkg/interceptor"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSpotInstrumentServiceServer
}

var (
	logger = loggingWrap.Default()
)

const (
	laddr = ":9999"
)

func main() {
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
	logger.Info("Listen at", "local address", laddr)

	err = grpcServer.Serve(listener)

	if err != nil {
		logger.Error("Serve", "error", err)
		return
	}
}

func (s *server) ViewMarkets(
	ctx context.Context,
	req *pb.ViewMarketsRequest,
) (
	*pb.ViewMarketsResponse,
	error,
) {
	markets, err := usecase.ViewMarkets(ctx, req)

	if err != nil {
		return nil, WrapError(err)
	}

	var resp pb.ViewMarketsResponse
	resp.Market = make([]*pb.Market, len(markets))
	market.ToProtobufMany(markets, resp.Market)
	return &resp, nil
}
