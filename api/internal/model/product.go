package model

import (
	"time"
)

// ProductStatus represents the status of the product
type ProductStatus string

const (
	// ProductStatusActive means the product is active
	ProductStatusActive ProductStatus = "ACTIVE"
	// ProductStatusDeleted means the product is deleted
	ProductStatusDeleted ProductStatus = "DELETED"
)

// String converts to string value
func (p ProductStatus) String() string {
	return string(p)
}

// IsValid checks if plan status is valid
func (p ProductStatus) IsValid() bool {
	switch p {
	case ProductStatusActive, ProductStatusDeleted:
		return true
	}
	return false
}

// Product represents the product to be sold
type Product struct {
	ID          int64
	Name        string
	Description string
	Status      ProductStatus
	Price       float64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateProductInput holds input params for creating the product
type CreateProductInput struct {
	Name  string
	Desc  string
	Price float64
	Stock int64
}

// UpdateProductInput holds input params for updating the product
type UpdateProductInput struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Stock       int64
	Status      ProductStatus
}
