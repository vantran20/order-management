package system

import (
	"context"
)

// CheckReadiness checks if the system is ready for operation or not
func (i impl) CheckReadiness(ctx context.Context) error {
	return i.repo.System().CheckDB(ctx)
}
