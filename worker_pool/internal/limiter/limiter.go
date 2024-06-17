package limiter

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type Limiter struct {
	counterPerHostRWMu   sync.RWMutex
	requestLimitsPerHost map[string]int
	actualRequestCount   map[string]int
}

func NewLimiter() *Limiter {
	return &Limiter{
		actualRequestCount:   make(map[string]int),
		requestLimitsPerHost: requestLimitsPerHost(),
	}
}

func (l *Limiter) ShouldBeThrottled(url string) bool {
	path := l.getPath(url)
	l.counterPerHostRWMu.RLock()
	counter := l.actualRequestCount[path]
	l.counterPerHostRWMu.RUnlock()
	if counter >= l.requestLimitsPerHost[path] {
		log.Println("LIMITER THROTTLED REQUEST %s", url)
		return true
	}

	return false
}

func (l *Limiter) UpdateActualRequestCount(url string) {
	host := l.getPath(url)
	l.counterPerHostRWMu.Lock()
	l.actualRequestCount[host]++
	l.counterPerHostRWMu.Unlock()
}

func (l *Limiter) getPath(link string) string {
	dividedLink := strings.Split(link, "/")
	return strings.Join(dividedLink[0:len(dividedLink)-1], "/")
}

func requestLimitsPerHost() map[string]int {
	return map[string]int{
		"https://jsonplaceholder.typicode.com/posts":    5,
		"https://jsonplaceholder.typicode.com/comments": 3,
		"https://jsonplaceholder.typicode.com/almubs":   4,
		"https://jsonplaceholder.typicode.com/todos":    2,
		"https://jsonplaceholder.typicode.com/users":    10,
	}
}

func (l *Limiter) Start(ctx context.Context) {
	for {
		select {
		case <- ctx.Done():
			return
		case <- time.After(1 * time.Second):
			l.NullifyActualRequestCount()
		}
	}
}

func (l *Limiter) NullifyActualRequestCount() {
	l.counterPerHostRWMu.Lock()
	for k, _ := range l.actualRequestCount {
		l.actualRequestCount[k] = 0
	}
	l.counterPerHostRWMu.Unlock()
}
