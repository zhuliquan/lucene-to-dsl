package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWildcardMatch(t *testing.T) {
	assert.True(t, WildcardMatch([]rune(""), []rune("")))
	assert.False(t, WildcardMatch([]rune("a"), []rune("")))
	assert.True(t, WildcardMatch([]rune(""), []rune("*")))
	assert.True(t, WildcardMatch([]rune("a"), []rune("a*")))
	assert.True(t, WildcardMatch([]rune("ab"), []rune("a?")))
	assert.True(t, WildcardMatch([]rune("abb"), []rune("a*b")))
	assert.False(t, WildcardMatch([]rune("a"), []rune("a?")))
	assert.True(t, WildcardMatch([]rune("我们"), []rune("我?")))
}
