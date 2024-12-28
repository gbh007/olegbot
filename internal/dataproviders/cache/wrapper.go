package cache

import (
	"log/slog"
	"sync"
)

type wrapper[K comparable, V any] struct {
	raw   map[K]V
	mutex sync.RWMutex

	name   string
	logger *slog.Logger
}

func newWrapper[K comparable, V any](
	name string,
	logger *slog.Logger,
) *wrapper[K, V] {
	return &wrapper[K, V]{
		raw:    make(map[K]V),
		mutex:  sync.RWMutex{},
		name:   name,
		logger: logger,
	}
}

func (w *wrapper[K, V]) Get(key K) (V, bool) {
	w.logger.Debug("cache Get", slog.String("cache", w.name), slog.Any("key", key))

	w.mutex.RLock()
	v, ok := w.raw[key]
	w.mutex.RUnlock()

	return v, ok
}

func (w *wrapper[K, V]) Set(key K, value V) {
	w.logger.Debug("cache Set", slog.String("cache", w.name), slog.Any("key", key))

	w.mutex.Lock()
	w.raw[key] = value
	w.mutex.Unlock()
}

func (w *wrapper[K, V]) Del(key K) {
	w.logger.Debug("cache Del", slog.String("cache", w.name), slog.Any("key", key))

	w.mutex.Lock()
	delete(w.raw, key)
	w.mutex.Unlock()
}

func (w *wrapper[K, V]) Clear() {
	w.logger.Debug("cache Clear", slog.String("cache", w.name))

	w.mutex.Lock()
	w.raw = make(map[K]V)
	w.mutex.Unlock()
}
