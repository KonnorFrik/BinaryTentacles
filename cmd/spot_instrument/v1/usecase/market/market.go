package market

import (
	"fmt"
	"sync"
	"time"
)

// Market - some market in stock exchange.
type Market struct {
	mut       sync.Mutex
	Id        string    `json:"id"`
	Enabled   bool      `json:"enabled"`
	DeletedAt time.Time `json:"deleted_at"`
}

// IsActive - check is market active.
// Must be enabled and not deleted.
func (m *Market) IsActive() bool {
	var emptyTime time.Time
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.Enabled && m.DeletedAt == emptyTime
}

// String - implement Stringer interface.
func (m *Market) String() string {
	m.mut.Lock()
	defer m.mut.Unlock()
	return fmt.Sprintf("Market(ID:%s, Enabled:%t, DeletedAt:%v", m.Id, m.Enabled, m.DeletedAt)
}
