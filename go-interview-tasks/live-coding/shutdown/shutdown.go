// Пакет shutdown — корректная остановка HTTP-сервера (graceful shutdown).
//
// По SIGINT/SIGTERM перестаём принимать новые соединения и ждём завершения текущих запросов.
package shutdown

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Serve слушает addr, обрабатывает handler; по Ctrl+C / SIGTERM — Shutdown с таймаутом.
func Serve(addr string, handler http.Handler, shutdownTimeout time.Duration) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // k8s шлёт SIGTERM при остановке pod
	defer signal.Stop(quit)
	return ServeUntil(addr, handler, shutdownTimeout, quit)
}

// ServeUntil то же, что Serve, но источник сигнала задаётся снаружи (удобно в тестах).
func ServeUntil(addr string, handler http.Handler, shutdownTimeout time.Duration, stop <-chan os.Signal) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return ServeListener(ln, handler, shutdownTimeout, stop)
}

// ServeListener запускает сервер на уже созданном listener.
func ServeListener(ln net.Listener, handler http.Handler, shutdownTimeout time.Duration, stop <-chan os.Signal) error {
	server := &http.Server{Handler: handler}

	errCh := make(chan error, 1)
	go func() {
		// Serve блокируется; ErrServerClosed — норма после Shutdown
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err // сервер упал до сигнала остановки
	case <-stop:
		// получили SIGTERM / SIGINT — начинаем graceful shutdown
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown: не принимает новые conn, ждёт активные запросы (до таймаута)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
		return err
	}
	return nil
}
