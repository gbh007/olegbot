package binds

import "fmt"

type AddModeratorRequest struct {
	BotID       int64  `json:"bot_id"`
	UserID      int64  `json:"user_id"`
	Description string `json:"description"`
}

func (req AddModeratorRequest) Validate() error {
	if req.BotID < 1 {
		return fmt.Errorf("%w: incorrect bot id %d", ErrInvalidData, req.BotID)
	}

	if req.UserID < 1 {
		return fmt.Errorf("%w: incorrect user id %d", ErrInvalidData, req.UserID)
	}

	return nil
}

type DeleteModeratorRequest struct {
	BotID  int64 `json:"bot_id"`
	UserID int64 `json:"user_id"`
}

func (req DeleteModeratorRequest) Validate() error {
	if req.BotID < 1 {
		return fmt.Errorf("%w: incorrect bot id %d", ErrInvalidData, req.BotID)
	}

	if req.UserID < 1 {
		return fmt.Errorf("%w: incorrect user id %d", ErrInvalidData, req.UserID)
	}

	return nil
}

type ListModeratorRequest struct {
	BotID int64 `json:"bot_id"`
}

func (req ListModeratorRequest) Validate() error {
	if req.BotID < 1 {
		return fmt.Errorf("%w: incorrect bot id %d", ErrInvalidData, req.BotID)
	}

	return nil
}
