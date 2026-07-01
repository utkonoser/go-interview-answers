//go:build solution

package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// fix: http.Get использует DefaultClient без таймаута → зависание на медленном upstream.
// fix: не проверяется resp.StatusCode — 500/404 вернутся как «успех».
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Fetch(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
