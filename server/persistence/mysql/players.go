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
	// The default PreferredInteractionMode should be equal to int(interaction.UnsetMode)
	createPlayersTable = `CREATE TABLE IF NOT EXISTS Players (
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Name VARCHAR(255),
		PreferredInteractionMode INT(1) DEFAULT 0,
		PRIMARY KEY (PlayerID)
	) ENGINE = INNODB;`

	// GamePlayerColors keeps track of what color each player is in a given game
	// The default Color should match int(model.UnsetColor) for player colors
	createGamePlayerColorsTable = `CREATE TABLE IF NOT EXISTS GamePlayerColors (
		GameID INT(1) UNSIGNED,
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Color TINYINT(1) UNSIGNED DEFAULT 0,
		PRIMARY KEY (GameID, PlayerID)
	) ENGINE = INNODB;`

	getPlayerName = `SELECT 
		Name
	FROM Players
	WHERE PlayerID = ? 
	;`

	getPlayerGames = `SELECT 
		GameID, Color
	FROM GamePlayerColors
	WHERE PlayerID = ? 
	;`

	getPlayerColorsForGame = `SELECT 
		PlayerID, Color
	FROM GamePlayerColors
	WHERE GameID = ?
	;`

	createPlayer = `INSERT INTO Players
		(PlayerID, Name)
	VALUES
		(?, ?)
	;`

	// TODO consider inserting a "not set" value for the color
	addPlayerToGame = `INSERT INTO GamePlayerColors
		(PlayerID, GameID)
	VALUES
		(?, ?)
	;`

	updatePlayerColor = `UPDATE GamePlayerColors
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
		createGamePlayerColorsTable,
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
		if err == sql.ErrNoRows {
			return model.Player{}, persistence.ErrPlayerNotFound
		}
		return model.Player{}, err
	}

	rows, err := ps.db.Query(getPlayerGames, id)
	if err != nil {
		return model.Player{}, err
	}

	games := map[model.GameID]model.PlayerColor{}

	for rows.Next() {
		var gameID model.GameID
		var color model.PlayerColor
		err = rows.Scan(&gameID, &color)
		if err != nil {
			return model.Player{}, err
		}

		games[gameID] = color
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
