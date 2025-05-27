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
					break
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
						return // возможно косяк тут и нужно сделать return/ break/ очистить канал/ переделать нахуй эту функцию чтобы она через sink.Once наращивала переменную по которой будут закрываться каналы

					case stageInput <- val: // возможно тоже не совсем верно
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
