package mysql

import (
	"context"
	"database/sql"
	"errors"
	"sync"
)

type txWrapper struct {
	db *sql.DB
	tx *sql.Tx

	lock sync.Mutex
}

func (t *txWrapper) start(ctx context.Context) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.tx != nil {
		return errors.New(`mysql transaction already started`)
	}

	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	t.tx = tx

	return nil
}

func (t *txWrapper) commit() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	defer func() {
		t.tx = nil
	}()

	if t.tx == nil {
		return errors.New(`mysql transaction not started`)
	}

	return t.tx.Commit()
}

func (t *txWrapper) rollback() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	defer func() {
		t.tx = nil
	}()

	if t.tx == nil {
		return errors.New(`mysql transaction not started`)
	}

	return t.tx.Rollback()
}

func (t *txWrapper) Exec(query string, ifs ...interface{}) (sql.Result, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.tx != nil {
		return t.tx.Exec(query, ifs...)
	}
	return t.db.Exec(query, ifs...)
}

func (t *txWrapper) ExecContext(ctx context.Context, query string, ifs ...interface{}) (sql.Result, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.tx != nil {
		return t.tx.ExecContext(ctx, query, ifs...)
	}
	return t.db.ExecContext(ctx, query, ifs...)
}

func (t *txWrapper) QueryRow(query string, ifs ...interface{}) *sql.Row {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.tx != nil {
		return t.tx.QueryRow(query, ifs...)
	}
	return t.db.QueryRow(query, ifs...)
}

func (t *txWrapper) Query(query string, ifs ...interface{}) (*sql.Rows, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.tx != nil {
		return t.tx.Query(query, ifs...)
	}
	return t.db.Query(query, ifs...)
}
