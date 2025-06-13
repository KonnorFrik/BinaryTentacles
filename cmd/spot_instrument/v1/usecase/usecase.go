package usecase

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	db "github.com/KonnorFrik/BinaryTentacles/pkg/fake_db"
	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"
)

var (
	fakeDb = db.New()
	logger = logging.Default()
)

func init() {
	var (
		mark *market.Market
		ctx  = context.Background()
	)
	mark = new(market.Market)
	mark.Enabled = true
	mark.Id = fakeDb.Create(ctx, mark)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[SpotInstrument/init]",
		slog.String("Created valid market with id", strconv.FormatUint(mark.Id, 10)),
	)

	mark = new(market.Market)
	mark.Enabled = true
	mark.DeletedAt = time.Now()
	mark.Id = fakeDb.Create(ctx, mark)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[SpotInstrument/init]",
		slog.String("Created invalid market with id", strconv.FormatUint(mark.Id, 10)),
	)

	mark = new(market.Market)
	mark.DeletedAt = time.Now()
	mark.Id = fakeDb.Create(ctx, mark)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[SpotInstrument/init]",
		slog.String("Created invalid market with id", strconv.FormatUint(mark.Id, 10)),
	)

	mark = new(market.Market)
	mark.Id = fakeDb.Create(ctx, mark)
	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[SpotInstrument/init]",
		slog.String("Created invalid market with id", strconv.FormatUint(mark.Id, 10)),
	)

	logger.LogAttrs(
		nil,
		slog.LevelInfo,
		"[SpotInstrument/init]",
		slog.String("fakeDB status", "filled with mock objects"),
	)
}

func ViewMarkets(
	ctx context.Context,
	req *pb.ViewMarketsRequest,
) (
	[]*market.Market,
	error,
) {
	var markets []*market.Market
	allMarkets := db.All[*market.Market](ctx, fakeDb)

	for _, m := range allMarkets {
		if m.IsActive() {
			markets = append(markets, m)
		}
	}

	return markets, nil
}
