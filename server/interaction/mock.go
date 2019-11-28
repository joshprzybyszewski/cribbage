// +build !prod

package interaction

import (
	"github.com/stretchr/testify/mock"

	"github.com/joshprzybyszewski/cribbage/model"
)

var _ Player = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func (m *Mock) ID() model.PlayerID {
	args := m.Called()
	return args.Get(0).(model.PlayerID)
}
func (m *Mock) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	args := m.Called(b, g, s)
	return args.Error(0)
}
func (m *Mock) NotifyMessage(g model.Game, s string) error {
	args := m.Called(g, s)
	return args.Error(0)
}
func (m *Mock) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	args := m.Called(g, msgs)
	return args.Error(0)
}
