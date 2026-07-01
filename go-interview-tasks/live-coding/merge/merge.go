// Пакет merge — объединение нескольких входных каналов в один выходной.
//
// На каждый вход — горутина-перекачка; out закрывается, когда все входы исчерпаны.
package merge

import "sync"

// Int сливает каналы в один. Закрытие out — после wg.Wait() всех перекачек.
func Int(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			out <- v // пересылаем в общий канал
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch) // ch передаём аргументом — иначе замыкание на последний ch цикла
	}

	go func() {
		wg.Wait()
		close(out) // только когда все входные каналы прочитаны до конца
	}()

	return out
}
