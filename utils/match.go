package utils

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