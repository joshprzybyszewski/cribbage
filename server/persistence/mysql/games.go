package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
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
	// Crib is a 4-byte int of the (up to 4) cards in the crib where every byte is each crib card.
	//   If this weren't just a fun project, I wouldn't try to be this tricky.
	// CurrentDealer is the PlayerID for the dealer
	// BlockingPlayers is a json encoded map of who's blocking and why
	// Hands is a json encoded map of slices for player hands
	// PeggedCards is the json-encoded slice of previously pegged cards
	// Action is the json encoded model.PlayerAction
	createGameTable = `CREATE TABLE IF NOT EXISTS Games (
		GameID INT UNSIGNED,
		NumActions INT UNSIGNED,
		ScoreBlue TINYINT UNSIGNED,
		ScoreRed TINYINT UNSIGNED,
		ScoreGreen TINYINT UNSIGNED,
		ScoreBlueLag TINYINT UNSIGNED,
		ScoreRedLag TINYINT UNSIGNED,
		ScoreGreenLag TINYINT UNSIGNED,
		Phase TINYINT UNSIGNED,
		CutCard SMALLINT,
		Crib INT,
		CurrentDealer VARCHAR(` + maxPlayerUUIDLenStr + `),
		BlockingPlayers BLOB,
		Hands BLOB,
		PeggedCards BLOB,
		Action BLOB,
		PRIMARY KEY (GameID, NumActions)
	) ENGINE = INNODB;`

	createGamePlayersTable = `CREATE TABLE IF NOT EXISTS GamePlayers (
		GameID INT UNSIGNED,
		Player1ID VARCHAR(` + maxPlayerUUIDLenStr + `) NOT NULL,
		Player2ID VARCHAR(` + maxPlayerUUIDLenStr + `) NOT NULL,
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
	ORDER BY
		NumActions DESC
	LIMIT 1;`

	queryGameAtNumActions = `SELECT 
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
	WHERE g.GameID = ? AND
		g.NumActions = ?
	;`

	queryPlayerActionsBefore = `SELECT 
		NumActions, Action
	FROM Games
	WHERE GameID = ? AND
		NumActions <= ?
	;`

	addPlayersToGamePlayers = `INSERT INTO GamePlayers
		(
			GameID, 
			Player1ID, Player2ID, Player3ID, Player4ID
		)
	VALUES
		(
			?,
			?, ?, ?, ?
		)
	;`

	insertGameAt = `INSERT INTO Games
		(
			GameID, NumActions, 
			ScoreBlue, ScoreRed, ScoreGreen,
			ScoreBlueLag, ScoreRedLag, ScoreGreenLag,
			Phase, CutCard, Crib,
			CurrentDealer,
			BlockingPlayers, Hands, PeggedCards, Action
		)
	VALUES
		(
			?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?,
			?, ?, ?, ?
		)
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
	db *txWrapper
}

func getGameService(
	db *txWrapper,
) persistence.GameService {

	return &gameService{
		db: db,
	}
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

	var p1ID, p2ID model.PlayerID
	var p3ID, p4ID *model.PlayerID
	var curDealerID model.PlayerID
	var scoreBlue, scoreRed, scoreGreen,
		lagScoreBlue, lagScoreRed, lagScoreGreen uint8
	var phase model.Phase
	var cribCardInts int32
	var cutCardInt int8
	var blockingPlayers, hands, peggedCards, action []byte
	var numActions uint32
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

	players, err := g.getPlayersForGame(p1ID, p2ID, p3ID, p4ID)
	if err != nil {
		return model.Game{}, err
	}

	pc, err := g.getPlayerColors(gID)
	if err != nil {
		return model.Game{}, err
	}
	addInPopulatedColor(curScores, lagScores, pc)

	cutCard, err := model.NewCardFromTinyInt(cutCardInt)
	if err != nil {
		// We interpret an error here to mean that there is no cut
		// card. Therefore, we set it to the empty card.
		cutCard = model.Card{}
	}

	cribCards := getCribCards(cribCardInts)

	bp, err := getBlockingPlayers(blockingPlayers)
	if err != nil {
		return model.Game{}, err
	}

	h, err := getHands(hands)
	if err != nil {
		return model.Game{}, err
	}

	p, err := getPeggedCards(peggedCards)
	if err != nil {
		return model.Game{}, err
	}

	pas, err := g.getActions(gID, int(numActions))
	if err != nil {
		return model.Game{}, err
	}

	game := model.Game{
		ID:              gID,
		CurrentScores:   curScores,
		LagScores:       lagScores,
		Players:         players,
		PlayerColors:    pc,
		Phase:           phase,
		CurrentDealer:   curDealerID,
		CutCard:         cutCard,
		Crib:            cribCards,
		BlockingPlayers: bp,
		Hands:           h,
		PeggedCards:     p,
		Actions:         pas,
	}

	return game, nil
}

func populateScores(
	scoreBlue, scoreRed, scoreGreen,
	lagScoreBlue, lagScoreRed, lagScoreGreen uint8,
) (cur, lag map[model.PlayerColor]int) {
	curScores := make(map[model.PlayerColor]int, 3)
	lagScores := make(map[model.PlayerColor]int, 3)
	if scoreBlue > 0 {
		curScores[model.Blue] = int(scoreBlue)
		lagScores[model.Blue] = int(lagScoreBlue)
	}
	if scoreRed > 0 {
		curScores[model.Red] = int(scoreRed)
		lagScores[model.Red] = int(lagScoreRed)
	}
	if scoreGreen > 0 {
		curScores[model.Green] = int(scoreGreen)
		lagScores[model.Green] = int(lagScoreGreen)
	}

	return curScores, lagScores
}

func addInPopulatedColor(
	curScores, lagScores map[model.PlayerColor]int,
	pc map[model.PlayerID]model.PlayerColor,
) {
	// if we know what color the players are, but we don't have point entries
	// for those colors in the scores maps, add zeros
	for _, color := range pc {
		if _, ok := curScores[color]; !ok {
			curScores[color] = 0
		}
		if _, ok := lagScores[color]; !ok {
			lagScores[color] = 0
		}
	}
}

func getCribCards(cribCardInt int32) []model.Card {
	var cribCards []model.Card
	var cci int8
	for i := uint(0); i < 4; i++ {
		cci = int8(cribCardInt >> (8 * i))
		c, err := model.NewCardFromTinyInt(cci)
		if err != nil {
			// If we've errored here, we assume it just means the card isn't set
			continue
		}
		cribCards = append(cribCards, c)
	}
	return cribCards
}

func serializeCribCards(crib []model.Card) int32 {
	val := int32(0)
	for i := uint(0); i < 4; i++ {
		ti := int8(model.NumCardsPerDeck + 1) // set it to an invalid num
		if int(i) < len(crib) {
			ti = crib[i].ToTinyInt()
		}
		val |= (int32(ti) << (8 * i))
	}
	return val
}

func getBlockingPlayers(ser []byte) (map[model.PlayerID]model.Blocker, error) {
	blockers := map[model.PlayerID]model.Blocker{}

	err := json.Unmarshal(ser, &blockers)
	if err != nil {
		return nil, err
	}

	return blockers, nil
}

func serializeBlockingPlayers(input map[model.PlayerID]model.Blocker) ([]byte, error) {
	return json.Marshal(input)
}

func getHands(ser []byte) (map[model.PlayerID][]model.Card, error) {
	hands := map[model.PlayerID][]model.Card{}

	err := json.Unmarshal(ser, &hands)
	if err != nil {
		return nil, err
	}

	return hands, nil
}

func serializeHands(input map[model.PlayerID][]model.Card) ([]byte, error) {
	return json.Marshal(input)
}

func getPeggedCards(ser []byte) ([]model.PeggedCard, error) {
	peggedCards := []model.PeggedCard{}

	err := json.Unmarshal(ser, &peggedCards)
	if err != nil {
		return nil, err
	}

	return peggedCards, nil
}

func serializePeggedCards(input []model.PeggedCard) ([]byte, error) {
	return json.Marshal(input)
}

func (g *gameService) getPlayersForGame(
	p1ID, p2ID model.PlayerID,
	p3ID, p4ID *model.PlayerID,
) ([]model.Player, error) {

	if len(p1ID) == 0 || len(p2ID) == 0 {
		return nil, errors.New(`at least two players required`)
	}

	pIDs := []model.PlayerID{
		p1ID, p2ID,
	}

	if p3ID != nil && len(*p3ID) > 0 {
		// The third and fourth players can only exist if the first two do
		pIDs = append(pIDs, *p3ID)
		if p4ID != nil && len(*p4ID) > 0 {
			pIDs = append(pIDs, *p4ID)
		}
	}

	players := make([]model.Player, len(pIDs))
	for i, pID := range pIDs {
		players[i].ID = pID
	}
	return players, nil

}

func (g *gameService) getPlayerColors(
	gID model.GameID,
) (map[model.PlayerID]model.PlayerColor, error) {

	// populate pc with the colors for each player
	pc := make(map[model.PlayerID]model.PlayerColor, 4)

	rows, err := g.db.Query(getPlayerColorsForGame, gID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pID model.PlayerID
		var color model.PlayerColor
		err := rows.Scan(&pID, &color)
		if err != nil {
			return nil, err
		}
		if color != model.UnsetColor {
			pc[pID] = color
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pc, nil
}

func (g *gameService) getActions(
	gID model.GameID,
	maxNumActions int,
) ([]model.PlayerAction, error) {

	rows, err := g.db.Query(queryPlayerActionsBefore, gID, maxNumActions)
	if err != nil {
		return nil, err
	}
	paMap := make(map[int][]byte, maxNumActions)
	var lenActionSlice, actionIndex int
	var serAction []byte
	for rows.Next() {
		err = rows.Scan(&lenActionSlice, &serAction)
		if err != nil {
			return nil, err
		}
		// we subtract one because the last action is serialized and paired
		// with the len of the action slice. Therefore, we need to say that
		// this action's index (into the action slice) is one fewer than the
		// number we persisted it at
		actionIndex = lenActionSlice - 1
		paMap[actionIndex] = serAction
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	pas := make([]model.PlayerAction, maxNumActions)
	for i := range pas {
		bytes, ok := paMap[i]
		if !ok {
			return nil, errors.New(`missing action`)
		}
		pa, err := getPlayerAction(bytes)
		if err != nil {
			return nil, err
		}
		pas[i] = pa
	}

	return pas, nil
}

func getPlayerAction(ser []byte) (model.PlayerAction, error) {
	return jsonutils.UnmarshalPlayerAction(ser)
}

func serializePlayerAction(input model.PlayerAction) ([]byte, error) {
	// Remember: this is complemented by jsonutils.UnmarshalPlayerAction
	// because we have to unmarshal into an interface
	return json.Marshal(input)
}

func (g *gameService) UpdatePlayerColor(id model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	// There should be nothing to do here because the player service should take care
	// of all of the persistence that needs to happen
	return nil
}

func (g *gameService) Begin(mg model.Game) error {
	ifs := []interface{}{
		mg.ID,
	}
	for _, p := range mg.Players {
		if len(p.ID) > maxPlayerUUIDLen {
			return persistence.ErrInvalidPlayerID
		}

		ifs = append(ifs, p.ID)
	}
	for len(ifs) < 5 {
		// the query expects 5 inputs. it'd be better to have variadic queries
		// but I don't want to write that right now.
		ifs = append(ifs, nil)
	}

	_, err := g.db.Exec(addPlayersToGamePlayers, ifs...)
	if err != nil {
		return err
	}

	return g.Save(mg)
}

func (g *gameService) Save(mg model.Game) error {
	if mg.ID > maxGameID {
		return persistence.ErrInvalidGameID
	}

	if len(mg.CurrentDealer) > maxPlayerUUIDLen {
		return persistence.ErrInvalidPlayerID
	}

	if err := persistence.ValidateLatestActionBelongs(mg); err != nil {
		return err
	}

	cut := mg.CutCard.ToTinyInt()
	crib := serializeCribCards(mg.Crib)

	bp, err := serializeBlockingPlayers(mg.BlockingPlayers)
	if err != nil {
		return err
	}
	h, err := serializeHands(mg.Hands)
	if err != nil {
		return err
	}
	pegged, err := serializePeggedCards(mg.PeggedCards)
	if err != nil {
		return err
	}
	var a []byte
	if ai := mg.NumActions() - 1; ai >= 0 {
		// get the last action in the slice of actions. Serialize it for saving
		a, err = serializePlayerAction(mg.Actions[ai])
		if err != nil {
			return err
		}
	}

	ifs := []interface{}{
		mg.ID, mg.NumActions(),
		uint8(mg.CurrentScores[model.Blue]), uint8(mg.CurrentScores[model.Red]), uint8(mg.CurrentScores[model.Green]),
		uint8(mg.LagScores[model.Blue]), uint8(mg.LagScores[model.Red]), uint8(mg.LagScores[model.Green]),
		mg.Phase, cut, crib,
		mg.CurrentDealer,
		bp, h, pegged, a,
	}
	_, err = g.db.Exec(insertGameAt, ifs...)
	if err != nil {
		return err
	}

	return nil
}
