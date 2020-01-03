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
	tx, err := ps.db.BeginTx(ps.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
	return model.Player{}, nil
}

func (ps *playerService) UpdateGameColor(id model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	return nil
}
