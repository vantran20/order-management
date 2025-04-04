package users

import (
	"context"

	"omg/api/internal/model"
)

// GetUsers retrieve all users data
func (i impl) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := i.repo.User().GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
