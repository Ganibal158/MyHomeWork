package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var (
		wg       sync.WaitGroup
		errCount int32
		done     = make(chan struct{}) // Канал по закрытию которого будут происходить завершение работы
		tasksCh  = make(chan Task)     // Канал для передачи задачи воркеру
		once     sync.Once
	)
	wg.Add(n)
	for i := 0; i < n; i++ { // Создаём n воркеров
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done: // Заврершение работы воркера при чтении/закрытии канала done
					return
				case task, ok := <-tasksCh: // Чтение задачи из канала
					if !ok {
						return
					}
					if err := task(); err != nil {
						if int(atomic.AddInt32(&errCount, 1)) >= m {
							once.Do(func() { // безопасное закрытие канала
								close(done)
							})
						}
					}
				}
			}
		}()
	}
	go func() { // Горутина для передачи задач через канал
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case <-done:
				return
			case tasksCh <- task:
			}
		}
	}()
	wg.Wait()
	if int(atomic.LoadInt32(&errCount)) >= m { // При превышении кол-ва ошибок, передаём ошибку
		return ErrErrorsLimitExceeded
	}
	return nil
}
