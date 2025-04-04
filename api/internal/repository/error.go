package repository

import "errors"

var (
	errNestedTx = errors.New("db txn nested in db txn")
)
