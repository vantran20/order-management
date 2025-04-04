package user

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
)

// GetByEmail retrieve the user data by email
func (i impl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	o, err := orm.Users(
		orm.UserWhere.Email.EQ(email),
	).One(ctx, i.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, pkgerrors.WithStack(ErrNotFound)
		}
		return model.User{}, pkgerrors.WithStack(err)
	}

	return toUser(o), nil
}
