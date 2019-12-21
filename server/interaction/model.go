package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	Unset Mode = iota
	Localhost
	Unknown
)

type Mode int

type Means struct {
	Mode Mode        `protobuf:"-" json:"-" bson:"mode"`
	Info interface{} `protobuf:"-" json:"-" bson:"info"`
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
