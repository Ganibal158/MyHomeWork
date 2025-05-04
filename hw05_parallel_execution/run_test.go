package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("m less zero", func(t *testing.T) {
		var runCounter int32
		tasks := make([]Task, 10) // Создаём десять задач, которые при запуске инкрементируют счётчик
		for i := 0; i < 10; i++ {
			tasks[i] = func() error {
				atomic.AddInt32(&runCounter, 1)
				return nil
			}
		}
		err := Run(tasks, 5, -1) // Согласно принятой логике, ни одна из задач не должна запустится
		if !errors.Is(err, ErrErrorsLimitExceeded) {
			t.Errorf("expected error %v, got %v", ErrErrorsLimitExceeded, err)
		}
		if atomic.LoadInt32(&runCounter) != 0 {
			t.Errorf("expected runCounter to be 0, got %d", runCounter)
		}
	})

	t.Run("Start with zero tusks", func(t *testing.T) {
		workersCount := 10
		maxErrorsCount := 1
		done := make(chan struct{})
		go func() { // Вызываем Run в отдельной горутине, чтобы проверить время её закрытия.
			err := Run(nil, workersCount, maxErrorsCount)
			require.NoError(t, err)
			close(done)
		}()
		select {
		case <-done: // канал закрылся вовремя, утечек быть не должно
		case <-time.After(1 * time.Second):
			t.Fatal("Run with nil tasks didn't return in time — possible goroutine leak")
		}
	})

	t.Run("Parallelism test", func(t *testing.T) {
		tasksCount := 50
		sleepPerTask := 100 * time.Millisecond
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(sleepPerTask)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}
		workersCount := 5
		maxErrorsCount := 1
		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsed := time.Since(start)
		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		expectedSequential := sleepPerTask * time.Duration(tasksCount)
		require.Less(t, elapsed, expectedSequential/2, "tasks likely ran sequentially")
	})
}
