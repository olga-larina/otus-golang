package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 { // максимум 0 ошибок
		return ErrErrorsLimitExceeded
	}

	// ошибка, может быть записана 1 раз
	var err error
	errSync := &sync.Once{}

	// помещаем задачи в канал и закрываем его
	tasksChan := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	// количество ошибок, обрабатываем атомарно
	var errCount int64

	// запускаем n горутин, добавляем в WaitGroup
	// в каждой горутине читаем задания из канала до тех пор, пока он не будет закрыт
	// если выполнение задачи завершилось ошибкой, то увеличиваем счётчик ошибок
	// перед запуском каждой задачи проверяем, не превышено ли количество ошибок
	// если превышено, то записываем ошибку и завершаем горутину
	taskWaitGroup := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		taskWaitGroup.Add(1)
		go func() {
			defer taskWaitGroup.Done()
			for task := range tasksChan {
				if atomic.LoadInt64(&errCount) >= int64(m) {
					errSync.Do(func() {
						err = ErrErrorsLimitExceeded
					})
					return
				}
				if err := task(); err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	// ожидаем завершения выполнения всех горутин
	taskWaitGroup.Wait()
	return err
}
