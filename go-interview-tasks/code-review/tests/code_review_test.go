//go:build solution

package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	racecounter "interviews/go-interview-tasks/code-review/01-race-counter"
	goroutineleak "interviews/go-interview-tasks/code-review/02-goroutine-leak"
	waitgroupmisuse "interviews/go-interview-tasks/code-review/03-waitgroup-misuse"
	httpclient "interviews/go-interview-tasks/code-review/04-http-client"
	concurrentmap "interviews/go-interview-tasks/code-review/05-concurrent-map"
	deferinloop "interviews/go-interview-tasks/code-review/06-defer-in-loop"
	errorhandling "interviews/go-interview-tasks/code-review/07-error-handling"
	adclick "interviews/go-interview-tasks/code-review/08-ad-click-tracker"
	wallet "interviews/go-interview-tasks/code-review/09-wallet-transfer"
	parindex "interviews/go-interview-tasks/code-review/10-parallel-search"
	requestmetrics "interviews/go-interview-tasks/code-review/11-request-metrics"
	bookingoutbox "interviews/go-interview-tasks/code-review/12-booking-outbox"
	stockreserve "interviews/go-interview-tasks/code-review/13-stock-reserve"
	stockcache "interviews/go-interview-tasks/code-review/14-stock-cache"
	orderpricing "interviews/go-interview-tasks/code-review/15-order-pricing"
)

func TestRaceCounter(t *testing.T) {
	t.Parallel()

	const goroutines = 100
	const perG = 100

	c := &racecounter.Counter{}
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range perG {
				c.Inc()
			}
		}()
	}
	wg.Wait()

	if got := c.Value(); got != goroutines*perG {
		t.Fatalf("counter = %d, want %d", got, goroutines*perG)
	}
}

func TestGoroutineLeak(t *testing.T) {
	t.Parallel()

	got := goroutineleak.DoubleAll([]int{1, 2, 3, 4, 5})
	want := []int{2, 4, 6, 8, 10}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestWaitGroupMisuse(t *testing.T) {
	t.Parallel()

	const n = 100
	var mu sync.Mutex
	done := 0

	waitgroupmisuse.Process(n, func(_ int) {
		mu.Lock()
		done++
		mu.Unlock()
	})

	if done != n {
		t.Fatalf("processed %d, want %d", done, n)
	}
}

func TestHTTPClient(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)

	body, err := httpclient.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "ok" {
		t.Fatalf("body = %q, want ok", body)
	}
}

func TestConcurrentMap(t *testing.T) {
	t.Parallel()

	c := concurrentmap.New()
	const n = 50
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for range n {
		go func() {
			defer wg.Done()
			c.Set("key", "value")
		}()
		go func() {
			defer wg.Done()
			_, _ = c.Get("key")
		}()
	}
	wg.Wait()
}

func TestDeferInLoop(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	paths := make([]string, 3)
	for i := range paths {
		paths[i] = filepath.Join(dir, string(rune('a'+i)))
		if err := os.WriteFile(paths[i], []byte{byte('A' + i)}, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := deferinloop.FirstBytes(paths)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{'A', 'B', 'C'}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "out.dat")
	data := []byte("hello")

	if err := errorhandling.Save(path, data); err != nil {
		t.Fatal(err)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(read) != string(data) {
		t.Fatalf("got %q, want %q", read, data)
	}
}

func TestAdClickDedup(t *testing.T) {
	t.Parallel()

	repo := adclick.NewRepository()
	pub := adclick.NewPublisher("clicks", 64)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = pub.Close(ctx)
	})

	svc := adclick.NewService(repo, pub, 5*time.Second)

	var wg sync.WaitGroup
	const n = 40
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			_, _ = svc.Track("user-1", "camp-1", "track-1")
		}()
	}
	wg.Wait()

	if repo.Len() != 1 {
		t.Fatalf("dedup failed: %d clicks stored, want 1", repo.Len())
	}
}

func TestWalletTransferConservation(t *testing.T) {
	t.Parallel()

	s := wallet.NewStore()
	s.Credit(1, 100)

	var wg sync.WaitGroup
	const n = 20
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			_ = s.Transfer(1, 2, 10)
		}()
	}
	wg.Wait()

	if s.Balance(1) < 0 {
		t.Fatalf("negative balance on user 1: %d", s.Balance(1))
	}
	if s.Balance(1)+s.Balance(2) != 100 {
		t.Fatalf("money lost: user1=%d user2=%d", s.Balance(1), s.Balance(2))
	}
}

func TestParallelSearch(t *testing.T) {
	t.Parallel()

	docs := []string{"alpha", "beta", "alphabet", "gamma"}
	got := parindex.FindIndexes(docs, "alp")
	if len(got) != 2 {
		t.Fatalf("got %v, want 2 matches", got)
	}

	if max := parindex.MaxEven([]int{1, 8, 3, 12, 5}); max != 12 {
		t.Fatalf("MaxEven = %d, want 12", max)
	}
}

func TestRequestMetrics(t *testing.T) {
	t.Parallel()

	m := requestmetrics.New()
	mux := http.NewServeMux()
	requestmetrics.RegisterRoutes(mux, m)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			resp, err := http.Get(srv.URL + "/api/v1/places")
			if err != nil {
				t.Error(err)
				return
			}
			_ = resp.Body.Close()
		}()
	}
	wg.Wait()

	if got := m.Total(); got != n {
		t.Fatalf("total = %v, want %d", got, n)
	}
}

func TestBookingOutboxAtomicity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := bookingoutbox.NewStore(map[int]bool{101: true})
	pub := bookingoutbox.NewPublisher()
	svc := bookingoutbox.NewService(store, pub)

	okBooking := bookingoutbox.Booking{ID: "b1", UserID: "u1", RoomID: 101}
	if err := svc.CreateBooking(ctx, okBooking); err != nil {
		t.Fatal(err)
	}
	if store.Count() != 1 || pub.Len() != 1 {
		t.Fatalf("booking=%d events=%d, want 1/1", store.Count(), pub.Len())
	}

	badBooking := bookingoutbox.Booking{ID: "b2", UserID: "u1", RoomID: 404}
	if err := svc.CreateBooking(ctx, badBooking); !errors.Is(err, bookingoutbox.ErrRoomUnavailable) {
		t.Fatalf("err = %v", err)
	}
	if store.Count() != 1 {
		t.Fatalf("bookings = %d, want 1", store.Count())
	}
	if pub.Len() != 1 {
		t.Fatalf("phantom event published: events = %d, want 1", pub.Len())
	}
}

func TestStockReserveNoOversell(t *testing.T) {
	t.Parallel()

	const sku = "BOOT-42"
	const stock = 10
	const workers = 20
	const each = 1

	store := stockreserve.NewStore(map[string]int{sku: stock})
	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()
			_ = store.Reserve(sku, each)
		}()
	}
	wg.Wait()

	reserved := store.Reserved(sku)
	if reserved > stock {
		t.Fatalf("oversell: reserved %d, stock %d", reserved, stock)
	}
	if reserved+store.Free(sku) != stock {
		t.Fatalf("invariant broken: reserved=%d free=%d stock=%d", reserved, store.Free(sku), stock)
	}
}

func TestStockReleaseBounds(t *testing.T) {
	t.Parallel()

	const sku = "COAT-1"
	store := stockreserve.NewStore(map[string]int{sku: 5})
	if err := store.Reserve(sku, 3); err != nil {
		t.Fatal(err)
	}
	if err := store.Release(sku, 10); err == nil {
		t.Fatal("expected error on over-release")
	}
	if store.Reserved(sku) != 3 {
		t.Fatalf("reserved = %d, want 3", store.Reserved(sku))
	}
}

func TestStockCacheSingleLoad(t *testing.T) {
	t.Parallel()

	const sku = "SKU-1"
	var loads atomic.Int32

	cache := stockcache.New(time.Minute)
	loader := func(string) (int, error) {
		loads.Add(1)
		time.Sleep(20 * time.Millisecond)
		return 100, nil
	}

	const n = 30
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			qty, err := cache.Get(sku, loader)
			if err != nil || qty != 100 {
				t.Errorf("qty=%d err=%v", qty, err)
			}
		}()
	}
	wg.Wait()

	if got := loads.Load(); got != 1 {
		t.Fatalf("loader called %d times, want 1", got)
	}
}

func TestOrderPricingKopecks(t *testing.T) {
	t.Parallel()

	// 1499.99 ₽ × 2 − 10% = 2699.98 ₽ = 269998 коп.
	lines := []orderpricing.Line{
		{PriceKopecks: orderpricing.RubToKopecks(1499.99), Qty: 2},
	}
	got := orderpricing.TotalKopecks(lines, 10)
	const want int64 = 269998
	if got != want {
		t.Fatalf("total = %d kopecks, want %d", got, want)
	}
}
