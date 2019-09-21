package round

type Stage int

const (
	Deal Stage = 1 << iota
	BuildCrib
	Cut
	Pegging
	Counting
)
