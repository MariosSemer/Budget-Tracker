package budget

import (
	"fmt"
	"sort"
	"sync"
)

type Store interface {
	Add(e Transaction) Transaction
	List() []Transaction
	Delete(id int64) error
}

type MemoryStore struct {
	mu       sync.RWMutex
	expenses map[int64]Transaction
	nextID   int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{expenses: make(map[int64]Transaction),
		nextID: 1,
	}
}

func (m *MemoryStore) Add(e Transaction) Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	e.ID = m.nextID
	m.nextID++
	m.expenses[e.ID] = e
	return e
}

func (m *MemoryStore) List() []Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]Transaction, 0, len(m.expenses))
	for _, e := range m.expenses {
		list = append(list, e)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (m *MemoryStore) Delete(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.expenses[id]; !ok {
		return fmt.Errorf("Transaction %d not found", id)
	}

	delete(m.expenses, id)

	return nil
}
