package binds

import "fmt"

type AddModeratorRequest struct {
	UserID      int64  `json:"user_id"`
	Description string `json:"description"`
}

func (req AddModeratorRequest) Validate() error {
	if req.UserID < 1 {
		return fmt.Errorf("%w: incorrect user id %d", ErrInvalidData, req.UserID)
	}

	return nil
}
