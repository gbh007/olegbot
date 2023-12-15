package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

func (r *Repository) Load(_ context.Context, filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("load: open file: %w", err)
	}

	data := make([]string, 0)

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return fmt.Errorf("load: decode: %w", err)
	}

	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	r.data = data

	return nil
}
