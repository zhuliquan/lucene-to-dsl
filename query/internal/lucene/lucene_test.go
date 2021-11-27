package query

import (
	"testing"

	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

func TestLucene(t *testing.T) {
	var luceneParser = participle.MustBuild(
		&Lucene{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name    string
		input   string
		wantErr bool
	}

	var testCases = []testCase{
		{
			name:    "TestLucene01",
			input:   `x:1 AND NOT x:2`,
			wantErr: false,
		},
		{
			name:    "TestLucene02",
			input:   `NOT (x:1 AND y:2) OR z:9`,
			wantErr: false,
		},
		{
			name:    "TestLucne03",
			input:   `(x:1 AND NOT y:2) AND (NOT x:8 AND k:90)`,
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var lucene = &Lucene{}
			if err := luceneParser.ParseString(tt.input, lucene); (err != nil) != tt.wantErr {
				t.Errorf("parser lucene, err: %+v", err)
			}
		})
	}
}
