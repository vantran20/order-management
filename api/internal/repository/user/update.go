package user

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) Update(ctx context.Context, m model.User) error {
	u, err := orm.FindUser(ctx, i.dbConn, m.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkgerrors.WithStack(ErrNotFound)
		}

		return pkgerrors.WithStack(err)
	}

	u.Name = m.Name
	u.Email = m.Email
	u.Password = m.Password
	u.Status = m.Status.String()

	if _, err = u.Update(ctx, i.dbConn, boil.Whitelist(
		orm.UserColumns.Name,
		orm.UserColumns.Email,
		orm.UserColumns.Password,
		orm.UserColumns.Status,
		orm.UserColumns.UpdatedAt,
	)); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}
