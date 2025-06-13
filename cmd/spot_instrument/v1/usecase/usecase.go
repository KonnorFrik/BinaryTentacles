package usecase

import (
	"context"

	"github.com/KonnorFrik/BinaryTentacles/cmd/spot_instrument/v1/usecase/market"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
	db "github.com/KonnorFrik/BinaryTentacles/pkg/fake_db"
)

var fakeDb = db.New()

func init() {
	// TODO: fill db with mock data
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
