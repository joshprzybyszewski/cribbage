package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayerColorStringConversions(t *testing.T) {
	for _, pc := range []PlayerColor{
		UnsetColor,
		Green,
		Blue,
		Red,
	} {
		assert.Equal(t, pc, NewPlayerColorFromString(pc.String()))
	}

	assert.Equal(t, `notacolor`, unknownPlayerColor.String())
	assert.Equal(t, `notacolor`, (PlayerColor)(4).String())
	assert.Equal(t, unknownPlayerColor, NewPlayerColorFromString(`other`))
}

func TestBlockerStringConversions(t *testing.T) {
	for _, b := range []Blocker{
		DealCards,
		CribCard,
		CutCard,
		PegCard,
		CountHand,
		CountCrib,
	} {
		assert.Equal(t, b, NewBlockerFromString(b.String()))
	}

	assert.Equal(t, `InvalidBlocker`, unknownBlocker.String())
	assert.Equal(t, `InvalidBlocker`, (Blocker)(6).String())
	assert.Equal(t, unknownBlocker, NewBlockerFromString(`other`))
}

func TestPhaseStringConversions(t *testing.T) {
	for _, p := range []Phase{
		Deal,
		BuildCribReady,
		BuildCrib,
		CutReady,
		Cut,
		PeggingReady,
		Pegging,
		CountingReady,
		Counting,
		CribCountingReady,
		CribCounting,
	} {
		assert.Equal(t, p, NewPhaseFromString(p.String()))
	}

	assert.Equal(t, `unknown`, unknownPhase.String())
	assert.Equal(t, `unknown`, (Phase)(12).String())
	assert.Equal(t, unknownPhase, NewPhaseFromString(`other`))
}

func TestPlayerActionTimeStamp(t *testing.T) {
	pa := PlayerAction{}

	t0 := time.Now()
	pa.SetTimeStamp(t0)

	assert.Equal(t, t0.Format(time.RFC3339), pa.TimestampStr)
}
