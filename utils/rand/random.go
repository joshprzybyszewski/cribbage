package rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var (
	randChars []rune = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_`)
)

func Intn(max int) int {
	return int(Int64n(int64(max)))
}

func Int64n(max int64) int64 {
	randBigInt, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		// rand.Int should never fail
		fmt.Printf("randIntn got error: %+v\n", err)
		// Thirteen isn't random, but it's prime. So at least it's got that going for it
		return 13
	}
	return randBigInt.Int64()
}

// Int returns a random integer between 0 and 10000
// Since this is a game of cribbage, we likely won't
// need any larger number
func Int() int {
	return int(Int64n(int64(1000)))
}

// Float64 returns a random number between 0 and 1
// It will be granular to the ten-thousandths since
// this is a cribbage game and does not need finer
// floats than that
func Float64() float64 {
	return float64(Int64n(10000)) / 10000.0
}

func String(n int) string {
	s := ``
	for len(s) < n {
		s += string(randChars[Intn(len(randChars))])
	}
	return s
}
