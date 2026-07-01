//go:build solution

package wallet

import (
	"errors"
	"sync"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

// fix: map без синхронизации — data race при конкурентных Credit/Transfer.
type Store struct {
	mu       sync.Mutex
	balances map[int64]int64
}

func NewStore() *Store {
	return &Store{balances: make(map[int64]int64)}
}

func (s *Store) Balance(userID int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.balances[userID]
}

func (s *Store) Credit(userID, amount int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.balances[userID] += amount
}

func (s *Store) Debit(userID, amount int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.debitLocked(userID, amount)
}

// fix: Transfer read-check-write не атомарен → отрицательный баланс при гонке.
func (s *Store) Transfer(from, to, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.debitLocked(from, amount); err != nil {
		return err
	}
	s.balances[to] += amount
	return nil
}

func (s *Store) debitLocked(userID, amount int64) error {
	if s.balances[userID] < amount {
		return ErrInsufficientFunds
	}
	s.balances[userID] -= amount
	return nil
}
