package dsl

import (
	"bytes"
	"math"
	"testing"
)

func TestPrefixNode(t *testing.T) {
	var s1 = "我们是"
	var s2 = "我们~"
	var s3 = "this"
	var s4 = "我们是this"
	var s5 = ""
	var s6 = string([]rune{math.MinInt32})

	t.Log([]byte(s1))
	t.Log([]rune(s1))

	t.Log([]byte(s2))
	t.Log([]rune(s2))

	t.Log([]byte(s3))
	t.Log([]rune(s3))

	t.Log([]byte(s4))
	t.Log([]rune(s4))
	t.Log(bytes.Compare([]byte(s1), []byte(s2)))
	t.Log(s5 < s6)
}
