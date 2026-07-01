// Пакет workerpool — пул воркеров: фиксированное число горутин обрабатывает очередь задач.
//
// Вместо «одна горутина на задачу» переиспользуем воркеров и ограничиваем параллелизм.
package workerpool

import "sync"

// Job — одна задача для воркера.
type Job func()

// Run запускает numWorkers горутин. Канал jobs закрывает отправитель — тогда воркеры завершаются.
func Run(jobs <-chan Job, numWorkers int) {
	if numWorkers < 1 {
		numWorkers = 1
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for range numWorkers {
		go func() {
			defer wg.Done()
			for job := range jobs { // читаем, пока jobs не закроют
				job()
			}
		}()
	}

	wg.Wait() // ждём, пока все воркеры обработают очередь
}
