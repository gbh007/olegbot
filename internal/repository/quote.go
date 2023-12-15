package repository

import (
	"context"
	"math/rand"
)

func (r *Repository) RandomQuote(_ context.Context) (string, error) {
	r.dataMutex.RLock()
	defer r.dataMutex.RUnlock()

	return r.data[rand.Intn(len(r.data))], nil
}

func (r *Repository) AddQuote(_ context.Context, text string) error {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	r.data = append(r.data, text)

	return nil
}
