package main

import (
	"testing"

	"github.com/gotify/plugin-api"
	"github.com/stretchr/testify/assert"
)

func TestAPICompatibility(t *testing.T) {
	assert.Implements(t, (*plugin.Plugin)(nil), new(MyPlugin))
	// Add other interfaces you intend to implement here
}
