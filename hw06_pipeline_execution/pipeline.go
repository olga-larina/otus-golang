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

	// первый стейдж получает на вход исходный канал, остальные стейджи - каналы, модицифированные предыдущими стейджами
	out := in

	for _, stage := range stages {
		out = stage(executeStage(out, done))
	}

	return out
}

func executeStage(in In, done In) Out {
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
