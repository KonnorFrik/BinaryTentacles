package market

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
)

type Market struct {
	mut       sync.Mutex
	Id        uint64
	Enabled   bool
	DeletedAt time.Time
}

func (m *Market) IsActive() bool {
	var emptyTime time.Time
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.Enabled && m.DeletedAt == emptyTime
}

func (m *Market) ToGrpcViewMarketResponse(
	markets []*Market,
	resp *pb.ViewMarketsResponse,
) *Market {
	resp.Market = resp.Market[:0]

	for _, value := range markets {
		var pbMarket = new(pb.Market)
		pbMarket.Id = value.Id
		resp.Market = append(resp.Market, pbMarket)
	}

	return m
}

func (m *Market) String() string {
	m.mut.Lock()
	defer m.mut.Unlock()
	return fmt.Sprintf("Market(ID:%d, Enabled:%t, DeletedAt:%v", m.Id, m.Enabled, m.DeletedAt)
}
