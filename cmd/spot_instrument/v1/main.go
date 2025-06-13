/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"context"
	"net"

	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
)

type server struct {
	pb.UnimplementedSpotInstrumentServiceServer
}

var (
	logger = logging.Default()
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
			logger.UnaryServerInterceptor,
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

// WrapError - wrap usecase error into gRPC error with codes
func WrapError(err error) error {
	var code = codes.Internal
	var msg string

	switch {
	// case errors.Is(err, usecase.ErrDoesNotExist):
	// 	code = codes.NotFound
	// case errors.Is(err, usecase.ErrAlreadyExist):
	// 	code = codes.AlreadyExists
	// case errors.Is(err, usecase.ErrInvalidData):
	// 	code = codes.InvalidArgument
	// case errors.Is(err, usecase.ErrDbNoAccess):
	// 	// default = Internal
	// case errors.Is(err, usecase.ErrUnknown):
	// default = Internal
	}

	return status.Error(code, msg)
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

	if len(markets) == 0 {
		return nil, status.Error(codes.Aborted, "Something went wrong")
	}

	var resp pb.ViewMarketsResponse
	markets[0].ToGrpcViewMarketResponse(markets, &resp)
	return &resp, nil
}
