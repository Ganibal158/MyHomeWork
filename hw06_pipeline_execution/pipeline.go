package hw06pipelineexecution

import (
	"log/slog"
	"strconv"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for pos, stage := range stages {
		out = processingWithDone(pos, out, done, stage)
	}
	return out
}

func processingWithDone(pos int, in In, done In, stage Stage) Out {
	stageName := strconv.Itoa(pos)
	stageOut := stage(in) // здесь мы создаём новую горутину, которая читает из in и отправляет значения в stageOut
	out := make(Bi)
	go func() {
		// получает значения из stageOut и отправляет их в out, если не получен done
		isDone := false

		for {
			if !isDone {
				select {
				case <-done:
					// не прекращаем слушать stageOut, просто закрываем out
					close(out)
					isDone = true
				case val, ok := <-stageOut:
					if !ok {
						close(out)
						return
					}
					out <- val
				}
			}
			if isDone {
				v, ok := <-stageOut
				if ok {
					slog.Info("processingWithDone got value after done", "stage", stageName, "value", v)
				} else {
					return
				}
			}
		}
	}()

	return out
}
