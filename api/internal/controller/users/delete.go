package users

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/user"
)

// Delete handle soft delete user
func (i impl) Delete(ctx context.Context, id int64) error {
	u, err := i.repo.User().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	u.Status = model.UserStatusDeleted

	if err = i.repo.User().Update(ctx, u); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}
