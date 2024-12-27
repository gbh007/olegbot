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
	BotID int64  `json:"bot_id"`
	Text  string `json:"text"`
}

func (req AddQuoteRequest) Validate() error {
	if req.BotID < 1 {
		return fmt.Errorf("%w: incorrect bot id", ErrInvalidData)
	}

	if req.Text == "" {
		return fmt.Errorf("%w: missing text", ErrInvalidData)
	}

	return nil
}

type GetQuoteRequest struct {
	ID int64 `json:"id"`
}

func (req GetQuoteRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}

type DeleteQuoteRequest struct {
	ID int64 `json:"id"`
}

func (req DeleteQuoteRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}

type QuoteListRequest struct {
	BotID int64 `json:"bot_id"`
}

func (req QuoteListRequest) Validate() error {
	if req.BotID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}
