package workerpool

import (
	"fmt"
	"sync"
)

const WorkerNum int = 5

type FileWriterI interface {
	WriteResult(filePath string, result []string) error
}

type WorkerPool struct {
	fileWriter           FileWriterI
	taskQueue            chan string
	failedTasks          []string
	successTasks         []string
	failedTasksMu        sync.Mutex
	successfailedTasksMu sync.Mutex
}

func NewWorkerPool(fileWriter FileWriterI) *WorkerPool {
	return &WorkerPool{
		fileWriter: fileWriter,
		taskQueue:  make(chan string, WorkerNum),
	}
}

func (wp *WorkerPool) GetTaskQueue() <-chan string {
	return wp.taskQueue
}

func (wp *WorkerPool) PutFailedTask(failedTask string, err error) {
	wp.failedTasksMu.Lock()
	wp.failedTasks = append(wp.failedTasks, fmt.Sprintf("Failed task: %s, Err: %s", failedTask, err.Error()))
	wp.failedTasksMu.Unlock()
}

func (wp *WorkerPool) PutSuccessTask(successTask string) {
	wp.successfailedTasksMu.Lock()
	wp.successTasks = append(wp.successTasks, fmt.Sprintf("Success task: %s", successTask))
	wp.successfailedTasksMu.Unlock()
}

func (wp *WorkerPool) AddTask(task string) {
	wp.taskQueue <- task
}

func (wp *WorkerPool) Result() {
	fmt.Println("FAILED:!!!!!")
	fmt.Println(len(wp.failedTasks))
	fmt.Println("SUCCESS:!!!!!")
	fmt.Println(len(wp.successTasks))
	wp.fileWriter.WriteResult("errors", wp.failedTasks)
	wp.fileWriter.WriteResult("success", wp.successTasks)
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
}
