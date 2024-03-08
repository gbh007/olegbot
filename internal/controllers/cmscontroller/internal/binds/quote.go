package binds

import "fmt"

type UpdateQuoteRequest struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

func (req UpdateQuoteRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id %d", ErrInvalidData, req.ID)
	}

	if req.Text == "" {
		return fmt.Errorf("%w: missing text", ErrInvalidData)
	}

	return nil
}

type AddQuoteRequest struct {
	Text string `json:"text"`
}

func (req AddQuoteRequest) Validate() error {
	if req.Text == "" {
		return fmt.Errorf("%w: missing text", ErrInvalidData)
	}

	return nil
}
