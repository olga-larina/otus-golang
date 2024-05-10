package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 { // максимум 0 ошибок
		return ErrErrorsLimitExceeded
	}
	if len(tasks) == 0 {
		return nil
	}

	// итоговая ошибка
	var err error

	// канал ошибок
	errChan := make(chan error, m)
	defer close(errChan)

	// помещаем задачи в канал, когда есть место, или завершаем работу при превышении количества ошибок
	tasksChan := make(chan Task)
	go func() {
		defer close(tasksChan)

		errCount := 0
		taskIdx := 0
		for taskIdx < len(tasks) {
			select {
			case <-errChan:
				errCount++
				if errCount >= m {
					err = ErrErrorsLimitExceeded
					return
				}
			case tasksChan <- tasks[taskIdx]:
				taskIdx++
			}
		}
	}()

	// запускаем n горутин, добавляем в WaitGroup
	// в каждой горутине читаем задания из канала до тех пор, пока он не будет закрыт
	// если выполнение задачи завершилось ошибкой, то отправляем ошибку в канал ошибок
	taskWaitGroup := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		taskWaitGroup.Add(1)
		go func() {
			defer taskWaitGroup.Done()
			for task := range tasksChan {
				if err := task(); err != nil {
					errChan <- err
				}
			}
		}()
	}

	// ожидаем завершения выполнения всех горутин
	taskWaitGroup.Wait()

	return err
}
