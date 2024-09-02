package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDsl(t *testing.T) {
	var d = DSL{}
	addValueForDSL(d, "foo", "bar")
	assert.Equal(t, DSL{"foo": "bar"}, d)
	addValueForDSL(d, "x", 0)
	assert.Equal(t, DSL{"foo": "bar"}, d)
	addValueForDSL(d, "x", 1)
	assert.Equal(t, DSL{"foo": "bar", "x": 1}, d)
	assert.Equal(t, "{\"foo\":\"bar\"}", DSL{"foo": "bar"}.String())
}

func TestLeafValue(t *testing.T) {
	assert.Equal(t, nil, EmptyValue)
}

func TestEmptyDsl(t *testing.T) {
	assert.Equal(t, DSL{}, EmptyDSL)
}
