package main

import (
	"context"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase"
	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
)

type server struct {
	pb.UnimplementedSpotInstrumentServiceServer
}

// ViewMarkets - return all available markets.
func (s *server) ViewMarkets(
	ctx context.Context,
	req *pb.ViewMarketsRequest,
) (
	*pb.ViewMarketsResponse,
	error,
) {
	const method = "ViewMarkets"
	markets, err := usecase.ViewMarkets(ctx, req)

	if err != nil {
		return nil, s.wrapError(err, method)
	}

	var resp pb.ViewMarketsResponse
	resp.Market = make([]*pb.Market, len(markets))
	market.ToProtobufMany(markets, resp.Market)
	return &resp, nil
}

// wrapError - log error if it not nil and call wrapError.
func (s *server) wrapError(err error, method string) error {
	if err == nil {
		return nil
	}

	logger.LogAttrs(
		nil,
		slog.LevelError,
		"spot_instrument/"+method,
		slog.String("error", err.Error()),
	)

	return wrapError(err)
}
