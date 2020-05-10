package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	createGameTable = `CREATE TABLE IF NOT EXISTS Games (
		GameID INT(1) UNSIGNED,
		ScoreBlue TINYINT(1) UNSIGNED,
		ScoreRed TINYINT(1) UNSIGNED,
		ScoreGreen TINYINT(1) UNSIGNED,
		ScoreBlueLag TINYINT(1) UNSIGNED,
		ScoreRedLag TINYINT(1) UNSIGNED,
		ScoreGreenLag TINYINT(1) UNSIGNED,
		Phase TINYINT(1) UNSIGNED,
		BlockingPlayers SMALLBLOB, -- json encoded map of who's blocking and why
		CurrentDealer VARCHAR(` + maxPlayerUUIDLenStr + `),
		Hands SMALLBLOB, -- json encoded map of slices for player hands
		Crib TINYINT(4), -- number representation of the cards in the crib
		CutCard TINYINT(1), -- number representation of the card that's been cut
		PeggedCards SMALLBLOB, -- the json-encoded slice of previously pegged cards
		NumActions INT(1) UNSIGNED, -- how many actions have occurred in the game before this one
		Action SMALLBLOB, -- the json encoded PlayerAction
		PRIMARY KEY (GameID, NumActions)
	) ENGINE = INNODB;`

	createGamePlayersTable = `CREATE TABLE IF NOT EXISTS GamePlayers (
		GameID INT(1) UNSIGNED,
		Player1ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Player2ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Player3ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Player4ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		PRIMARY KEY (GameID)
	) ENGINE = INNODB;`

	queryLatestGame = `SELECT 
		g.GameID,
		gp.Player1ID, gp.Player2ID, gp.Player3ID, gp.Player4ID,
		g.ScoreBlue, g.ScoreRed, g.ScoreGreen,
		g.ScoreBlueLag, g.ScoreRedLag, g.ScoreGreenLag,
		g.Phase, g.BlockingPlayers, g.CurrentDealer,
		g.Hands, g.Crib, g.CutCard,
		g.PeggedCards,
		g.NumActions, g.Action
	FROM Games g
	INNER JOIN GamePlayers gp
		ON g.GameID = gp.GameID
		WHERE g.GameID = ? 
	SORT DESC NumActions
	LIMIT 1;`

	queryGameAtNumActions = `SELECT 
		g.GameID,
		gp.Player1ID, gp.Player2ID, gp.Player3ID, gp.Player4ID,
		g.ScoreBlue, g.ScoreRed, g.ScoreGreen,
		g.ScoreBlueLag, g.ScoreRedLag, g.ScoreGreenLag,
		g.Phase, g.BlockingPlayers, g.CurrentDealer,
		g.Hands, g.Crib, g.CutCard,
		g.PeggedCards,
		g.NumActions, g.Action
	FROM Games g
	INNER JOIN GamePlayers gp
		ON g.GameID = gp.GameID
		WHERE g.GameID = ?
		WHERE g.NumActions = ?
	;`
)

var (
	createStmts = []string{
		createGameTable,
		createGamePlayersTable,
	}
)

var _ persistence.GameService = (*gameService)(nil)

type gameService struct {
	db *sql.DB
}

func getGameService(
	ctx context.Context,
	db *sql.DB,
) (persistence.GameService, error) {

	for _, createStmt := range createStmts {
		res, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
		fmt.Printf("res: %+v\n", res)
	}

	return &gameService{
		db: db,
	}, nil
}

func (g *gameService) Get(id model.GameID) (model.Game, error) {
	r := g.db.QueryRow(queryLatestGame, id)
	return g.populateGameFromRow(r)
}

func (g *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	r := g.db.QueryRow(queryGameAtNumActions, id, numActions)
	return g.populateGameFromRow(r)
}

func (g *gameService) populateGameFromRow(r *sql.Row) (model.Game, error) {
	var gameID uint32
	var p1ID, p2ID, p3ID, p4ID string
	var scoreBlue, scoreRed, scoreGreen, lagScoreBlue, lagScoreRed, lagScoreGreen int
	var phase uint8
	var blockingPlayers []byte
	var curDealerID string
	var hands []byte
	var cribCards []int8 = make([]int8, 4)
	var cutCard int8
	var peggedCards []byte
	var numActions uint32
	var action []byte
	err := r.Scan(
		&gameID,
		&p1ID, &p2ID, &p3ID, &p4ID,
		&scoreBlue, &scoreRed, &scoreGreen,
		&lagScoreBlue, &lagScoreRed, &lagScoreGreen,
		&phase, &blockingPlayers, &curDealerID,
		&hands, &cribCards, &cutCard,
		&peggedCards,
		&numActions, &action,
	)
	if err != nil {
		return model.Game{}, err
	}

	// TODO choose what colors based on the colors defined by the players table

	curScores := make(map[model.PlayerColor]int, 3)
	lagScores := make(map[model.PlayerColor]int, 3)
	if scoreBlue > 0 {
		curScores[model.Blue] = scoreBlue
		lagScores[model.Blue] = lagScoreBlue
	}
	if scoreRed > 0 {
		curScores[model.Red] = scoreRed
		lagScores[model.Red] = lagScoreRed
	}
	if scoreGreen > 0 {
		curScores[model.Green] = scoreGreen
		lagScores[model.Green] = lagScoreGreen
	}

	game := model.Game{
		ID:            model.GameID(gameID),
		Players:       nil, // []model.Player
		PlayerColors:  nil, //map[model.PlayerID]model.PlayerColor
		CurrentScores: curScores,
		LagScores:     lagScores,
		Phase:         model.Phase(phase),
		//BlockingPlayers map[PlayerID]Blocker
		//CurrentDealer PlayerID
		//Hands map[PlayerID][]Card
		// Crib []Card
		// CutCard Card
		// PeggedCards []PeggedCard
		// Actions []PlayerAction
	}

	return game, nil
}

func (g *gameService) UpdatePlayerColor(id model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	// TODO persist the updated player color
	return nil
}

func (g *gameService) Save(mg model.Game) error {
	// TODO persist the game in the DB
	return nil
}
