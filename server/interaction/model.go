package interaction

import (
	"errors"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	UnsetMode Mode = 0
	Localhost Mode = 1
	NPC       Mode = 2
	Unknown   Mode = 3
)

type Mode int

type Means struct {
	Mode Mode        `protobuf:"-" json:"-" bson:"mode"`
	Info interface{} `protobuf:"-" json:"-" bson:"info"`
}

func (m *Means) AddSerializedInfo(serInfo []byte) error {
	switch m.Mode {
	case UnsetMode, Unknown:
		// nothing we know how to do for these
		return nil
	case Localhost:
		// the local host player expects a string as the info to tell us which port to connect to
		m.Info = string(serInfo)
		return nil
	case NPC:
		// serInfo should represent an action handler for the NPC.
		// It should be overwritten elsewhere to npcActionHandler
		return nil
	default:
		return fmt.Errorf(`unsupported Mode: %v`, m.Mode)

	}
}

func (m *Means) GetSerializedInfo() ([]byte, error) {
	switch m.Mode {
	case UnsetMode, Unknown:
		// nothing we know how to do here either
		return nil, nil
	case Localhost:
		str, ok := m.Info.(string)
		if !ok {
			return nil, errors.New(`localhost player should have a string as its info`)
		}
		return []byte(str), nil
	case NPC:
		// Info should represent an ActionHandler for the NPC.
		// It should be a pointer to a struct that implements this interface
		// so we can't serialize it.
		return nil, nil
	default:
		return nil, fmt.Errorf(`unsupported Mode: %v`, m.Mode)
	}
}

type PlayerMeans struct {
	PlayerID      model.PlayerID `protobuf:"-" json:"-" bson:"playerID"`
	PreferredMode Mode           `protobuf:"-" json:"-" bson:"preferredMode"`
	Interactions  []Means        `protobuf:"-" json:"-" bson:"interactions"`
}

func (pm PlayerMeans) getMeans(m Mode) Means {
	for _, i := range pm.Interactions {
		if i.Mode == m {
			return i
		}
	}

	return Means{}
}
