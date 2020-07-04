package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvironment(t *testing.T) {
	defer os.Setenv(`deploy`, os.Getenv(`deploy`))

	os.Setenv(`deploy`, ``)
	assert.Equal(t, `default`, getEnvironment())

	os.Setenv(`deploy`, `prod`)
	assert.Equal(t, `prod`, getEnvironment())
}
