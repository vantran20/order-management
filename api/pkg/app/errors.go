package app

import (
	"errors"
)

var (
	// ErrInvalidAppConfig means the Config is invalid
	ErrInvalidAppConfig = errors.New("invalid app config")
)
