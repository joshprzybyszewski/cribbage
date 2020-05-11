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
}

func New(ctx context.Context, config Config) (persistence.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/cribbage`,
		config.DSNUser,
		config.DSNPassword,
		config.DSNHost,
		config.DSNPort,
	)
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

	return persistence.New(
		gs,
		ps,
		is,
	), nil
}
