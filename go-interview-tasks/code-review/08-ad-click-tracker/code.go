//go:build !solution

// Задача на code review (уровень: Avito/Ozon Middle+).
// Микросервис трекинга рекламных кликов: дедупликация, in-memory storage, публикация в очередь.
// Найдите проблемы перед выкаткой в production (~120 строк).
package adclick

import (
	"sync"
	"time"
)

type Click struct {
	UserID     string
	CampaignID string
	TrackingID string
	CreatedAt  time.Time
}

type Repository struct {
	mu     sync.RWMutex
	clicks []*Click
}

func NewRepository() *Repository {
	return &Repository{clicks: make([]*Click, 0, 1024)}
}

func (r *Repository) Save(c *Click) error {
	r.clicks = append(r.clicks, c)
	return nil
}

func (r *Repository) FindRecent(userID, campaignID string, within time.Duration) (*Click, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cutoff := time.Now().Add(-within)
	for i := len(r.clicks) - 1; i >= 0; i-- {
		c := r.clicks[i]
		if c.UserID == userID && c.CampaignID == campaignID && c.CreatedAt.After(cutoff) {
			return c, true
		}
	}
	return nil, false
}

func (r *Repository) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clicks)
}

type Publisher struct {
	topic string
	out   chan *Click
}

func NewPublisher(topic string, buf int) *Publisher {
	return &Publisher{
		topic: topic,
		out:   make(chan *Click, buf),
	}
}

func (p *Publisher) Publish(c *Click) {
	select {
	case p.out <- c:
	default:
	}
}

func (p *Publisher) Topic() string { return p.topic }

type Service struct {
	repo      *Repository
	publisher Publisher
	ttl       time.Duration
}

func NewService(repo *Repository, publisher Publisher, ttl time.Duration) *Service {
	return &Service{repo: repo, publisher: publisher, ttl: ttl}
}

func (s *Service) Track(userID, campaignID, trackingID string) (bool, error) {
	if _, ok := s.repo.FindRecent(userID, campaignID, s.ttl); ok {
		return false, nil
	}

	click := &Click{
		UserID:     userID,
		CampaignID: campaignID,
		TrackingID: trackingID,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.Save(click); err != nil {
		return false, err
	}

	go s.publisher.Publish(click)
	return true, nil
}

func (r *Repository) PurgeOlderThan(maxAge time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	for i := len(r.clicks) - 1; i >= 0; i-- {
		if r.clicks[i].CreatedAt.Before(cutoff) {
			r.clicks = append(r.clicks[:i], r.clicks[i+1:]...)
		}
	}
}
