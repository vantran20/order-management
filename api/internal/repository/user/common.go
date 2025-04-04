package user

import (
	"omg/api/internal/model"
	"omg/api/internal/repository/orm"
)

func toUser(o *orm.User) model.User {
	return model.User{
		ID:        o.ID,
		Name:      o.Name,
		Email:     o.Email,
		Password:  o.Password,
		Status:    model.UserStatus(o.Status),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
