package repository

import "sync"

type Repository struct {
	data      []string
	dataMutex *sync.RWMutex
}

func New() *Repository {
	return &Repository{
		data:      make([]string, 0),
		dataMutex: &sync.RWMutex{},
	}
}
