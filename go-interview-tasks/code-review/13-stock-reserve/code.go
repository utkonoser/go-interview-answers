//go:build !solution

// Задача на code review (уровень: Lamoda Middle+).
// Сервис резервирования товара на складе — домен из классического тестового Lamoda
// (reserve / release / stock, см. github.com/6jodeci/lmd-tt и аналоги).
// Найдите проблемы перед Black Friday.
package stockreserve

import "errors"

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
	items map[string]*Item
}

func NewStore(skus map[string]int) *Store {
	s := &Store{items: make(map[string]*Item, len(skus))}
	for sku, qty := range skus {
		s.items[sku] = &Item{SKU: sku, Available: qty}
	}
	return s
}

func (s *Store) Reserve(sku string, qty int) error {
	if qty <= 0 {
		return errors.New("qty must be positive")
	}

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

func (s *Store) Release(sku string, qty int) error {
	item, ok := s.items[sku]
	if !ok {
		return ErrUnknownProduct
	}

	item.Reserved -= qty
	return nil
}

func (s *Store) Free(sku string) int {
	item, ok := s.items[sku]
	if !ok {
		return 0
	}
	return item.Available - item.Reserved
}

func (s *Store) Reserved(sku string) int {
	item, ok := s.items[sku]
	if !ok {
		return 0
	}
	return item.Reserved
}
