//go:build !solution

// Задача на code review (уровень: 2ГИС Application Review).
// Сервис бронирования: сохранить бронь и отправить событие в шину (Kafka).
// На собесе спрашивают про transactional outbox и слои REST-сервиса.
package bookingoutbox

import (
	"context"
	"encoding/json"
	"errors"
)

var ErrRoomUnavailable = errors.New("room unavailable")

type Booking struct {
	ID     string
	UserID string
	RoomID int
}

type Store struct {
	bookings []Booking
	rooms    map[int]bool
}

func NewStore(rooms map[int]bool) *Store {
	return &Store{rooms: rooms, bookings: make([]Booking, 0)}
}

func (s *Store) Save(b Booking) error {
	if !s.rooms[b.RoomID] {
		return ErrRoomUnavailable
	}
	s.bookings = append(s.bookings, b)
	return nil
}

func (s *Store) Count() int {
	return len(s.bookings)
}

type Publisher struct {
	events [][]byte
}

func NewPublisher() *Publisher {
	return &Publisher{events: make([][]byte, 0)}
}

func (p *Publisher) Publish(_ context.Context, payload []byte) error {
	p.events = append(p.events, payload)
	return nil
}

func (p *Publisher) Len() int {
	return len(p.events)
}

type Service struct {
	store *Store
	pub   *Publisher
}

func NewService(store *Store, pub *Publisher) *Service {
	return &Service{store: store, pub: pub}
}

func (s *Service) CreateBooking(ctx context.Context, b Booking) error {
	payload, err := json.Marshal(b)
	if err != nil {
		return err
	}

	if err := s.pub.Publish(ctx, payload); err != nil {
		return err
	}

	return s.store.Save(b)
}
