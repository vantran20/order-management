package inventory

import (
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")
)
