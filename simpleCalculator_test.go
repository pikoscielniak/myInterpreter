package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var token Token

func setup(t *testing.T) {
	token = Token{
		Type: INTEGER,
		Value:"3",
	}
}

func Test(t *testing.T) {
	setup(t)

	assert.Equal(t, "Token({INTEGER}, {3})", token.String())
}
