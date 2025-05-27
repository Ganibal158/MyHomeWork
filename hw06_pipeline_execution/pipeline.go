package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = processingWithDone(out, done, stage)
	}
	return out
}

func processingWithDone(in In, done In, stage Stage) Out {
	stageOut := make(Bi)
	go func() {
		defer close(stageOut)
		stageInput := make(Bi)
		go func() {
			defer close(stageInput)
			for {
				select {
				case <-done:
					break // выход при закрытии канала
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
						return // завершение при закрытии канала с отправкой выходного канала

					case stageInput <- val:
					}
				}
			}
		}()
		for val := range stage(stageInput) {
			select {
			case <-done:
				return
			case stageOut <- val:
			}
		}
	}()

	return stageOut
}
