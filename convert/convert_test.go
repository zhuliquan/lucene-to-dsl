package convert

import (
	"math"
	"reflect"
	"testing"

	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	op "github.com/zhuliquan/lucene_parser/operator"
	term "github.com/zhuliquan/lucene_parser/term"
)

func Test_convertToGroup(t *testing.T) {
	type args struct {
		field     *term.Field
		termGroup *term.TermGroup
		property  *mapping.Property
	}
	tests := []struct {
		name    string
		args    args
		want    dsl.DSLNode
		wantErr bool
	}{
		{
			name: "test_empty_01",
			args: args{
				field:     nil,
				termGroup: nil,
			},
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
			want:    nil,
			wantErr: true,
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
				property: &mapping.Property{
					Type: mapping.TEXT_FIELD_TYPE,
				},
			},
			want: &dsl.MatchPhraseNode{
				EqNode: dsl.EqNode{
					Field: "x",
					Type:  mapping.TEXT_FIELD_TYPE,
					Value: "78",
				},
			},
			wantErr: false,
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
				property: &mapping.Property{
					Type: mapping.TEXT_FIELD_TYPE,
				},
			},
			want: &dsl.QueryStringNode{
				EqNode: dsl.EqNode{
					Field: "x",
					Type:  mapping.TEXT_FIELD_TYPE,
					Value: "78",
				},
				Boost: 0.8,
			},
			wantErr: false,
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
				property: &mapping.Property{
					Type: mapping.TEXT_FIELD_TYPE,
				},
			},
			want: &dsl.MatchPhraseNode{
				EqNode: dsl.EqNode{
					Field: "x",
					Type:  mapping.TEXT_FIELD_TYPE,
					Value: "78",
				},
				Boost: 0.8,
			},
			wantErr: true,
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
				property: &mapping.Property{
					Type: mapping.INTEGER_FIELD_TYPE,
				},
			},
			want: &dsl.RangeNode{
				Field:       "x",
				ValueType:   mapping.INTEGER_FIELD_TYPE,
				LeftValue:   int32(78),
				RightValue:  math.MaxInt32,
				LeftCmpSym:  dsl.GT,
				RightCmpSym: dsl.LT,
				Boost:       0.8,
			},
			wantErr: false,
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
				property: &mapping.Property{
					Type: mapping.TEXT_FIELD_TYPE,
				},
			},
			want: &dsl.RangeNode{
				Field:       "x",
				ValueType:   mapping.TEXT_FIELD_TYPE,
				LeftValue:   "",
				RightValue:  "2006-01-01",
				LeftCmpSym:  dsl.GT,
				RightCmpSym: dsl.LTE,
				Boost:       0.8,
			},
			wantErr: false,
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
				property: &mapping.Property{
					Type: mapping.INTEGER_FIELD_TYPE,
				},
			},
			want: &dsl.RangeNode{
				Field:       "x",
				ValueType:   mapping.INTEGER_FIELD_TYPE,
				LeftValue:   int32(78),
				RightValue:  int32(100),
				LeftCmpSym:  dsl.GT,
				RightCmpSym: dsl.LT,
				Boost:       0.8,
			},
			wantErr: true,
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
				property: &mapping.Property{
					Type: mapping.INTEGER_FIELD_TYPE,
				},
			},
			// {*, 50} {78, 100} {178,*}
			want: &dsl.OrDSLNode{
				MinimumShouldMatch: 1,
				Nodes: map[string][]dsl.DSLNode{
					"x": {
						&dsl.RangeNode{
							Field:       "x",
							ValueType:   mapping.INTEGER_FIELD_TYPE,
							LeftValue:   dsl.MinInt[32],
							RightValue:  int32(50),
							LeftCmpSym:  dsl.GT,
							RightCmpSym: dsl.LT,
							Boost:       0.8,
						},
						&dsl.RangeNode{
							Field:       "x",
							ValueType:   mapping.INTEGER_FIELD_TYPE,
							LeftValue:   int32(78),
							RightValue:  int32(100),
							LeftCmpSym:  dsl.GT,
							RightCmpSym: dsl.LT,
							Boost:       0.8,
						},
						&dsl.RangeNode{
							Field:       "x",
							ValueType:   mapping.INTEGER_FIELD_TYPE,
							LeftValue:   int32(178),
							RightValue:  dsl.MaxInt[32],
							LeftCmpSym:  dsl.GT,
							RightCmpSym: dsl.LT,
							Boost:       0.8,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := convertToGroup(tt.args.field, &term.Term{TermGroup: tt.args.termGroup}, tt.args.property); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToGroup() = %v, want %v", got, tt.want)
			} else if (err != nil) != tt.wantErr {
				t.Errorf("convertToGroup(), err = %v, want err: %v", err, tt.wantErr)
			}
		})
	}
}
