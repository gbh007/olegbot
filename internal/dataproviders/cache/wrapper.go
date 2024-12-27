package cache

import "sync"

type wrapper[K comparable, V any] struct {
	raw   map[K]V
	mutex sync.RWMutex
}

func newWrapper[K comparable, V any]() *wrapper[K, V] {
	return &wrapper[K, V]{
		raw:   make(map[K]V),
		mutex: sync.RWMutex{},
	}
}

func (w *wrapper[K, V]) Get(key K) (V, bool) {
	w.mutex.RLock()
	v, ok := w.raw[key]
	w.mutex.RUnlock()

	return v, ok
}

func (w *wrapper[K, V]) Set(key K, value V) {
	w.mutex.Lock()
	w.raw[key] = value
	w.mutex.Unlock()
}

func (w *wrapper[K, V]) Del(key K) {
	w.mutex.Lock()
	delete(w.raw, key)
	w.mutex.Unlock()
}

func (w *wrapper[K, V]) Clear() {
	w.mutex.Lock()
	w.raw = make(map[K]V)
	w.mutex.Unlock()
}
