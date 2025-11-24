package storage

import (
	"fmt"
	"mini-quicko/pkg/models"
	"sync"
)

const MaxHistoryEntries = 100

type Storage interface {
	AddHistoryEntry(productID string, entry models.PriceHistoryEntry)
	GetHistory(productID string) (*models.PriceHistory, error)
}

type MemoryStorage struct {
	mu      sync.RWMutex
	history map[string][]models.PriceHistoryEntry
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		history: make(map[string][]models.PriceHistoryEntry),
	}
}

func (m *MemoryStorage) AddHistoryEntry(productID string, entry models.PriceHistoryEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.history[productID]; !exists {
		m.history[productID] = make([]models.PriceHistoryEntry, 0)
	}

	m.history[productID] = append(m.history[productID], entry)

	if len(m.history[productID]) > MaxHistoryEntries {
		m.history[productID] = m.history[productID][1:]
	}
}

func (m *MemoryStorage) GetHistory(productID string) (*models.PriceHistory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries, exists := m.history[productID]
	if !exists {
		return nil, fmt.Errorf("no history found for product: %s", productID)
	}

	return &models.PriceHistory{
		ProductID: productID,
		History:   entries,
	}, nil
}
