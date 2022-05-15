package mapping

import (
	"reflect"
	"testing"
)

func Test_checkTypeSupportLucene(t *testing.T) {
	type args struct {
		typ FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_true",
			args: args{typ: KEYWORD_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_false",
			args: args{typ: SHAPE_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkTypeSupportLucene(tt.args.typ); got != tt.want {
				t.Errorf("checkTypeSupportLucene() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAliasMap(t *testing.T) {
	type args struct {
		m *Mapping
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test_error_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: ALIAS_FIELD_TYPE,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: ALIAS_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: ALIAS_FIELD_TYPE,
										Path: "name.first",
									},
								},
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_04",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"alias_type": {
							Type: ALIAS_FIELD_TYPE,
							Path: "other_field",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_ok_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name_2": {
							Type: ALIAS_FIELD_TYPE,
							Path: "name_second",
						},
						"name_second": {
							Type: TEXT_FIELD_TYPE,
						},
					},
				},
			},
			want:    map[string]string{"name_2": "name_second"},
			wantErr: false,
		},
		{
			name: "test_ok_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: ALIAS_FIELD_TYPE,
										Path: "name_first",
									},
								},
							},
						},
						"name_first": {
							Type: TEXT_FIELD_TYPE,
						},
					},
				},
			},
			want:    map[string]string{"name.first": "name_first"},
			wantErr: false,
		},
		{
			name: "test_ok_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: TEXT_FIELD_TYPE,
									},
								},
							},
						},
						"name_first": {
							Type: ALIAS_FIELD_TYPE,
							Path: "name.first",
						},
					},
				},
			},
			want:    map[string]string{"name_first": "name.first"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pm = &PropertyMapping{
				_mapping:  tt.args.m,
				_cacheMap: map[string]*Property{},
				_aliasMap: map[string]string{},
			}
			got, err := getAliasMap(pm)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAliasMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAliasMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchProperty(t *testing.T) {
	type args struct {
		m      *Mapping
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    *Property
		wantErr bool
	}{
		{
			name: "test_error_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"y": {Type: TEXT_FIELD_TYPE},
					},
				},
				target: "x",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: TEXT_FIELD_TYPE,
						},
					},
				},
				target: "x.y",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: TEXT_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
				target: "x",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_04",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: OBJECT_FIELD_TYPE,
										Mapping: Mapping{
											Properties: map[string]*Property{
												"z": {
													Type: TEXT_FIELD_TYPE,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				target: "x.y",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_ok_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: FLATTENED_FIELD_TYPE,
						},
					},
				},
				target: "x",
			},
			want: &Property{
				Type: FLATTENED_FIELD_TYPE,
			},
			wantErr: false,
		},
		{
			name: "test_ok_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: TEXT_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
				target: "x.y",
			},
			want: &Property{
				Type: TEXT_FIELD_TYPE,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m = &PropertyMapping{
				_mapping: tt.args.m,
			}
			got, err := getProperty(m, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchProperty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
