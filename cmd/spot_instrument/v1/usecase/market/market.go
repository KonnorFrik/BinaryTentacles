package market

import (
	"fmt"
	"sync"
	"time"
)

type Market struct {
	mut       sync.Mutex
	Id        string    `json:"id"`
	Enabled   bool      `json:"enabled"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (m *Market) IsActive() bool {
	var emptyTime time.Time
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.Enabled && m.DeletedAt == emptyTime
}

func (m *Market) String() string {
	m.mut.Lock()
	defer m.mut.Unlock()
	return fmt.Sprintf("Market(ID:%s, Enabled:%t, DeletedAt:%v", m.Id, m.Enabled, m.DeletedAt)
}
