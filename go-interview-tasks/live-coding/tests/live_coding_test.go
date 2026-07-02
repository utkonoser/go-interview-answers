package tests

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"interviews/go-interview-tasks/live-coding/counter"
	"interviews/go-interview-tasks/live-coding/errgroup"
	"interviews/go-interview-tasks/live-coding/lru"
	"interviews/go-interview-tasks/live-coding/merge"
	"interviews/go-interview-tasks/live-coding/producerconsumer"
	"interviews/go-interview-tasks/live-coding/ratelimit"
	"interviews/go-interview-tasks/live-coding/retry"
	"interviews/go-interview-tasks/live-coding/semaphore"
	"interviews/go-interview-tasks/live-coding/shutdown"
	"interviews/go-interview-tasks/live-coding/singleflight"
	"interviews/go-interview-tasks/live-coding/timeout"
	"interviews/go-interview-tasks/live-coding/workerpool"
)

func TestLRU(t *testing.T) {
	t.Parallel()

	c := lru.New(2)
	c.Put(1, 10)
	c.Put(2, 20)

	if v, ok := c.Get(1); !ok || v != 10 {
		t.Fatalf("Get(1) = %d, %v", v, ok)
	}

	c.Put(3, 30)

	if _, ok := c.Get(2); ok {
		t.Fatal("key 2 should be evicted")
	}
	if v, ok := c.Get(3); !ok || v != 30 {
		t.Fatalf("Get(3) = %d, %v", v, ok)
	}
}

func TestAtomicCounter(t *testing.T) {
	t.Parallel()

	var c counter.Atomic
	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()

	if c.Value() != n {
		t.Fatalf("value = %d, want %d", c.Value(), n)
	}
}

func TestWorkerPool(t *testing.T) {
	t.Parallel()

	jobs := make(chan workerpool.Job, 10)
	var processed atomic.Int32

	go func() {
		for i := range 10 {
			i := i
			jobs <- func() { processed.Add(int32(i)) }
		}
		close(jobs)
	}()

	workerpool.Run(jobs, 3)

	if processed.Load() != 45 {
		t.Fatalf("sum = %d, want 45", processed.Load())
	}
}

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	rl := ratelimit.New(2, time.Second)
	defer rl.Stop()

	if !rl.Allow() || !rl.Allow() {
		t.Fatal("expected two immediate allows from full bucket")
	}
	if rl.Allow() {
		t.Fatal("third allow should fail on empty bucket")
	}
}

func TestProducerConsumer(t *testing.T) {
	t.Parallel()

	var sum atomic.Int32
	producerconsumer.Run(10, func(i int) {
		sum.Add(int32(i))
	})
	if sum.Load() != 45 {
		t.Fatalf("sum = %d, want 45", sum.Load())
	}
}

func TestMerge(t *testing.T) {
	t.Parallel()

	ch1 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	close(ch1)

	ch2 := make(chan int, 1)
	ch2 <- 3
	close(ch2)

	var got []int
	for v := range merge.Int(ch1, ch2) {
		got = append(got, v)
	}
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
}

func TestTimeoutOK(t *testing.T) {
	t.Parallel()

	err := timeout.Do(func() error { return nil }, time.Second)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTimeoutExceeded(t *testing.T) {
	t.Parallel()

	err := timeout.Do(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}, 5*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
}

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()

	stop := make(chan os.Signal, 1)
	addrCh := make(chan string, 1)

	go func() {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Error(err)
			return
		}
		addrCh <- ln.Addr().String()
		err = shutdown.ServeListener(ln, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), time.Second, stop)
		if err != nil {
			t.Error(err)
		}
	}()

	var addr string
	select {
	case addr = <-addrCh:
	case <-time.After(time.Second):
		t.Fatal("server did not start")
	}

	resp, err := http.Get("http://" + addr)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	stop <- os.Interrupt
	time.Sleep(50 * time.Millisecond)
}

func TestErrgroupCancelOnError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	errSentinel := errors.New("boom")

	started := make(chan struct{}, 1)
	err := errgroup.RunParallel(ctx,
		func(ctx context.Context) error {
			close(started)
			return errSentinel
		},
		func(ctx context.Context) error {
			<-started
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second):
				t.Fatal("second task should be cancelled")
				return nil
			}
		},
	)
	if !errors.Is(err, errSentinel) {
		t.Fatalf("err = %v", err)
	}
}

func TestErrgroupParallelFetch(t *testing.T) {
	t.Parallel()

	if err := errgroup.ParallelFetch(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSingleflightOneCall(t *testing.T) {
	t.Parallel()

	var g singleflight.Group
	var calls atomic.Int32

	const n = 20
	var wg sync.WaitGroup
	wg.Add(n)

	for range n {
		go func() {
			defer wg.Done()
			_, _ = g.Load("sku-1", func() (any, error) {
				calls.Add(1)
				time.Sleep(20 * time.Millisecond)
				return 42, nil
			})
		}()
	}
	wg.Wait()

	if calls.Load() != 1 {
		t.Fatalf("loader called %d times, want 1", calls.Load())
	}
}

func TestSemaphoreLimitsConcurrency(t *testing.T) {
	t.Parallel()

	sem := semaphore.New(2)
	var active atomic.Int32
	var peak atomic.Int32

	const n = 10
	var wg sync.WaitGroup
	wg.Add(n)

	ctx := context.Background()
	for range n {
		go func() {
			defer wg.Done()
			if err := sem.Acquire(ctx); err != nil {
				t.Error(err)
				return
			}
			defer sem.Release()

			cur := active.Add(1)
			for {
				old := peak.Load()
				if cur <= old || peak.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			active.Add(-1)
		}()
	}
	wg.Wait()

	if peak.Load() > 2 {
		t.Fatalf("peak concurrency = %d, want <= 2", peak.Load())
	}
}

func TestSemaphoreFetchURLsParallel(t *testing.T) {
	t.Parallel()

	urls := []string{"/a", "/b", "/c", "/d", "/e"}
	if err := semaphore.FetchURLsParallel(context.Background(), urls, 2); err != nil {
		t.Fatal(err)
	}
}

func TestRetrySuccessAfterFailure(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	err := retry.Do(context.Background(), retry.Config{
		MaxAttempts: 3,
		BaseDelay:   time.Millisecond,
		MaxDelay:    10 * time.Millisecond,
	}, func() error {
		if attempts.Add(1) < 2 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("attempts = %d, want 2", attempts.Load())
	}
}

func TestRetryRespectsContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := retry.Do(ctx, retry.Config{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    time.Second,
	}, func() error {
		return errors.New("always fail")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
}
