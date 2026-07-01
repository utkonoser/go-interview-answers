//go:build !solution

// Задача на code review (уровень: Avito/Tinkoff Middle).
// Микросервис баланса пользователей: зачисление, списание, перевод.
// Паттерн из тестовых Avito (internship_backend_2022, job-backend-trainee-assignment).
package wallet

import "errors"

var ErrInsufficientFunds = errors.New("insufficient funds")

type Store struct {
	balances map[int64]int64
}

func NewStore() *Store {
	return &Store{balances: make(map[int64]int64)}
}

func (s *Store) Balance(userID int64) int64 {
	return s.balances[userID]
}

func (s *Store) Credit(userID, amount int64) {
	s.balances[userID] += amount
}

func (s *Store) Debit(userID, amount int64) error {
	if s.balances[userID] < amount {
		return ErrInsufficientFunds
	}
	s.balances[userID] -= amount
	return nil
}

func (s *Store) Transfer(from, to, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if s.balances[from] < amount {
		return ErrInsufficientFunds
	}
	s.balances[from] -= amount
	s.balances[to] += amount
	return nil
}
