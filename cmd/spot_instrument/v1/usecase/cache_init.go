package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	"github.com/google/uuid"
)

func newMarket(enabled bool, delAt time.Time) *market.Market {
	var mark = market.Market{
		Enabled:   enabled,
		DeletedAt: delAt,
		Id:        uuid.NewString(),
	}
	return &mark
}

func fill() {
	var (
		mark *market.Market
		ctx  = context.Background()
	)
	mark = &market.Market{Enabled: true, DeletedAt: time.Time{}, Id: "5d6f8857-fafe-432c-8380-2b340ec03bb7"}
	markBytes, _ := json.Marshal(mark)
	err := marketCache.Set(ctx, mark.Id, string(markBytes), time.Hour)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[SpotInstrument/init]",
			slog.String("Create market error", err.Error()),
		)

	} else {
		logger.LogAttrs(
			nil,
			slog.LevelInfo,
			"[SpotInstrument/init]",
			slog.String("Created valid market with id", mark.Id),
		)
	}

	mark = newMarket(true, time.Now())
	markBytes, _ = json.Marshal(mark)
	err = marketCache.Set(ctx, mark.Id, string(markBytes), time.Hour)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[SpotInstrument/init]",
			slog.String("Create market error", err.Error()),
		)

	} else {
		logger.LogAttrs(
			nil,
			slog.LevelInfo,
			"[SpotInstrument/init]",
			slog.String("Created invalid market with id", mark.Id),
		)
	}

	mark = newMarket(false, time.Now())
	markBytes, _ = json.Marshal(mark)
	err = marketCache.Set(ctx, mark.Id, string(markBytes), time.Hour)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[SpotInstrument/init]",
			slog.String("Create market error", err.Error()),
		)

	} else {
		logger.LogAttrs(
			nil,
			slog.LevelInfo,
			"[SpotInstrument/init]",
			slog.String("Created invalid market with id", mark.Id),
		)
	}

	mark = newMarket(false, time.Time{})
	markBytes, _ = json.Marshal(mark)
	err = marketCache.Set(ctx, mark.Id, string(markBytes), time.Hour)

	if err != nil {
		logger.LogAttrs(
			nil,
			slog.LevelError,
			"[SpotInstrument/init]",
			slog.String("Create market error", err.Error()),
		)

	} else {
		logger.LogAttrs(
			nil,
			slog.LevelInfo,
			"[SpotInstrument/init]",
			slog.String("Created invalid market with id", mark.Id),
		)
	}
}
