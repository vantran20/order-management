package user

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetUsers retrieve all users active data
func (i impl) GetUsers(ctx context.Context) ([]model.User, error) {
	o, err := orm.Users(
		orm.UserWhere.Status.NEQ(model.UserStatusDeleted.String()),
		qm.OrderBy(orm.UserColumns.CreatedAt+" DESC"),
	).All(ctx, i.dbConn)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	var users []model.User
	for _, u := range o {
		users = append(users, toUser(u))
	}
	return users, nil
}
