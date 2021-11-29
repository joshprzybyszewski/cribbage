package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCORSConfig(t *testing.T) {
	c := getCORSConfig()
	assert.Equal(t, []string{`GET`, `POST`}, c.AllowMethods)
	assert.Equal(t, []string{`https://hobbycribbage.com`, `http://localhost:3000`}, c.AllowOrigins)
}
