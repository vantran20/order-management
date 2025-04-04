package users

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/user"
)

// GetByID handles get user data by user id
func (i impl) GetByID(ctx context.Context, id int64) (model.User, error) {
	u, err := i.repo.User().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	return u, nil
}
