package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTicket(t *testing.T) {
	prefix, title := getDefaultTitle([]string{"ABC"}, "abc-1234-branch-name")
	assert.Equal(t, "ABC-1234", prefix)
    assert.Equal(t, "Branchname", title)
}

func TestWithInvalidTicket(t *testing.T) {
	prefix, title := getDefaultTitle([]string{"ABC"}, "hi-1234-branch-name")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "Hi 1234 branch name", title)
}

func TestWithoutTicket(t *testing.T) {
	prefix, title := getDefaultTitle([]string{"ABC"}, "branch-name")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "Branch name", title)
}

func TestWithMixedCaseTicket(t *testing.T) {
	prefix, title := getDefaultTitle([]string{"aBc"}, "Abc-1234-branch-name")
	assert.Equal(t, "ABC-1234", prefix)
	assert.Equal(t, "Branch name", title)
}

func TestWithMultiplePrefixes(t *testing.T) {
	prefixes := []string{"abc", "def"}
	prefix, title := getDefaultTitle(prefixes, "def-1234-branch-name")
	assert.Equal(t, "DEF-1234", prefix)
	assert.Equal(t, "Branch name", title)

	prefix, title = getDefaultTitle(prefixes, "def-1234-branch-name")
	assert.Equal(t, "DEF-1234", prefix)
	assert.Equal(t, "Branch name", title)
}
