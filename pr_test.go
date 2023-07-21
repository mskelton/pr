package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTicket(t *testing.T) {
	prefix, title := getDefaultTitle("fcs-1234-branch-name")

	assert.Equal(t, "FCS-1234", prefix)
	assert.Equal(t, "Branch name", title)
}

func TestWithInvalidTicket(t *testing.T) {
	prefix, title := getDefaultTitle("hi-1234-branch-name")

	assert.Equal(t, "", prefix)
	assert.Equal(t, "Hi 1234 branch name", title)
}

func TestWithoutTicket(t *testing.T) {
	prefix, title := getDefaultTitle("branch-name")

	assert.Equal(t, "", prefix)
	assert.Equal(t, "Branch name", title)
}
