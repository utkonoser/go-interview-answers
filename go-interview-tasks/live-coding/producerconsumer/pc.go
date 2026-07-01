// Пакет producerconsumer — классический паттерн: один пишет в канал, другой читает.
//
// Закрытие канала отправителем — сигнал consumer'у, что данных больше не будет.
package producerconsumer

// Produce кладёт значения 0..n-1 в канал и закрывает его.
func Produce(items chan<- int, n int) {
	for i := range n {
		items <- i
	}
	close(items) // без close consumer зависнет в range
}

// Consume читает канал до закрытия и вызывает fn на каждом элементе.
func Consume(items <-chan int, fn func(int)) {
	for item := range items { // выходим, когда producer закроет канал
		fn(item)
	}
}

// Run запускает producer и consumer, блокируется до конца обработки.
func Run(n int, fn func(int)) {
	items := make(chan int) // без буфера — синхронная передача
	done := make(chan struct{})

	go func() {
		Produce(items, n)
	}()

	go func() {
		Consume(items, fn)
		close(done) // consumer закончил — можно разбудить main
	}()

	<-done // main ждёт завершения consumer
}
