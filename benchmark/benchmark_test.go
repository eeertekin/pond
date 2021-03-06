package benchmark

import (
	"sync"
	"testing"
	"time"

	"github.com/alitto/pond"
	"github.com/gammazero/workerpool"
	"github.com/panjf2000/ants/v2"
)

const (
	taskCount    = 1000000
	taskDuration = 10 * time.Millisecond
	workerCount  = 200000
)

func testFunc() {
	time.Sleep(taskDuration)
}

func BenchmarkPond(b *testing.B) {
	var wg sync.WaitGroup
	pool := pond.New(workerCount, taskCount)
	defer pool.StopAndWait()

	// Submit tasks
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(taskCount)
		for i := 0; i < taskCount; i++ {
			pool.Submit(func() {
				testFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup

	// Submit tasks
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(taskCount)
		for i := 0; i < taskCount; i++ {
			go func() {
				testFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoroutinePool(b *testing.B) {
	var wg sync.WaitGroup

	// Submit tasks
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskChan := make(chan func())
		wg.Add(workerCount)
		// Start worker goroutines
		for i := 0; i < workerCount; i++ {
			go func() {
				for task := range taskChan {
					task()
				}
				wg.Done()
			}()
		}
		// Submit tasks
		for i := 0; i < taskCount; i++ {
			taskChan <- testFunc
		}
		close(taskChan)
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkBufferedGoroutinePool(b *testing.B) {
	var wg sync.WaitGroup

	// Submit tasks
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskChan := make(chan func(), taskCount)
		wg.Add(workerCount)
		// Start worker goroutines
		for i := 0; i < workerCount; i++ {
			go func() {
				for task := range taskChan {
					task()
				}
				wg.Done()
			}()
		}
		// Submit tasks
		for i := 0; i < taskCount; i++ {
			taskChan <- testFunc
		}
		close(taskChan)
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkGammazeroWorkerpool(b *testing.B) {
	var wg sync.WaitGroup
	wp := workerpool.New(workerCount)
	defer wp.StopWait()

	// Submit tasks
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(taskCount)
		for i := 0; i < taskCount; i++ {
			wp.Submit(func() {
				testFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkAnts(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(workerCount, ants.WithExpiryDuration(10*time.Second))
	defer p.Release()

	// Submit tasks
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(taskCount)
		for i := 0; i < taskCount; i++ {
			_ = p.Submit(func() {
				testFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}
