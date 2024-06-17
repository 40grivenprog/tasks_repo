// Название пакетов в множественном числе или не ?
// Размер буффера
// Как лучше называть интерфейсы с большой или малой ? по моей логике лучше с малой так как нам не нужно их экспортить

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	filewriter "github.com/pool_with_limiter/internal/file_writer"
	"github.com/pool_with_limiter/internal/limiter"
	"github.com/pool_with_limiter/internal/worker"
	workerpool "github.com/pool_with_limiter/internal/worker_pool"
)

const WorkersNum int = 5
const TasksCount = 30

var links []string = []string{
	"https://jsonplaceholder.typicode.com/posts",
	"https://jsonplaceholder.typicode.com/comments",
	"https://jsonplaceholder.typicode.com/almubs",
	"https://jsonplaceholder.typicode.com/todos",
	"https://jsonplaceholder.typicode.com/users",
}

type workerPoolI interface {
	AddTask(task string)
	Stop()
	Result()
}

type workerI interface {
	Start(ctx context.Context)
}

type limiterI interface {
	Start(ctx context.Context)
}

func getInputLinks() []string {
	urls := make([]string, TasksCount)
	for i := 0; i < TasksCount; i++ {
		urls[i] = fmt.Sprintf("%s/%d", links[rand.Intn(4)], rand.Intn(30))
	}
	return urls
}

func populateTaskQueue(ctx context.Context, wp workerPoolI) {
	links := getInputLinks()
	for _, link := range links {
		select {
		case <-ctx.Done():
			return
		default:
			wp.AddTask(link)
		}
	}
	wp.Stop()
}

func startLimiter(ctx context.Context, l limiterI) {
	l.Start(ctx)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileWriter := filewriter.NewFileWriter()
	wp := workerpool.NewWorkerPool(fileWriter)
	go populateTaskQueue(ctx, wp)
	limiter := limiter.NewLimiter()
	go startLimiter(ctx, limiter)
	wg := new(sync.WaitGroup)
	wg.Add(WorkersNum)

	workers := make([]workerI, WorkersNum)

	for i, _ := range workers {
		workers[i] = worker.NewWorker(wg, wp, limiter)
	}
	for _, worker := range workers {
		go worker.Start(ctx)
	}

	wg.Wait()
	wp.Result()
}
