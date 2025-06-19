package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"

	redCache "github.com/KonnorFrik/BinaryTentacles/pkg/cache/redis"
	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"
)

var (
	marketCache *redCache.Cache
	logger      = logging.Default()
)

// init a redis connection for store markets.
func init() {
	ctx := context.Background()
	config, err := redCache.NewConfig(
		redCache.WithDB(1),
	)

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[SpotInstrument/usecase/init redis]",
			slog.String("Read config", err.Error()),
		)
		return
	}

	marketCache, err = redCache.New(ctx, config, redCache.WithSlog(logger.Logger))

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[SpotInstrument/usecase/init redis]",
			slog.String("Connection", err.Error()),
		)
		return
	}

	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"[SpotInstrument/usecase/init redis]",
		slog.String("Connection", "Successfull"),
	)

	fill()
}

// ViewMarkets - return available markets logic.
func ViewMarkets(
	ctx context.Context,
	req *pb.ViewMarketsRequest,
) (
	[]*market.Market,
	error,
) {
	all, err := marketCache.Values(ctx)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if len(all) == 0 {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[SpotInstrument/ViewMarkets/GetJSONMarkets",
			slog.String("error", "Got 0 markets from redis"),
		)
		return nil, ErrNoMarkets
	}

	var marketsJson = make([]string, 0, len(all))

	for _, v := range all {
		mark, ok := v.(string)

		if ok {
			marketsJson = append(marketsJson, mark)

		} else {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"[SpotInstrument/ViewMarkets/ConvertToStr",
				slog.String("error", fmt.Sprintf("<%T, %[1]+v> Not string", v)),
			)
		}
	}

	var (
		markets = make([]*market.Market, len(marketsJson))
		ind     int
	)

	for _, v := range marketsJson {
		err = json.Unmarshal([]byte(v), &markets[ind])

		if err == nil {
			ind++

		} else {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"[SpotInstrument/ViewMarkets/Unmarshal",
				slog.String("error", err.Error()),
			)
		}
	}

	return markets, nil
}
