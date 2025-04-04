package users

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

// Create handles user registration
func (i impl) Create(ctx context.Context, inp model.CreateUserInput) (model.User, error) {
	// Check if user with this email already exists
	_, err := i.repo.User().GetByEmail(ctx, inp.Email)
	if err != nil {
		if !errors.Is(err, user.ErrNotFound) {
			return model.User{}, err
		}
	} else {
		return model.User{}, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inp.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, ErrHashedPassword
	}

	m := model.User{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: string(hashedPassword),
		Status:   model.UserStatusActive,
	}

	return i.repo.User().CreateUser(ctx, m)
}
