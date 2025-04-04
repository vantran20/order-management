package users

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/user"
)

// GetByEmail handles get user data by email
func (i impl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := i.repo.User().GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	return u, nil
}
