package memory

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DB = (*memDB)(nil)

type memDB struct {
	persistence.ServicesWrapper

	lock sync.Mutex
}

func New() persistence.DB {
	sw := persistence.NewServicesWrapper(
		getGameService(),
		getPlayerService(),
		getInteractionService(),
	)

	return &memDB{
		ServicesWrapper: sw,
	}
}

func (mdb *memDB) Start() error {
	mdb.lock.Lock()
	defer mdb.lock.Unlock()

	return nil
}

func (mdb *memDB) Commit() error {
	mdb.lock.Lock()
	defer mdb.lock.Unlock()

	return nil
}

func (mdb *memDB) Rollback() error {
	mdb.lock.Lock()
	defer mdb.lock.Unlock()

	return nil
}
