package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase"
	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
)

type Option func(*server) error

type server struct {
	pb.UnimplementedSpotInstrumentServiceServer
	logger *slog.Logger
}

func New(opts ...Option) (*server, error) {
	var srv server

	for _, opt := range opts {
		if e := opt(&srv); e != nil {
			return nil, e
		}
	}

	return &srv, nil
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

// IsAvailable - check is one market available.
func (s *server) IsAvailable(
	ctx context.Context,
	req *pb.IsAvailableRequest,
) (
	*pb.IsAvailableResponse,
	error,
) {
	const method = "IsAvailable"
	available, err := usecase.IsAvailable(ctx, req)

	if err != nil {
		return nil, s.wrapError(err, method)
	}

	var resp pb.IsAvailableResponse
	resp.IsAvailable = available
	return &resp, nil
}

// wrapError - log error if it not nil and call wrapError.
func (s *server) wrapError(err error, method string) error {
	if err == nil {
		return nil
	}

	if s.logger != nil {
		s.logger.LogAttrs(
			nil,
			slog.LevelError,
			"spot_instrument/"+method,
			slog.String("error", err.Error()),
		)
	}

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
