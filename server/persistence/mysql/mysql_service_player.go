package mysql

import (
	"context"
	"encoding/json"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type dbPlayer struct {
	id      model.PlayerID
	name    string
	gIDJson []byte
}

func playerToDBPlayer(p model.Player) (dbPlayer, error) {
	res := dbPlayer{
		id:   p.ID,
		name: p.Name,
	}
	if len(p.Games) > 0 {
		gIDs := make([]model.GameID, 0, len(p.Games))
		for gID := range p.Games {
			gIDs = append(gIDs, gID)
		}
		j, err := json.Marshal(gIDs)
		if err != nil {
			return dbPlayer{}, err
		}
		res.gIDJson = j
	}
	return res, nil
}

func dbPlayerToPlayer(p dbPlayer, gs gameService) (model.Player, error) {
	res := model.Player{
		ID:   p.id,
		Name: p.name,
	}
	var gameIDs []model.GameID
	if len(p.gIDJson) > 0 {
		err := json.Unmarshal(p.gIDJson, &gameIDs)
		if err != nil {
			return model.Player{}, err
		}
		gameMap := make(map[model.GameID]model.PlayerColor, len(gameIDs))
		for _, id := range gameIDs {
			g, err := gs.Get(id)
			if err != nil {
				return model.Player{}, err
			}
			gameMap[id] = g.PlayerColors[res.ID]
		}
		res.Games = gameMap
	}
	return res, nil
}

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	gs  gameService
	ctx context.Context
	db  *sql.DB
}

func getPlayerService(gs gameService, ctx context.Context, db *sql.DB) (*playerService, error) {
	// TODO create the player table in the db if it doesn't exist
	return &playerService{
		gs:  gs,
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
	dbp, err := playerToDBPlayer(p)
	if err != nil {
		return err
	}

	tx, err := ps.beginTx()
	if err != nil {
		return err
	}
	// according to https://golang.org/pkg/database/sql/#Tx.ExecContext:
	// The rollback will be ignored if the tx has been committed later in the function.
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ps.ctx, `INSERT INTO `+playerTableName+` VALUES ( ?, ?, ? )`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ps.ctx, dbp.id, dbp.name, dbp.gIDJson)

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

	dbp := dbPlayer{}
	err = r.Scan(&dbp.id, &dbp.name, &dbp.gIDJson)
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
	p, err := dbPlayerToPlayer(dbp, ps.gs)
	if err != nil {
		return model.Player{}, err
	}
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

	dbp, err := playerToDBPlayer(p)
	if err != nil {
		return err
	}

	tx, err := ps.beginTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ps.ctx, `UPDATE `+playerTableName+` SET gameIDs = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ps.ctx, dbp.gIDJson, dbp.id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
