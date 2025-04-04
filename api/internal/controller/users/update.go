package users

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

func (i impl) Update(ctx context.Context, inp model.UpdateUserInput) (model.User, error) {
	// Check if user with this id already exists
	u, err := i.repo.User().GetByID(ctx, inp.ID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inp.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, ErrHashedPassword
	}

	userUpToDate := model.User{
		ID:       u.ID,
		Name:     inp.Name,
		Email:    inp.Email,
		Status:   inp.Status,
		Password: string(hashedPassword),
	}

	if err = i.repo.User().Update(ctx, userUpToDate); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return userUpToDate, nil
}
