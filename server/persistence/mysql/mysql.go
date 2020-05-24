package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint:golint
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DBFactory = (*mysqlDBFactory)(nil)

type mysqlDBFactory struct {
	db *sql.DB
}

func NewFactory(ctx context.Context, config Config) (persistence.DBFactory, error) {
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

	allCreateStmts := make([]string, 0, len(gamesCreateStmts)+len(playersCreateStmts)+len(interactionCreateStmts))
	allCreateStmts = append(allCreateStmts, gamesCreateStmts...)
	allCreateStmts = append(allCreateStmts, playersCreateStmts...)
	allCreateStmts = append(allCreateStmts, interactionCreateStmts...)

	for _, createStmt := range allCreateStmts {
		_, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
	}

	return &mysqlDBFactory{
		db: db,
	}, nil
}

func (dbf *mysqlDBFactory) New(ctx context.Context) (persistence.DB, error) {
	dbWrapper := txWrapper{
		db: dbf.db,
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
	}

	return &mw, nil
}

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
