//go:build solution

package adclick

import (
	"context"
	"errors"
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

// fix: Save без Lock — data race с FindRecent/Purge.
func (r *Repository) Save(c *Click) error {
	r.mu.Lock()
	defer r.mu.Unlock()
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

// fix: TOCTOU в Track — атомарная операция под одним Lock.
func (r *Repository) SaveIfAbsent(userID, campaignID, trackingID string, within time.Duration) (*Click, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-within)
	for i := len(r.clicks) - 1; i >= 0; i-- {
		c := r.clicks[i]
		if c.UserID == userID && c.CampaignID == campaignID && c.CreatedAt.After(cutoff) {
			return c, false, nil
		}
	}

	click := &Click{
		UserID:     userID,
		CampaignID: campaignID,
		TrackingID: trackingID,
		CreatedAt:  time.Now(),
	}
	r.clicks = append(r.clicks, click)
	return click, true, nil
}

func (r *Repository) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clicks)
}

type Publisher struct {
	topic string
	out   chan *Click
	done  chan struct{}
}

func NewPublisher(topic string, buf int) *Publisher {
	p := &Publisher{
		topic: topic,
		out:   make(chan *Click, buf),
		done:  make(chan struct{}),
	}
	go p.run()
	return p
}

func (p *Publisher) run() {
	defer close(p.done)
	for range p.out {
		// ponytail: в проде — отправка в Kafka; здесь no-op consumer
	}
}

// fix: нет consumer → буфер заполняется, default молча дропает клики.
func (p *Publisher) Publish(c *Click) error {
	select {
	case p.out <- c:
		return nil
	default:
		return errors.New("publisher buffer full")
	}
}

func (p *Publisher) Close(ctx context.Context) error {
	close(p.out)
	select {
	case <-p.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Publisher) Topic() string { return p.topic }

type Service struct {
	repo      *Repository
	publisher *Publisher
	ttl       time.Duration
}

// fix: Publisher по значению — антипаттерн для struct с channel/mutex.
func NewService(repo *Repository, publisher *Publisher, ttl time.Duration) *Service {
	return &Service{repo: repo, publisher: publisher, ttl: ttl}
}

func (s *Service) Track(userID, campaignID, trackingID string) (bool, error) {
	click, created, err := s.repo.SaveIfAbsent(userID, campaignID, trackingID, s.ttl)
	if err != nil || !created {
		return created, err
	}

	// fix: лишняя горутина для неблокирующего Publish.
	if err := s.publisher.Publish(click); err != nil {
		return true, err
	}
	return true, nil
}

func (r *Repository) PurgeOlderThan(maxAge time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	keep := r.clicks[:0]
	for _, c := range r.clicks {
		if c.CreatedAt.After(cutoff) {
			keep = append(keep, c)
		}
	}
	r.clicks = keep
}

func (r *Repository) StartPurger(ctx context.Context, interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.PurgeOlderThan(maxAge)
			}
		}
	}()
}
