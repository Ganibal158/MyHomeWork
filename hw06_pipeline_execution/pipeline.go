package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = wrapWithDone(stage, out, done)
	}
	return out
}

func wrapWithDone(stage Stage, in In, done In) Out {
	var wg sync.WaitGroup
	stageOut := make(Bi)
	go func() {
		defer close(stageOut)
		stageInput := make(Bi)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(stageInput)
			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
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
		wg.Wait()
	}()

	return stageOut
}
