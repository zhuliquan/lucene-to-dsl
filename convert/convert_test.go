package convert

import (
	"reflect"
	"testing"

	lucene "github.com/zhuliquan/lucene_parser"
	op "github.com/zhuliquan/lucene_parser/operator"
	term "github.com/zhuliquan/lucene_parser/term"
)

func Test_convertTermGroupToLucene(t *testing.T) {
	type args struct {
		field     *term.Field
		termGroup *term.TermGroup
	}
	tests := []struct {
		name string
		args args
		want *lucene.Lucene
	}{
		{
			name: "test_empty_01",
			args: args{
				field:     nil,
				termGroup: nil,
			},
			want: nil,
		},
		{
			name: "test_empty_02",
			args: args{
				field: &term.Field{
					Value: []string{"x"},
				},
				termGroup: &term.TermGroup{
					LogicTermGroup: nil,
					BoostSymbol:    "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_03",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{},
					BoostSymbol:    "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_04",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_05",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_06",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_07",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								ParenTermGroup: &term.ParenTermGroup{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_08",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								ParenTermGroup: &term.ParenTermGroup{
									SubTermGroup: &term.LogicTermGroup{},
								},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_empty_09",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								ParenTermGroup: &term.ParenTermGroup{
									SubTermGroup: &term.LogicTermGroup{},
								},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: nil,
		},
		{
			name: "test_half_empty_00",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									PhraseTerm: &term.PhraseTerm{Chars: []string{"78"}},
								},
							},
							AnSTermGroup: []*term.AnSTermGroup{
								{
									AndSymbol: &op.AndSymbol{Symbol: "AND"},
									AndTermGroup: &term.AndTermGroup{
										TermGroupElem: nil,
									},
								},
								nil,
							},
						},
						OSTermGroup: []*term.OSTermGroup{nil},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								FuzzyTerm: &term.FuzzyTerm{
									PhraseTerm:  &term.PhraseTerm{Chars: []string{"78"}},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "test_half_empty_01",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									SingleTerm: &term.SingleTerm{Begin: "78"},
								},
							},
						},
						OSTermGroup: []*term.OSTermGroup{
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								FuzzyTerm: &term.FuzzyTerm{
									SingleTerm:  &term.SingleTerm{Begin: "78"},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "test_half_empty_02",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									PhraseTerm: &term.PhraseTerm{Chars: []string{"78"}},
								},
							},
						},
						OSTermGroup: []*term.OSTermGroup{
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								FuzzyTerm: &term.FuzzyTerm{
									PhraseTerm:  &term.PhraseTerm{Chars: []string{"78"}},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "test_half_empty_03",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									SRangeTerm: &term.SRangeTerm{
										Symbol: ">",
										Value: &term.RangeValue{
											SingleValue: []string{"78"},
										},
									},
								},
							},
						},
						OSTermGroup: []*term.OSTermGroup{
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								RangeTerm: &term.RangeTerm{
									SRangeTerm: &term.SRangeTerm{
										Symbol: ">",
										Value: &term.RangeValue{
											SingleValue: []string{"78"},
										},
									},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "test_half_empty_04",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									DRangeTerm: &term.DRangeTerm{
										LBRACKET: "{",
										LValue:   &term.RangeValue{InfinityVal: "*"},
										RValue:   &term.RangeValue{PhraseValue: []string{"2006", "-", "01", "-", "01"}},
										RBRACKET: "]",
									},
								},
							},
						},
						OSTermGroup: []*term.OSTermGroup{
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								RangeTerm: &term.RangeTerm{
									DRangeTerm: &term.DRangeTerm{
										LBRACKET: "{",
										LValue:   &term.RangeValue{InfinityVal: "*"},
										RValue:   &term.RangeValue{PhraseValue: []string{"2006", "-", "01", "-", "01"}},
										RBRACKET: "]",
									},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "test_half_empty_05",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									SRangeTerm: &term.SRangeTerm{
										Symbol: ">",
										Value: &term.RangeValue{
											SingleValue: []string{"78"},
										},
									},
								},
							},
							AnSTermGroup: []*term.AnSTermGroup{
								{
									AndSymbol: &op.AndSymbol{Symbol: "&&"},
									AndTermGroup: &term.AndTermGroup{
										TermGroupElem: &term.TermGroupElem{
											SRangeTerm: &term.SRangeTerm{
												Symbol: "<",
												Value: &term.RangeValue{
													SingleValue: []string{"100"},
												},
											},
										},
									},
								},
							},
						},
						OSTermGroup: []*term.OSTermGroup{
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								RangeTerm: &term.RangeTerm{
									SRangeTerm: &term.SRangeTerm{
										Symbol: ">",
										Value: &term.RangeValue{
											SingleValue: []string{"78"},
										},
									},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
					AnSQuery: []*lucene.AnSQuery{
						{
							AndSymbol: &op.AndSymbol{Symbol: "&&"},
							AndQuery: &lucene.AndQuery{
								FieldQuery: &lucene.FieldQuery{
									Field: &term.Field{Value: []string{"x"}},
									Term: &term.Term{
										RangeTerm: &term.RangeTerm{
											SRangeTerm: &term.SRangeTerm{
												Symbol: "<",
												Value: &term.RangeValue{
													SingleValue: []string{"100"},
												},
											},
											BoostSymbol: "^0.8",
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
			name: "test_ok",
			args: args{
				field: &term.Field{Value: []string{"x"}},
				termGroup: &term.TermGroup{
					LogicTermGroup: &term.LogicTermGroup{
						OrTermGroup: &term.OrTermGroup{
							AndTermGroup: &term.AndTermGroup{
								TermGroupElem: &term.TermGroupElem{
									SRangeTerm: &term.SRangeTerm{
										Symbol: ">",
										Value: &term.RangeValue{
											SingleValue: []string{"78"},
										},
									},
								},
							},
							AnSTermGroup: []*term.AnSTermGroup{
								{
									AndSymbol: &op.AndSymbol{Symbol: "&&"},
									AndTermGroup: &term.AndTermGroup{
										TermGroupElem: &term.TermGroupElem{
											SRangeTerm: &term.SRangeTerm{
												Symbol: "<",
												Value: &term.RangeValue{
													SingleValue: []string{"100"},
												},
											},
										},
									},
								},
							},
						},
						OSTermGroup: []*term.OSTermGroup{
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{
									AndTermGroup: &term.AndTermGroup{
										TermGroupElem: &term.TermGroupElem{
											SRangeTerm: &term.SRangeTerm{
												Symbol: ">",
												Value: &term.RangeValue{
													SingleValue: []string{"178"},
												},
											},
										},
									},
								},
							},
							{
								OrSymbol: &op.OrSymbol{
									Symbol: "||",
								},
								OrTermGroup: &term.OrTermGroup{
									AndTermGroup: &term.AndTermGroup{
										ParenTermGroup: &term.ParenTermGroup{
											SubTermGroup: &term.LogicTermGroup{
												OrTermGroup: &term.OrTermGroup{
													AndTermGroup: &term.AndTermGroup{
														TermGroupElem: &term.TermGroupElem{
															SRangeTerm: &term.SRangeTerm{
																Symbol: ">",
																Value: &term.RangeValue{
																	SingleValue: []string{"10"},
																},
															},
														},
													},
													AnSTermGroup: []*term.AnSTermGroup{
														{
															AndSymbol: &op.AndSymbol{Symbol: " AND "},
															AndTermGroup: &term.AndTermGroup{
																TermGroupElem: &term.TermGroupElem{
																	SRangeTerm: &term.SRangeTerm{
																		Symbol: "<",
																		Value: &term.RangeValue{
																			SingleValue: []string{"50"},
																		},
																	},
																},
															},
														},
													},
												},
												OSTermGroup: []*term.OSTermGroup{},
											},
										},
									},
								},
							},
						},
					},
					BoostSymbol: "^0.8",
				},
			},
			want: &lucene.Lucene{
				OrQuery: &lucene.OrQuery{
					AndQuery: &lucene.AndQuery{
						FieldQuery: &lucene.FieldQuery{
							Field: &term.Field{Value: []string{"x"}},
							Term: &term.Term{
								RangeTerm: &term.RangeTerm{
									SRangeTerm: &term.SRangeTerm{
										Symbol: ">",
										Value: &term.RangeValue{
											SingleValue: []string{"78"},
										},
									},
									BoostSymbol: "^0.8",
								},
							},
						},
					},
					AnSQuery: []*lucene.AnSQuery{
						{
							AndSymbol: &op.AndSymbol{Symbol: "&&"},
							AndQuery: &lucene.AndQuery{
								FieldQuery: &lucene.FieldQuery{
									Field: &term.Field{Value: []string{"x"}},
									Term: &term.Term{
										RangeTerm: &term.RangeTerm{
											SRangeTerm: &term.SRangeTerm{
												Symbol: "<",
												Value: &term.RangeValue{
													SingleValue: []string{"100"},
												},
											},
											BoostSymbol: "^0.8",
										},
									},
								},
							},
						},
					},
				},
				OSQuery: []*lucene.OSQuery{
					{
						OrSymbol: &op.OrSymbol{
							Symbol: "||",
						},
						OrQuery: &lucene.OrQuery{
							AndQuery: &lucene.AndQuery{
								FieldQuery: &lucene.FieldQuery{
									Field: &term.Field{Value: []string{"x"}},
									Term: &term.Term{
										RangeTerm: &term.RangeTerm{
											SRangeTerm: &term.SRangeTerm{
												Symbol: ">",
												Value: &term.RangeValue{
													SingleValue: []string{"178"},
												},
											},
											BoostSymbol: "^0.8",
										},
									},
								},
							},
						},
					},
					{
						OrSymbol: &op.OrSymbol{
							Symbol: "||",
						},
						OrQuery: &lucene.OrQuery{
							AndQuery: &lucene.AndQuery{
								ParenQuery: &lucene.ParenQuery{
									SubQuery: &lucene.Lucene{
										OrQuery: &lucene.OrQuery{
											AndQuery: &lucene.AndQuery{
												FieldQuery: &lucene.FieldQuery{
													Field: &term.Field{Value: []string{"x"}},
													Term: &term.Term{
														RangeTerm: &term.RangeTerm{
															SRangeTerm: &term.SRangeTerm{
																Symbol: ">",
																Value: &term.RangeValue{
																	SingleValue: []string{"10"},
																},
															},
															BoostSymbol: "^0.8",
														},
													},
												},
											},
											AnSQuery: []*lucene.AnSQuery{
												{
													AndSymbol: &op.AndSymbol{Symbol: " AND "},
													AndQuery: &lucene.AndQuery{
														FieldQuery: &lucene.FieldQuery{
															Field: &term.Field{Value: []string{"x"}},
															Term: &term.Term{
																RangeTerm: &term.RangeTerm{
																	SRangeTerm: &term.SRangeTerm{
																		Symbol: "<",
																		Value: &term.RangeValue{
																			SingleValue: []string{"50"},
																		},
																	},
																	BoostSymbol: "^0.8",
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
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTermGroupToLucene(tt.args.field, tt.args.termGroup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTermGroupToLucene() = %v, want %v", got, tt.want)
			}
		})
	}
}
