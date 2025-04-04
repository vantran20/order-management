package system

import (
	"context"

	pkgerrors "github.com/pkg/errors"
)

// CheckDB will check if calls to DB are successful or not
func (i impl) CheckDB(ctx context.Context) error {
	_, err := i.dbConn.ExecContext(ctx, "SELECT 1")

	return pkgerrors.WithStack(err)
}
