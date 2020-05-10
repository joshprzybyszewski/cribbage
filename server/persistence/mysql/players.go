package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	// Players stores info about Players that we need to keep.
	createPlayersTable = `CREATE TABLE IF NOT EXISTS Players (
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Name VARCHAR(255),
		PreferredInteractionMode INT(1),
		PRIMARY KEY (PlayerID)
	) ENGINE = INNODB;`

	createPlayersGameColorsTable = `CREATE TABLE IF NOT EXISTS PlayersGameColors (
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		GameID INT(1) UNSIGNED,
		Color TINYINT(1) UNSIGNED,
		PRIMARY KEY (PlayerID, GameID)
	) ENGINE = INNODB;`

	getPlayerName = `SELECT 
		Name
	FROM Players
		WHERE PlayerID = ? 
	;`

	getPlayerGames = `SELECT 
		GameID, Color
	FROM PlayersGameColors
		WHERE PlayerID = ? 
	;`

	createPlayer = `INSERT INTO Players
		(PlayerID, Name)
	VALUES
		(?, ?)
	;`

	// TODO consider inserting a "not set" value for the color
	addPlayerToGame = `INSERT INTO PlayersGameColors
		(PlayerID, GameID)
	VALUES
		(?, ?)
	;`

	updatePlayerColor = `UPDATE PlayersGameColors
	SET
		Color = ?
	WHERE
		PlayerID = ? AND
		GameID = ?
	;`
)

var (
	playersCreateStmts = []string{
		createPlayersTable,
		createPlayersGameColorsTable,
	}
)

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	db *sql.DB
}

func getPlayerService(
	ctx context.Context,
	db *sql.DB,
) (persistence.PlayerService, error) {

	for _, createStmt := range playersCreateStmts {
		_, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
	}

	return &playerService{
		db: db,
	}, nil
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	r := ps.db.QueryRow(getPlayerName, id)
	var name string
	err := r.Scan(
		&name,
	)
	if err != nil {
		return model.Player{}, err
	}

	rows, err := ps.db.Query(getPlayerGames, id)
	if err != nil {
		return model.Player{}, err
	}

	games := map[model.GameID]model.PlayerColor{}

	for rows.Next() {
		var gameID uint32
		var color uint8
		err = rows.Scan(&gameID, &color)
		if err != nil {
			return model.Player{}, err
		}

		games[model.GameID(gameID)] = model.PlayerColor(color)
	}

	if err := rows.Err(); err != nil {
		return model.Player{}, err
	}

	return model.Player{
		ID:    id,
		Name:  name,
		Games: games,
	}, nil
}

func (ps *playerService) Create(p model.Player) error {
	res, err := ps.db.Exec(createPlayer, p.ID, p.Name)
	if err != nil {
		return err
	}
	fmt.Printf("playerService.Create res: %+v\n", res)

	return nil
}

func (ps *playerService) UpdateGameColor(pID model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	res, err := ps.db.Exec(updatePlayerColor, color, pID, gID)
	if err != nil {
		return err
	}
	fmt.Printf("UpdateGameColor res: %+v\n", res)

	return nil
}
