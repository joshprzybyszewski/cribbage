package memory

import (
	"context"
	"sync"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DBFactory = memDBF{}

type memDBF struct{}

func NewFactory() (persistence.DBFactory, error) {
	return memDBF{}, nil
}

var _ persistence.DB = (*memDB)(nil)

type memDB struct {
	persistence.ServicesWrapper

	lock sync.Mutex
}

func (memDBF) New(context.Context) (persistence.DB, error) {
	sw := persistence.NewServicesWrapper(
		getGameService(),
		getPlayerService(),
		getInteractionService(),
	)

	return &memDB{
		ServicesWrapper: sw,
	}, nil
}

func (mdb *memDB) Close() error {
	mdb.ServicesWrapper = nil

	return nil
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

func Clear() {
	gservice = nil
	pservice = nil
	iservice = nil
}
