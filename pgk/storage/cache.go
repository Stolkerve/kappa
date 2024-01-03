package storage

import (
	"sync"
)

var Cache Queue

type Queue struct {
	items []Function
	l     sync.RWMutex
}

const MAX_FUNCTIONS_IN_CACHE = 10

func (q *Queue) Push(data Function) {
	q.l.Lock()
	defer q.l.Unlock()

	if len(q.items) == MAX_FUNCTIONS_IN_CACHE {
		q.items = q.items[1:]
	}
	q.items = append(q.items, data)
}

func (q *Queue) Get(id string) []byte {
	q.l.RLock()
	defer q.l.RUnlock()

	for _, v := range q.items {
		if v.ID == id {
			return v.Wasm
		}
	}
	return nil
}
