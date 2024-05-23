package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	out := generator(in, done)

	for _, stage := range stages {
		out = stage(executeStage(out, done))
	}

	return out
}

// генерация канала исходных данных с учётом получения сигналов о завершении.
func generator(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}
	}()

	return out
}

// обработка каналов для стейджа.
func executeStage(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		// необходимо прочитать все значения из входящего канала
		// иначе предыдущий stage может зависнуть на записи, ожидая места в канале
		for v := range in {
			// пропускаем запись прочитанных значений, если получен сигнал завершения
			select {
			case <-done:
				continue
			default:
			}

			// ожидаем возможности записи в канал или сигнала завершения
			select {
			case out <- v:
			case <-done:
				continue
			}
		}
	}()

	return out
}
