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
	// TODO create the player table in the db if it doesn't exist
	return &playerService{
		ctx: ctx,
		db:  db,
	}, nil
}

func (ps *playerService) beginTx() (*sql.Tx, error) {
	// https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels
	// let's use the highest level of isolation for now (which apparently uses read and write locks)
	return ps.db.BeginTx(ps.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

func (ps *playerService) Create(p model.Player) error {
	tx, err := ps.beginTx()
	if err != nil {
		return err
	}
	// according to https://golang.org/pkg/database/sql/#Tx.ExecContext:
	// The rollback will be ignored if the tx has been committed later in the function.
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ps.ctx, `INSERT INTO `+playerTableName+` VALUES ( ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ps.ctx, p.ID, p.Name)

	if err != nil {
		if me, ok := err.(*mysql.MySQLError); !ok {
			return err
		} else if me.Number == sqlErrCodeDuplicateEntry {
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
	tx, err := ps.beginTx()
	if err != nil {
		return model.Player{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ps.ctx, `SELECT * FROM `+playerTableName+` WHERE id=?`)
	if err != nil {
		return model.Player{}, err
	}
	defer stmt.Close()

	r := stmt.QueryRowContext(ps.ctx, id)

	p := model.Player{}
	err = r.Scan(&p.ID, &p.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			return model.Player{}, err
		}
		return model.Player{}, persistence.ErrPlayerNotFound
	}
	// transaction is done now
	err = tx.Commit()
	if err != nil {
		return model.Player{}, err
	}
	// TODO query the games table to populate Games on p
	return p, nil
}

func (ps *playerService) UpdateGameColor(pID model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	p, err := ps.Get(pID)
	if err != nil {
		return err
	}

	if p.Games == nil {
		p.Games = make(map[model.GameID]model.PlayerColor, 1)
	}

	if c, ok := p.Games[gID]; ok {
		if c != color {
			return persistence.ErrPlayerColorMismatch
		}

		// Nothing to do; the player already knows its color
		return nil
	}
	p.Games[gID] = color
	// TODO update this game in the games table
	return nil
}
