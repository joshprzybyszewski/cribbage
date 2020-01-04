package mysql

import (
	"context"
	"database/sql"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	sqlErrCodeDuplicateEntry = 1062

	dbName          = `cribbage`
	playerTableName = `players`
)

func New(ctx context.Context) (persistence.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/cribbage")
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, db)
	if err != nil {
		return nil, err
	}
	return persistence.New(nil, ps, nil), nil
}
