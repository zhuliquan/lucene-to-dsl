package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckTypeSupportLucene(t *testing.T) {
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
			got := checkTypeSupportLucene(tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractFieldAliasMap(t *testing.T) {
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
				fieldMapping:  tt.args.m,
				propertyCache: map[string]*Property{},
				fieldAliasMap: map[string]string{},
			}
			got, err := extractFieldAliasMap(pm)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFetchProperty(t *testing.T) {
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
			name: "test_err_05",
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
												"z.a": {
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
				target: "x.y.z",
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
			name: "test_ok_02",
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
				fieldMapping: tt.args.m,
			}
			got, err := getProperty(m, tt.args.target)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckIntType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_check_int_01",
			args: args{t: INTEGER_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_int_02",
			args: args{t: INTEGER_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_byte",
			args: args{t: BYTE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_short",
			args: args{t: SHORT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_long_01",
			args: args{t: LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_long_02",
			args: args{t: LONG_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_other",
			args: args{t: DOUBLE_RANGE_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIntType(tt.args.t); got != tt.want {
				t.Errorf("CheckIntType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUIntType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_uint64",
			args: args{t: UNSIGNED_LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUIntType(tt.args.t); got != tt.want {
				t.Errorf("CheckUIntType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckFloatType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_float16",
			args: args{t: HALF_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float32_01",
			args: args{t: FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float32_02",
			args: args{t: FLOAT_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float64_01",
			args: args{t: DOUBLE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float64_02",
			args: args{t: DOUBLE_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float128",
			args: args{t: SCALED_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckFloatType(tt.args.t); got != tt.want {
				t.Errorf("CheckFloatType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDateType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_check_date_01",
			args: args{t: DATE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_date_02",
			args: args{t: DATE_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_date_03",
			args: args{t: DATE_NANOS_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDateType(tt.args.t); got != tt.want {
				t.Errorf("CheckDateType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckIPType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_ip_01",
			args: args{t: IP_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_ip_02",
			args: args{t: IP_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIPType(tt.args.t); got != tt.want {
				t.Errorf("CheckIPType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckVersionType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_version",
			args: args{t: VERSION_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other_version",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckVersionType(tt.args.t); got != tt.want {
				t.Errorf("CheckVersionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
