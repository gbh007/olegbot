package binds

import (
	"app/internal/controllers/cmscontroller/internal/render"
	"fmt"
)

type CreateBotRequest render.Bot

func (req CreateBotRequest) Validate() error {
	if req.Token == "" {
		return fmt.Errorf("%w: incorrect token", ErrInvalidData)
	}

	if req.Name == "" {
		return fmt.Errorf("%w: incorrect name", ErrInvalidData)
	}

	if req.Tag == "" {
		return fmt.Errorf("%w: incorrect tag", ErrInvalidData)
	}

	return nil
}

type UpdateBotRequest render.Bot

func (req UpdateBotRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	if req.Token == "" {
		return fmt.Errorf("%w: incorrect token", ErrInvalidData)
	}

	if req.Name == "" {
		return fmt.Errorf("%w: incorrect name", ErrInvalidData)
	}

	if req.Tag == "" {
		return fmt.Errorf("%w: incorrect tag", ErrInvalidData)
	}

	return nil
}

type DeleteBotRequest struct {
	ID int64 `json:"id"`
}

func (req DeleteBotRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}

type GetBotRequest struct {
	ID int64 `json:"id"`
}

func (req GetBotRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}

type StartBotRequest struct {
	ID int64 `json:"id"`
}

func (req StartBotRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}

type StopBotRequest struct {
	ID int64 `json:"id"`
}

func (req StopBotRequest) Validate() error {
	if req.ID < 1 {
		return fmt.Errorf("%w: incorrect id", ErrInvalidData)
	}

	return nil
}
