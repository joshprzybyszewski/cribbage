package mysql

import (
	"context"
	"database/sql"
)

type txWrapper struct {
	db *sql.DB
	tx *sql.Tx
}

func (t *txWrapper) Exec(query string, ifs ...interface{}) (sql.Result, error) {
	if t.tx != nil {
		return t.tx.Exec(query, ifs...)
	}
	return t.db.Exec(query, ifs...)
}

func (t *txWrapper) ExecContext(ctx context.Context, query string, ifs ...interface{}) (sql.Result, error) {
	if t.tx != nil {
		return t.tx.ExecContext(ctx, query, ifs...)
	}
	return t.db.ExecContext(ctx, query, ifs...)
}

func (t *txWrapper) QueryRow(query string, ifs ...interface{}) *sql.Row {
	if t.tx != nil {
		return t.tx.QueryRow(query, ifs...)
	}
	return t.db.QueryRow(query, ifs...)
}

func (t *txWrapper) Query(query string, ifs ...interface{}) (*sql.Rows, error) {
	if t.tx != nil {
		return t.tx.Query(query, ifs...)
	}
	return t.db.Query(query, ifs...)
}
