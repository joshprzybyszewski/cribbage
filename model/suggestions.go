package model

type TossSummary struct {
	Kept   []Card
	Tossed []Card

	HandStats TossStats
	CribStats TossStats
}

type TossStats interface {
	Min() int
	Median() float64
	Avg() float64
	Max() int
}
