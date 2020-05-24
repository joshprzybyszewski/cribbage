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

var _ persistence.DB = (*mysqlWrapper)(nil)

type mysqlWrapper struct {
	persistence.ServicesWrapper

	txWrapper *txWrapper

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

	txWrapper := txWrapper{
		db: db,
	}

	gs, err := getGameService(ctx, &txWrapper)
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, &txWrapper)
	if err != nil {
		return nil, err
	}
	is, err := getInteractionService(ctx, &txWrapper)
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
		txWrapper:       &txWrapper,
		ctx:             ctx,
	}

	return &mw, nil
}

func (mw *mysqlWrapper) Close() error {
	return mw.txWrapper.db.Close()
}

func (mw *mysqlWrapper) Start() error {
	if mw.txWrapper.tx != nil {
		return errors.New(`mysql transaction already started`)
	}

	tx, err := mw.txWrapper.db.BeginTx(mw.ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	mw.txWrapper.tx = tx
	return nil
}

func (mw *mysqlWrapper) Commit() error {
	if mw.txWrapper.tx == nil {
		return errors.New(`mysql transaction not started`)
	}

	return mw.txWrapper.tx.Commit()
}

func (mw *mysqlWrapper) Rollback() error {
	if mw.txWrapper.tx == nil {
		return errors.New(`mysql transaction not started`)
	}

	return mw.txWrapper.tx.Rollback()
}
