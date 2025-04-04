package user

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
)

// GetByID retrieve the user data by user id
func (i impl) GetByID(ctx context.Context, id int64) (model.User, error) {
	o, err := orm.Users(
		orm.UserWhere.ID.EQ(id),
	).One(ctx, i.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, pkgerrors.WithStack(ErrNotFound)
		}
		return model.User{}, pkgerrors.WithStack(err)
	}

	return toUser(o), nil
}
