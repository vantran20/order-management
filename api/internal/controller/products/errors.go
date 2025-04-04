package products

import (
	"errors"
)

var (
	ErrNotFound             = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductDeleted       = errors.New("product deleted")
)
