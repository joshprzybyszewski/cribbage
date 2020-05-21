package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint:golint
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

type Config struct {
	DSNUser     string
	DSNPassword string
	DSNHost     string
	DSNPort     int
	DSNParams   string

	DatabaseName string
}

type txOrDB struct {
	db *sql.DB
	tx *sql.Tx
}

func (t *txOrDB) Exec(query string, ifs []interface{}) error {
	if t.tx != nil {
		return t.tx.Exec(query, ifs...)
	}
	return t.db.Exec(query, ifs...)
}

var _ persistence.DB = (*mysqlWrapper)(nil)

type mysqlWrapper struct {
	persistence.ServicesWrapper

	dt txOrDB

	ctx context.Context
}

func New(ctx context.Context, config Config) (persistence.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)`,
		config.DSNUser,
		config.DSNPassword,
		config.DSNHost,
		config.DSNPort,
	)
	if len(config.DatabaseName) > 0 {
		dsn += `/` + config.DatabaseName
	}
	if len(config.DSNParams) > 0 {
		dsn += `?` + config.DSNParams
	}
	db, err := sql.Open(`mysql`, dsn)
	if err != nil {
		return nil, err
	}

	gs, err := getGameService(ctx, db)
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, db)
	if err != nil {
		return nil, err
	}
	is, err := getInteractionService(ctx, db)
	if err != nil {
		return nil, err
	}

	sw := persistence.NewServicesWrapper(
		gs,
		ps,
		is,
	)

	mw := mysqlWrapper{
		ServicesWrapper: sw,
		dt:              txOrDB{db: db},
		ctx:             ctx,
	}

	return &mw, nil
}

func (mw *mysqlWrapper) Close() error {
	return mw.dt.db.Close()
}

func (mw *mysqlWrapper) Start() error {
	if mw.dt.tx != nil {
		return errors.New(`mysql transaction already started`)
	}

	tx, err := mw.db.BeginTx(mw.ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	mw.dt.tx = tx
	return nil
}

func (mw *mysqlWrapper) Commit() error {
	if mw.dt.tx == nil {
		return errors.New(`mysql transaction not started`)
	}

	return mw.dt.tx.Commit()
}

func (mw *mysqlWrapper) Rollback() error {
	if mw.dt.tx == nil {
		return errors.New(`mysql transaction not started`)
	}

	return mw.dt.tx.Rollback()
}
