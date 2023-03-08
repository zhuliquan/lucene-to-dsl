package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	op "github.com/zhuliquan/lucene_parser/operator"
	term "github.com/zhuliquan/lucene_parser/term"
)

func TestConvertToGroup(t *testing.T) {
	okNode := dsl.NewBoolNode(dsl.NewRangeNode(
		dsl.NewRgNode(
			dsl.NewFieldNode(dsl.NewLfNode(), "x"),
			dsl.NewValueType(mapping.INTEGER_FIELD_TYPE, true),
			int64(78), dsl.MaxInt[32], dsl.GT, dsl.LT,
		),
		dsl.WithBoost(0.8),
	), dsl.OR)
	okNode, _ = okNode.UnionJoin(dsl.NewRangeNode(
		dsl.NewRgNode(
			dsl.NewFieldNode(dsl.NewLfNode(), "x"),
			dsl.NewValueType(mapping.INTEGER_FIELD_TYPE, true),
			int64(10), int64(50), dsl.GT, dsl.LT,
		),
		dsl.WithBoost(0.8),
	))
	type args struct {
		field     *term.Field
		termGroup *term.TermGroup
		property  *mapping.Property
	}
	tests := []struct {
		name    string
		args    args
		want    dsl.AstNode
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
			want: dsl.NewMatchPhraseNode(
				dsl.NewKVNode(
					dsl.NewFieldNode(dsl.NewLfNode(), "x"),
					dsl.NewValueNode("78", dsl.NewValueType(mapping.TEXT_FIELD_TYPE, true)),
				),
				dsl.WithBoost(0.8),
			),
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
			want: dsl.NewQueryStringNode(
				dsl.NewKVNode(
					dsl.NewFieldNode(dsl.NewLfNode(), "x"),
					dsl.NewValueNode("78", dsl.NewValueType(mapping.TEXT_FIELD_TYPE, true)),
				),
				dsl.WithBoost(0.8),
			),
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
			want: dsl.NewMatchPhraseNode(
				dsl.NewKVNode(
					dsl.NewFieldNode(dsl.NewLfNode(), "x"),
					dsl.NewValueNode("78", dsl.NewValueType(mapping.TEXT_FIELD_TYPE, true)),
				),
				dsl.WithBoost(0.8),
			),
			wantErr: false,
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
			want: dsl.NewRangeNode(
				dsl.NewRgNode(
					dsl.NewFieldNode(dsl.NewLfNode(), "x"),
					dsl.NewValueType(mapping.INTEGER_FIELD_TYPE, true),
					int64(78), dsl.MaxInt[32], dsl.GT, dsl.LT,
				),
				dsl.WithBoost(0.8),
			),
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
			want: dsl.NewRangeNode(
				dsl.NewRgNode(
					dsl.NewFieldNode(dsl.NewLfNode(), "x"),
					dsl.NewValueType(mapping.TEXT_FIELD_TYPE, true),
					"", "2006-01-01", dsl.GT, dsl.LTE,
				),
				dsl.WithBoost(0.8),
			),
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
			want: dsl.NewRangeNode(
				dsl.NewRgNode(
					dsl.NewFieldNode(dsl.NewLfNode(), "x"),
					dsl.NewValueType(mapping.INTEGER_FIELD_TYPE, true),
					int64(78), int64(100), dsl.GT, dsl.LT,
				),
				dsl.WithBoost(0.8),
			),
			wantErr: false,
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
																	SingleValue: []string{"90"},
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
																			SingleValue: []string{"180"},
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
			want:    okNode,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &converter{}
			got, err := cc.convertToGroup(tt.args.field, &term.Term{TermGroup: tt.args.termGroup}, tt.args.property)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
