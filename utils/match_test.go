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

func TestWildCardPattern(t *testing.T) {
	m := NewWildCardPattern("a*")
	assert.True(t, m.Match([]byte("aa")))
	m = NewWildCardPattern("你好*")
	assert.True(t, m.Match([]byte("你好中国")))
}

func TestPrefixPattern(t *testing.T) {
	m := NewPrefixPattern("a")
	assert.True(t, m.Match([]byte("ab")))
	m = NewPrefixPattern("你好")
	assert.True(t, m.Match([]byte("你好中国")))
}
