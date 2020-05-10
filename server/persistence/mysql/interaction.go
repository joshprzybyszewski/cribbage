package mysql

import (
	"context"
	"database/sql"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	createInteractionTable = `CREATE TABLE IF NOT EXISTS Interactions (
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Mode INT(1),
		Means BLOB,
		PRIMARY KEY (PlayerID)
	) ENGINE = INNODB;`
)

var (
	interactionCreateStmts = []string{
		createInteractionTable,
	}
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	db *sql.DB
}

func getInteractionService(
	ctx context.Context,
	db *sql.DB,
) (persistence.InteractionService, error) {

	for _, createStmt := range interactionCreateStmts {
		_, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
	}

	return &interactionService{
		db: db,
	}, nil
}

func (s *interactionService) Get(id model.PlayerID) (interaction.PlayerMeans, error) {
	result := interaction.PlayerMeans{}
	// TODO get the means from the DB

	return result, nil
}

func (s *interactionService) Create(pm interaction.PlayerMeans) error {
	// TODO insert the means into the DB
	return nil
}

func (s *interactionService) Update(pm interaction.PlayerMeans) error {
	// TODO replace or add to the existing means
	return nil
}
