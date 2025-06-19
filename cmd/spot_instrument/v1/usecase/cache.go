package usecase

import (
	"context"
	"log/slog"

	redCache "github.com/KonnorFrik/BinaryTentacles/pkg/cache/redis"
)

var (
	marketCache *redCache.Cache
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

func ShutdownMarketCache(ctx context.Context) error {
	return marketCache.Close(ctx)
}
