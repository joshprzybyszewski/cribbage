package memory

import (
	"context"
	"sync"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DBFactory = memDBF{}

type memDBF struct {
	db *memDB
}

func NewFactory() (persistence.DBFactory, error) {
	return memDBF{}, nil
}

func (dbf memDBF) Close() error {
	dbf.db = nil

	return nil
}

var _ persistence.DB = (*memDB)(nil)

type memDB struct {
	persistence.ServicesWrapper

	lock sync.Mutex
}

func (dbf memDBF) New(context.Context) (persistence.DB, error) {
	if dbf.db != nil {
		return dbf.db, nil
	}

	sw := persistence.NewServicesWrapper(
		getGameService(),
		getPlayerService(),
		getInteractionService(),
	)

	dbf.db = &memDB{
		ServicesWrapper: sw,
	}
	return dbf.db, nil
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
