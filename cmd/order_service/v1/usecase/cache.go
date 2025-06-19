package usecase

import (
	"context"
	"log/slog"

	redCache "github.com/KonnorFrik/BinaryTentacles/pkg/cache/redis"
)

var (
	orderCache *redCache.Cache
)

// init a redis connection for store orders
func init() {
	ctx := context.Background()
	config, err := redCache.NewConfig(
		redCache.WithDB(0),
	)

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[OrderSevice/usecase/init redis]",
			slog.String("Read config", err.Error()),
		)
		return
	}

	orderCache, err = redCache.New(ctx, config, redCache.WithSlog(logger.Logger))

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[OrderSevice/usecase/init redis]",
			slog.String("Connection", err.Error()),
		)
		return
	}

	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"[OrderSevice/usecase/init redis]",
		slog.String("Connection", "Successfull"),
	)
}

func ShutdownOrderCache(ctx context.Context) error {
	err := orderCache.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}
