package mysql

import (
	"context"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	ctx context.Context
	db  *sql.DB
}

func getPlayerService(ctx context.Context, db *sql.DB) (*playerService, error) {
	return &playerService{
		ctx: ctx,
		db:  db,
	}, nil
}

func (ps *playerService) Create(p model.Player) error {
	tx, err := ps.beginTx()
	if err != nil {
		return err
	}
	// according to https://golang.org/pkg/database/sql/#Tx.ExecContext:
	// The rollback will be ignored if the tx has been committed later in the function.
	defer tx.Rollback()

	// TODO how to use a constant for the table name?
	stmt, err := tx.PrepareContext(ps.ctx, `INSERT INTO players VALUES ( ?, ? )`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ps.ctx, p.ID, p.Name)

	if err != nil {
		if me, ok := err.(*mysql.MySQLError); !ok {
			return err
		} else if me.Number == ErrNumDuplicateEntry {
			return persistence.ErrPlayerAlreadyExists
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	result := model.Player{}
	tx, err := ps.beginTx()
	if err != nil {
		return model.Player{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ps.ctx, `SELECT * FROM players WHERE id=?`)
	if err != nil {
		return model.Player{}, err
	}
	defer stmt.Close()

	r := stmt.QueryRowContext(ps.ctx, id)
	err = r.Scan(&result.ID, &result.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			return model.Player{}, err
		}
		return model.Player{}, persistence.ErrPlayerNotFound
	}

	return result, nil
}

func (ps *playerService) UpdateGameColor(id model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	return nil
}

func (ps *playerService) beginTx() (*sql.Tx, error) {
	// https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels
	// let's use the highest level of isolation for now (which apparently uses read and write locks)
	return ps.db.BeginTx(ps.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
}
