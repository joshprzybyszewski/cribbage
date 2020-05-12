package mysql

import (
	"context"
	"database/sql"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	// Players stores info about Players that we need to keep.
	// The default PreferredInteractionMode should be equal to int(interaction.UnsetMode)
	createPlayersTable = `CREATE TABLE IF NOT EXISTS Players (
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Name VARCHAR(` + maxPlayerNameLenStr + `),
		PreferredInteractionMode INT DEFAULT 0,
		PRIMARY KEY (PlayerID)
	) ENGINE = INNODB;`

	// GamePlayerColors keeps track of what color each player is in a given game
	// The default Color should match int(model.UnsetColor) for player colors
	createGamePlayerColorsTable = `CREATE TABLE IF NOT EXISTS GamePlayerColors (
		GameID INT UNSIGNED,
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Color TINYINT UNSIGNED DEFAULT 0,
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

	addPlayerToGame = `INSERT INTO GamePlayerColors
		(GameID, PlayerID)
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
	if len(p.ID) > maxPlayerUUIDLen {
		return persistence.ErrInvalidPlayerID
	}

	if len(p.Name) > maxPlayerNameLen {
		return persistence.ErrInvalidPlayerName
	}

	_, err := ps.db.Exec(createPlayer, p.ID, p.Name)
	err = convertMysqlError(err)
	if err != nil {
		if err == errDuplicateEntry {
			return persistence.ErrPlayerAlreadyExists
		}
		return err
	}

	return nil
}

func (ps *playerService) BeginGame(gID model.GameID, players []model.Player) error {
	for _, p := range players {
		if len(p.ID) > maxPlayerUUIDLen {
			return persistence.ErrInvalidPlayerID
		}

		_, err := ps.db.Exec(addPlayerToGame, gID, p.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *playerService) UpdateGameColor(pID model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	_, err := ps.db.Exec(updatePlayerColor, color, pID, gID)
	if err != nil {
		return err
	}

	return nil
}
