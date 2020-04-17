package memory

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var iservice *interactionService
var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	lock sync.Mutex

	interactions map[model.PlayerID]interaction.PlayerMeans
}

func getInteractionService() persistence.InteractionService {
	if iservice == nil {
		iservice = &interactionService{
			interactions: map[model.PlayerID]interaction.PlayerMeans{},
		}
	}
	return iservice
}

func clearInteractionService() {
	iservice.interactions = map[model.PlayerID]interaction.PlayerMeans{}
}

func (is *interactionService) Get(id model.PlayerID) (interaction.PlayerMeans, error) {
	is.lock.Lock()
	defer is.lock.Unlock()

	if i, ok := is.interactions[id]; ok {
		return i, nil
	}
	return interaction.PlayerMeans{}, persistence.ErrInteractionNotFound
}

func (is *interactionService) Create(pm interaction.PlayerMeans) error {
	is.lock.Lock()
	defer is.lock.Unlock()

	pID := pm.PlayerID
	if _, ok := is.interactions[pID]; ok {
		return persistence.ErrInteractionAlreadyExists
	}

	is.interactions[pID] = pm
	return nil
}

func (is *interactionService) Update(pm interaction.PlayerMeans) error {
	is.lock.Lock()
	defer is.lock.Unlock()

	is.interactions[pm.PlayerID] = pm
	return nil
}
