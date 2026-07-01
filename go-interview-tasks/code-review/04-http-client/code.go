//go:build !solution

// Задача на code review: HTTP-клиент для вызова внутренних сервисов.
package httpclient

import (
	"io"
	"net/http"
)

// Fetch загружает тело ответа по URL.
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
