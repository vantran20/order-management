package generator

import (
	"omg/api/pkg/snowflake"

	pkgerrors "github.com/pkg/errors"
)

// Snowflake generators per table.
var (
	// UserIDSNF the snowflake generator for User table's ID in DB
	UserIDSNF *snowflake.Generator
	// ProductIDSNF the snowflake generator for Product table's ID in DB
	ProductIDSNF *snowflake.Generator
	// OrderIDSNF the snowflake generator for Order table's ID in DB
	OrderIDSNF *snowflake.Generator
	// OrderItemIDSNF the snowflake generator for Order Item table's ID in DB
	OrderItemIDSNF *snowflake.Generator
)

// InitSnowflakeGenerators initializes all the snowflake generators
func InitSnowflakeGenerators() error {
	var err error

	if UserIDSNF == nil {
		UserIDSNF, err = snowflake.New()
		if err != nil {
			return pkgerrors.WithStack(err)
		}
	}

	if ProductIDSNF == nil {
		ProductIDSNF, err = snowflake.New()
		if err != nil {
			return pkgerrors.WithStack(err)
		}
	}

	if OrderIDSNF == nil {
		OrderIDSNF, err = snowflake.New()
		if err != nil {
			return pkgerrors.WithStack(err)
		}
	}

	if OrderItemIDSNF == nil {
		OrderItemIDSNF, err = snowflake.New()
		if err != nil {
			return pkgerrors.WithStack(err)
		}
	}

	return nil
}
