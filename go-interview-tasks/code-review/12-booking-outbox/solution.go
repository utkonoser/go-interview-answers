//go:build solution

package bookingoutbox

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var ErrRoomUnavailable = errors.New("room unavailable")

type Booking struct {
	ID     string
	UserID string
	RoomID int
}

type outboxRow struct {
	topic   string
	payload []byte
}

type Store struct {
	mu       sync.Mutex
	bookings []Booking
	rooms    map[int]bool
	outbox   []outboxRow
}

func NewStore(rooms map[int]bool) *Store {
	return &Store{
		rooms:    rooms,
		bookings: make([]Booking, 0),
		outbox:   make([]outboxRow, 0),
	}
}

// fix: Save и outbox должны быть атомарны (одна транзакция / один Lock).
func (s *Store) CreateBooking(b Booking) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.rooms[b.RoomID] {
		return ErrRoomUnavailable
	}

	payload, err := json.Marshal(b)
	if err != nil {
		return err
	}

	s.bookings = append(s.bookings, b)
	s.outbox = append(s.outbox, outboxRow{topic: "booking.created", payload: payload})
	return nil
}

func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.bookings)
}

func (s *Store) OutboxLen() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.outbox)
}

func (s *Store) DrainOutbox() [][]byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.outbox) == 0 {
		return nil
	}

	rows := make([][]byte, len(s.outbox))
	for i, row := range s.outbox {
		rows[i] = append([]byte(nil), row.payload...)
	}
	s.outbox = s.outbox[:0]
	return rows
}

type Publisher struct {
	mu     sync.Mutex
	events [][]byte
}

func NewPublisher() *Publisher {
	return &Publisher{events: make([][]byte, 0)}
}

func (p *Publisher) Publish(_ context.Context, payload []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, append([]byte(nil), payload...))
	return nil
}

func (p *Publisher) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.events)
}

type Service struct {
	store *Store
	pub   *Publisher
}

func NewService(store *Store, pub *Publisher) *Service {
	return &Service{store: store, pub: pub}
}

// fix: publish до commit в БД — событие уйдёт, а бронь может не сохраниться.
// fix: transactional outbox — relay читает outbox и шлёт в Kafka отдельным воркером.
func (s *Service) CreateBooking(ctx context.Context, b Booking) error {
	if err := s.store.CreateBooking(b); err != nil {
		return err
	}
	return s.relayOutbox(ctx)
}

func (s *Service) relayOutbox(ctx context.Context) error {
	for _, payload := range s.store.DrainOutbox() {
		if err := s.pub.Publish(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}

// RelayPending повторно отправляет outbox (идемпотентный consumer на стороне подписчиков).
func (s *Service) RelayPending(ctx context.Context) error {
	return s.relayOutbox(ctx)
}
