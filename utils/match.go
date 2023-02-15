package utils

import (
	"strings"
)

type PatternMatcher interface {
	Match([]byte) bool
}

// prefix matcher
type prefixPattern struct {
	pattern string
}

func NewPrefixPattern(pattern string) PatternMatcher {
	return &prefixPattern{pattern: pattern}
}

func (p *prefixPattern) Match(text []byte) bool {
	return strings.HasPrefix(string(text), p.pattern)
}

// wildcard matcher
type wildcardPattern struct {
	pattern []rune
}

func NewWildCardPattern(pattern string) PatternMatcher {
	return &wildcardPattern{pattern: []rune(pattern)}
}

func (w *wildcardPattern) Match(text []byte) bool {
	return WildcardMatch([]rune(string(text)), w.pattern)
}

// wildcard match text and pattern
func WildcardMatch(text []rune, pattern []rune) bool {
	var n, m = len(text), len(pattern)
	var dp = make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]bool, m+1)
	}

	dp[0][0] = true
	for i := 1; i <= n; i++ {
		dp[i][0] = false
	}
	for j := 1; j <= m; j++ {
		if pattern[j-1] == '*' {
			dp[0][j] = dp[0][j-1]
		}
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if pattern[j-1] == '?' || pattern[j-1] == text[i-1] {
				dp[i][j] = dp[i-1][j-1]
			} else if pattern[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else {
				dp[i][j] = false
			}
		}
	}
	return dp[n][m]
}
