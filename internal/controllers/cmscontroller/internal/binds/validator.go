package binds

import "errors"

var (
	ErrMissingValidator = errors.New("missing validator")
	ErrInvalidData      = errors.New("invalid data")
)

type modelValidator interface {
	Validate() error
}

type Validator struct{}

func (Validator) Validate(i interface{}) error {
	if v, ok := i.(modelValidator); ok {
		return v.Validate()
	}

	return ErrMissingValidator
}
