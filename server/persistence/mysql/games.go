package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	// Games stores the state of a game at a given time.
	//   Each Action will update a games state and we keep a full history of all actions.
	// The columns act as follows:
	// GameID is a UUID to identify a game
	// NumActions is how many actions have occurred in the game before this one
	// ScoreBlue, ScoreRed, and ScoreGreen are the scores for each color
	// ScoreBlueLag, ScoreRedLag, and ScoreGreenLag are the previous scores for each color
	// Phase is the model.Phase that the game is currently in
	// CutCard is a number representation of the card that's been cut
	// Crib is a number representation of the (up to 4) cards in the crib
	// CurrentDealer is the PlayerID for the dealer
	// BlockingPlayers is a json encoded map of who's blocking and why
	// Hands is a json encoded map of slices for player hands
	// PeggedCards is the json-encoded slice of previously pegged cards
	// Action is the json encoded model.PlayerAction
	createGameTable = `CREATE TABLE IF NOT EXISTS Games (
		GameID INT(1) UNSIGNED,
		NumActions INT(1) UNSIGNED,
		ScoreBlue TINYINT(1) UNSIGNED,
		ScoreRed TINYINT(1) UNSIGNED,
		ScoreGreen TINYINT(1) UNSIGNED,
		ScoreBlueLag TINYINT(1) UNSIGNED,
		ScoreRedLag TINYINT(1) UNSIGNED,
		ScoreGreenLag TINYINT(1) UNSIGNED,
		Phase TINYINT(1) UNSIGNED,
		CutCard TINYINT(1) UNSIGNED,
		Crib TINYINT(4) UNSIGNED,
		CurrentDealer VARCHAR(` + maxPlayerUUIDLenStr + `),
		BlockingPlayers BLOB,
		Hands BLOB,
		PeggedCards BLOB,
		Action BLOB,
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
	gamesCreateStmts = []string{
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

	for _, createStmt := range gamesCreateStmts {
		_, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
	}

	return &gameService{
		db: db,
	}, nil
}

func (g *gameService) Get(id model.GameID) (model.Game, error) {
	r := g.db.QueryRow(queryLatestGame, id)
	return g.populateGameFromRow(id, r)
}

func (g *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	r := g.db.QueryRow(queryGameAtNumActions, id, numActions)
	return g.populateGameFromRow(id, r)
}

func (g *gameService) populateGameFromRow(
	gID model.GameID,
	r *sql.Row,
) (model.Game, error) {

	var p1ID, p2ID, p3ID, p4ID string
	var scoreBlue, scoreRed, scoreGreen int
	var lagScoreBlue, lagScoreRed, lagScoreGreen int
	var phase uint8
	var blockingPlayers []byte
	var curDealerID string
	var hands []byte
	var cribCardInts []int8 = make([]int8, 4)
	var cutCardInt int8
	var peggedCards []byte
	var numActions uint32
	var action []byte
	err := r.Scan(
		&p1ID, &p2ID, &p3ID, &p4ID,
		&scoreBlue, &scoreRed, &scoreGreen,
		&lagScoreBlue, &lagScoreRed, &lagScoreGreen,
		&phase, &blockingPlayers, &curDealerID,
		&hands, &cribCardInts, &cutCardInt,
		&peggedCards,
		&numActions, &action,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Game{}, persistence.ErrGameNotFound
		}
		return model.Game{}, err
	}

	curScores, lagScores := populateScores(
		scoreBlue, scoreRed, scoreGreen,
		lagScoreBlue, lagScoreRed, lagScoreGreen,
	)

	players, pc, err := g.getPlayerFields(gID, p1ID, p2ID, p3ID, p4ID)
	if err != nil {
		return model.Game{}, err
	}

	cutCard, err := model.NewCardFromTinyInt(cutCardInt)
	if err != nil {
		// If we've errored here, just ignore it and continue
		fmt.Printf("errored card for cut: %+v\n", err)
	}

	var cribCards []model.Card
	for _, cci := range cribCardInts {
		c, err := model.NewCardFromTinyInt(cci)
		if err != nil {
			// If we've errored here, just ignore it and continue
			fmt.Printf("errored card while building crib: %+v\n", err)
			continue
		}
		cribCards = append(cribCards, c)
	}

	game := model.Game{
		ID:            gID,
		CurrentScores: curScores,
		LagScores:     lagScores,
		Players:       players,
		PlayerColors:  pc,
		Phase:         model.Phase(phase),
		CurrentDealer: model.PlayerID(curDealerID),
		CutCard:       cutCard,
		Crib:          cribCards,
		//BlockingPlayers map[PlayerID]Blocker
		//Hands map[PlayerID][]Card
		// PeggedCards []PeggedCard
		// Actions []PlayerAction
	}

	return game, nil
}

func populateScores(
	scoreBlue, scoreRed, scoreGreen,
	lagScoreBlue, lagScoreRed, lagScoreGreen int,
) (cur, lag map[model.PlayerColor]int) {
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

	return curScores, lagScores
}

func (g *gameService) getPlayerFields(
	id model.GameID,
	p1ID, p2ID, p3ID, p4ID string,
) ([]model.Player, map[model.PlayerID]model.PlayerColor, error) {

	pIDs := make([]model.PlayerID, 0, 4)
	if len(p1ID) > 0 {
		pIDs = append(pIDs, model.PlayerID(p1ID))
	}
	if len(p2ID) > 0 {
		pIDs = append(pIDs, model.PlayerID(p2ID))
	}
	if len(p3ID) > 0 {
		pIDs = append(pIDs, model.PlayerID(p3ID))
	}
	if len(p4ID) > 0 {
		pIDs = append(pIDs, model.PlayerID(p4ID))
	}

	players := make([]model.Player, len(pIDs))
	pc := make(map[model.PlayerID]model.PlayerColor, len(pIDs))
	for i, pID := range pIDs {
		// TODO get the entire "player", not just the ID
		// TODO populate pc with the colors for each player
		players[i].ID = pID
	}

	return players, pc, nil
}

func (g *gameService) UpdatePlayerColor(id model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	// There should be nothing to do here because the player service should take care
	// of all of the persistence that needs to happen
	return nil
}

func (g *gameService) Save(mg model.Game) error {
	// TODO persist the game in the DB
	// TODO call addPlayerToGame for each player in this game
	return nil
}
