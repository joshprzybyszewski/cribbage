package memory

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DB = (*memDB)(nil)

type memDB struct {
	persistence.ServicesWrapper

	gs *gameService
	ps *playerService
	is *interactionService

	lock sync.Mutex
}

func New() persistence.DB {
	gs := getGameService()
	ps := getPlayerService()
	is := getInteractionService()
	sw := persistence.NewServicesWrapper(
		gs, ps, is,
	)

	return &memDB{
		ServicesWrapper: sw,
		gs:              gs.(*gameService),
		ps:              ps.(*playerService),
		is:              is.(*interactionService),
	}
}

func (mdb *memDB) Start() error {
	mdb.lock.Lock()
	defer mdb.lock.Unlock()

	mdb.gs = mdb.gs.Copy()
	mdb.ps = mdb.ps.Copy()
	mdb.is = mdb.is.Copy()

	sw := persistence.NewServicesWrapper(
		mdb.gs, mdb.ps, mdb.is,
	)

	mdb.ServicesWrapper = sw
	return nil
}

func (mdb *memDB) Commit() error {
	mdb.lock.Lock()
	defer mdb.lock.Unlock()

	err := saveGameService(mdb.gs)
	if err != nil {
		return err
	}
	err = savePlayerService(mdb.ps)
	if err != nil {
		return err
	}
	err = saveInteractionService(mdb.is)
	if err != nil {
		return err
	}
	mdb.ServicesWrapper = nil

	return nil
}

func (mdb *memDB) Rollback() error {
	mdb.lock.Lock()
	defer mdb.lock.Unlock()

	gs := getGameService()
	ps := getPlayerService()
	is := getInteractionService()
	sw := persistence.NewServicesWrapper(
		gs, ps, is,
	)
	mdb.ServicesWrapper = sw
	return nil
}
