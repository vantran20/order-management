package model

import "time"

// UserStatus represents the status of the user
type UserStatus string

const (
	// UserStatusActive means the user is active
	UserStatusActive UserStatus = "ACTIVE"
	// UserStatusDeleted means the user is deleted
	UserStatusDeleted UserStatus = "DELETED"
)

// String converts to string value
func (u UserStatus) String() string {
	return string(u)
}

// IsValid checks if plan status is valid
func (u UserStatus) IsValid() bool {
	switch u {
	case UserStatusActive, UserStatusDeleted:
		return true
	}
	return false
}

// User presents the user struct
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserInput presents creat user input
type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

// UpdateUserInput presents update user input
type UpdateUserInput struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Status   UserStatus
}
