package term

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

func TestTermGroup(t *testing.T) {
	var termParser = participle.MustBuild(
		&LogicTermGroup{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *LogicTermGroup
	}
	var testCases = []testCase{
		{
			name:  "TestTermGroup07",
			input: `((quick AND fox) OR (brown AND fox) OR fox) AND NOT news`,
			want: &LogicTermGroup{
				OrTermGroup: &OrTermGroup{
					AndTermGroup: &AndTermGroup{
						ParenTermGroup: &ParenTermGroup{
							SubTermGroup: &LogicTermGroup{
								OrTermGroup: &OrTermGroup{
									AndTermGroup: &AndTermGroup{
										ParenTermGroup: &ParenTermGroup{
											SubTermGroup: &LogicTermGroup{
												OrTermGroup: &OrTermGroup{
													AndTermGroup: &AndTermGroup{
														TermGroupElem: &TermGroupElem{
															SingleTerm: &SingleTerm{Value: []string{"quick"}},
														},
													},
													AnSTermGroup: []*AnSTermGroup{
														{
															AndSymbol: &op.AndSymbol{Symbol: "AND"},
															AndTermGroup: &AndTermGroup{
																TermGroupElem: &TermGroupElem{
																	SingleTerm: &SingleTerm{Value: []string{"fox"}},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								OSTermGroup: []*OSTermGroup{
									{
										OrSymbol: &op.OrSymbol{Symbol: "OR"},
										OrTermGroup: &OrTermGroup{
											AndTermGroup: &AndTermGroup{
												ParenTermGroup: &ParenTermGroup{
													SubTermGroup: &LogicTermGroup{
														OrTermGroup: &OrTermGroup{
															AndTermGroup: &AndTermGroup{
																TermGroupElem: &TermGroupElem{
																	SingleTerm: &SingleTerm{Value: []string{"brown"}},
																},
															},
															AnSTermGroup: []*AnSTermGroup{
																{
																	AndSymbol: &op.AndSymbol{Symbol: "AND"},
																	AndTermGroup: &AndTermGroup{
																		TermGroupElem: &TermGroupElem{
																			SingleTerm: &SingleTerm{Value: []string{"fox"}},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									{
										OrSymbol: &op.OrSymbol{Symbol: "OR"},
										OrTermGroup: &OrTermGroup{
											AndTermGroup: &AndTermGroup{
												ParenTermGroup: &ParenTermGroup{
													SubTermGroup: &LogicTermGroup{
														OrTermGroup: &OrTermGroup{
															AndTermGroup: &AndTermGroup{
																TermGroupElem: &TermGroupElem{
																	SingleTerm: &SingleTerm{Value: []string{"fox"}},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					AnSTermGroup: []*AnSTermGroup{
						{
							AndSymbol: &op.AndSymbol{Symbol: "AND"},
							NotSymbol: &op.NotSymbol{Symbol: "NOT"},
							AndTermGroup: &AndTermGroup{
								TermGroupElem: &TermGroupElem{
									SingleTerm: &SingleTerm{Value: []string{"news"}},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &LogicTermGroup{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Log(tt.want.String())
				t.Log(out.String())
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}
