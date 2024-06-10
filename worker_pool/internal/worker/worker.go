package worker

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type WorkerPoolI interface {
	GetTaskQueue() <-chan string
	PutFailedTask(failedTask string, err error)
	PutSuccessTask(task string)
}

type LimiterI interface {
	ShouldBeThrottled(url string) bool
	UpdateActualRequestCount(url string)
}

type Worker struct {
	wp      WorkerPoolI
	limiter LimiterI
	wg      *sync.WaitGroup
}

func NewWorker(wg *sync.WaitGroup, wp WorkerPoolI, limiter LimiterI) *Worker {
	return &Worker{wg: wg, wp: wp, limiter: limiter}
}

func (w *Worker) Start(ctx context.Context) {
	log.Println("!!!!!!Worker started!!!!!!")
	for task := range w.wp.GetTaskQueue() {
		select {
		case <-ctx.Done():
			return
		default:
			log.Println("WORKER PROCESSING TASK")
			if err := w.downloadUrl(ctx, task); err != nil {
				w.wp.PutFailedTask(task, err)
			} else {
				w.wp.PutSuccessTask(task)
			}
		}
	}
	w.wg.Done()
	log.Println("!!!!!Worker DONE!!!!!!")
}

func (w *Worker) downloadUrl(ctx context.Context, url string) error {
	for i := 0; i <= 5; i++ {
		if w.limiter.ShouldBeThrottled(url) {
			continue
		}
		time.Sleep(time.Duration(i) * 100 * time.Millisecond)
		if rand.Intn(3) == 0 {
			return fmt.Errorf("error downloading %s", url)
		}
		w.limiter.UpdateActualRequestCount(url)
		return nil
	}
	return nil
}
