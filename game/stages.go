package game

type RoundStage int

const (
	Deal RoundStage = 1 << iota
	BuildCrib
	Cut
	Pegging
	Counting
	CribCounting
	Done
)
