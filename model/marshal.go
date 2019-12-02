package model

import (
	"strconv"
)

func (g GameID) String() string {
	return strconv.Itoa(int(g))
}

// var _ encoding.TextMarshaler = (*GameID)(nil)
// var _ encoding.TextUnmarshaler = (*GameID)(nil)

// func (gID *GameID) MarshalText() ([]byte, error) {
// 	s := strconv.Itoa(int(*gID))
// 	return []byte(s), nil
// }

// func (gID *GameID) UnmarshalText(input []byte) error {
// 	n, err := strconv.Atoi(string(input))
// 	if err != nil {
// 		return err
// 	}
// 	*gID = GameID(n)
// 	return nil
// }

// var _ encoding.TextMarshaler = (*PlayerID)(nil)
// var _ encoding.TextUnmarshaler = (*PlayerID)(nil)

// func (pID *PlayerID) MarshalText() ([]byte, error) {
// 	return []byte(*pID), nil
// }

// func (pID *PlayerID) UnmarshalText(input []byte) error {
// 	*pID = PlayerID(input)
// 	return nil
// }

// var _ encoding.TextMarshaler = (*PlayerColor)(nil)
// var _ encoding.TextUnmarshaler = (*PlayerColor)(nil)

// func (pc *PlayerColor) MarshalText() ([]byte, error) {
// 	return []byte(pc.String()), nil
// }

// func (pc *PlayerColor) UnmarshalText(input []byte) error {
// 	inputStr := string(input)
// 	for i := Green; i <= 5; i++ {
// 		if i.String() == inputStr {
// 			*pc = i
// 			return nil
// 		}
// 	}
// 	return nil
// }
