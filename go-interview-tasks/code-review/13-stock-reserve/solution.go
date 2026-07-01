//go:build solution

package stockreserve

import (
	"errors"
	"sync"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrUnknownProduct    = errors.New("unknown product")
)

type Item struct {
	SKU       string
	Available int
	Reserved  int
}

type Store struct {
	mu    sync.Mutex
	items map[string]*Item
}

func NewStore(skus map[string]int) *Store {
	s := &Store{items: make(map[string]*Item, len(skus))}
	for sku, qty := range skus {
		s.items[sku] = &Item{SKU: sku, Available: qty}
	}
	return s
}

// fix: check-then-act без Lock → oversell при конкурентных Reserve.
func (s *Store) Reserve(sku string, qty int) error {
	if qty <= 0 {
		return errors.New("qty must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[sku]
	if !ok {
		return ErrUnknownProduct
	}

	free := item.Available - item.Reserved
	if free < qty {
		return ErrInsufficientStock
	}

	item.Reserved += qty
	return nil
}

// fix: Release не проверяет qty > Reserved → Reserved уходит в минус.
func (s *Store) Release(sku string, qty int) error {
	if qty <= 0 {
		return errors.New("qty must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[sku]
	if !ok {
		return ErrUnknownProduct
	}
	if item.Reserved < qty {
		return ErrInsufficientStock
	}

	item.Reserved -= qty
	return nil
}

func (s *Store) Free(sku string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[sku]
	if !ok {
		return 0
	}
	return item.Available - item.Reserved
}

func (s *Store) Reserved(sku string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[sku]
	if !ok {
		return 0
	}
	return item.Reserved
}
