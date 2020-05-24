package mysql

import (
	"context"
	"database/sql"
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

	is *interactionService
	gs *gameService
	ps *playerService
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

	dbWrapper := txWrapper{
		db: db,
	}

	gs, err := getGameService(ctx, &dbWrapper)
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, &dbWrapper)
	if err != nil {
		return nil, err
	}
	is, err := getInteractionService(ctx, &dbWrapper)
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
		txWrapper:       &dbWrapper,
		ctx:             ctx,
		gs:              gs,
		ps:              ps,
		is:              is,
	}

	return &mw, nil
}

func (mw *mysqlWrapper) Close() error {
	return mw.txWrapper.db.Close()
}

func (mw *mysqlWrapper) Start() error {
	return mw.txWrapper.start(mw.ctx)
}

func (mw *mysqlWrapper) Commit() error {
	// we don't expect this to be called if Start() was never called
	return mw.txWrapper.commit()
}

func (mw *mysqlWrapper) Rollback() error {
	// we don't expect this to be called if Start() was never called
	return mw.txWrapper.rollback()
}

func (mw *mysqlWrapper) Clone() persistence.DB {
	cpy := *mw

	cpy.txWrapper = &txWrapper{
		db: mw.txWrapper.db,
	}
	cpy.is.db = cpy.txWrapper
	cpy.gs.db = cpy.txWrapper
	cpy.ps.db = cpy.txWrapper

	return &cpy
}
