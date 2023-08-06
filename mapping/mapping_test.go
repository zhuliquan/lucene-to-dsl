package mapping

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestLoadMappingFile(t *testing.T) {
	type args struct {
		mappingPath string
	}
	tests := []struct {
		name        string
		args        args
		wantMapping *Mapping
		wantErr     bool
	}{
		{
			name:        "test_load_file_error",
			args:        args{mappingPath: "./test_mapping_file/not_exists_mapping.json"},
			wantMapping: nil,
			wantErr:     true,
		},
		{
			name:        "test_wrong_mapping_file",
			args:        args{mappingPath: "./test_mapping_file/wrong_mapping.json"},
			wantMapping: nil,
			wantErr:     true,
		},
		{
			name: "test_load_alias_mapping",
			args: args{mappingPath: "./test_mapping_file/alias_mapping.json"},
			wantMapping: &Mapping{
				Source:  &Source{Enabled: true},
				Dynamic: BoolDynamic(true),
				Properties: map[string]*Property{
					"host_name": {
						Type: KEYWORD_FIELD_TYPE,
					},
					"host_name_alias": {
						Type: ALIAS_FIELD_TYPE,
						Path: "host_name",
					},
				},
			},
			wantErr: false,
		},
		{
			name:        "test_load_wrong_alias_mapping",
			args:        args{mappingPath: "./test_mapping_file/wrong_alias_mapping.json"},
			wantMapping: nil,
			wantErr:     true,
		},
		{
			name: "test_load_keyword_mapping",
			args: args{mappingPath: "./test_mapping_file/keyword_mapping.json"},
			wantMapping: &Mapping{
				Source:  &Source{Enabled: true},
				Dynamic: StringDynamic("true"),
				Properties: map[string]*Property{
					"host_name": {
						Type: KEYWORD_FIELD_TYPE,
					},
					"created_at": {
						Type:   DATE_FIELD_TYPE,
						Format: "EEE MMM dd HH:mm:ss Z yyyy",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_load_object_mapping",
			args: args{mappingPath: "./test_mapping_file/object_mapping.json"},
			wantMapping: &Mapping{
				Dynamic: BoolDynamic(false),
				Properties: map[string]*Property{
					"region": {Type: KEYWORD_FIELD_TYPE},
					"manager": {
						Type: OBJECT_FIELD_TYPE,
						Mapping: Mapping{
							Dynamic: BoolDynamic(false),
							Properties: map[string]*Property{
								"age": {Type: INTEGER_FIELD_TYPE},
								"name": {
									Type: OBJECT_FIELD_TYPE,
									Mapping: Mapping{
										Dynamic: BoolDynamic(false),
										Properties: map[string]*Property{
											"first": {Type: "text"},
											"last":  {Type: "text"},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_load_flattened_mapping",
			args: args{mappingPath: "./test_mapping_file/flattened_mapping.json"},
			wantMapping: &Mapping{
				Dynamic: BoolDynamic(true),
				Properties: map[string]*Property{
					"title":  {Type: TEXT_FIELD_TYPE},
					"labels": {Type: FLATTENED_FIELD_TYPE},
				},
			},
			wantErr: false,
		},
		{
			name: "test_load_nested_mapping",
			args: args{mappingPath: "./test_mapping_file/nested_mapping.json"},
			wantMapping: &Mapping{
				Dynamic: BoolDynamic(true),
				Properties: map[string]*Property{
					"user": {Type: NESTED_FIELD_TYPE},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadMappingFile(tt.args.mappingPath, nil)
			assert.Equal(t, tt.wantErr, (err != nil))
			if got != nil {
				assert.Equal(t, tt.wantMapping, got.fieldMapping)
			}
		})
	}
}

func TestMappingString(t *testing.T) {
	m, err := LoadMappingFile("./test_mapping_file/keyword_mapping.json", nil)
	assert.Nil(t, err)
	t.Logf("got mapping file: %s", m.fieldMapping.String())
}

func TestGetProperty(t *testing.T) {
	pm, err := LoadMappingFile("./test_mapping_file/property_mapping.json", nil)
	if err != nil {
		t.Errorf("expect don't got error")
	}
	type testCase struct {
		name  string
		field string
		prop  map[string]*Property
	}

	for _, tt := range []testCase{
		{
			name:  "test_get_alias",
			field: "host_name_alias",
			prop: map[string]*Property{
				"host_name": {
					Type: "keyword",
				},
			},
		},
		{
			name:  "test_get_created_at",
			field: "created_at",
			prop: map[string]*Property{
				"created_at": {
					Type:   "date",
					Format: "EEE MMM dd HH:mm:ss Z yyyy",
				},
			},
		},
		{
			name:  "test_retry_get_created_at",
			field: "created_at",
			prop: map[string]*Property{
				"created_at": {
					Type:   "date",
					Format: "EEE MMM dd HH:mm:ss Z yyyy",
				},
			},
		},
		{
			name:  "test_not_find_error",
			field: "not_find_field",
			prop:  map[string]*Property{},
		},
		{
			name:  "test_not_support_lucene",
			field: "shape_field",
			prop:  map[string]*Property{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			prop, _ := pm.GetProperty(tt.field)
			assert.Equal(t, tt.prop, prop)
		})
	}
}

func TestGetExtFuncs(t *testing.T) {
	pm, err := LoadMappingFile("./test_mapping_file/property_mapping.json", nil)
	assert.Nil(t, err)
	f := pm.GetExtFuncs("field")
	assert.Nil(t, f)
	pm, err = LoadMappingFile("./test_mapping_file/property_mapping.json", map[string]ConvertFunc{})
	assert.Nil(t, err)
	f = pm.GetExtFuncs("field")
	assert.Nil(t, f)
}

func TestUnmarshalMapping(t *testing.T) {
	s := `{
		"dynamic": "true",
		"properties": {
			"x": {"type": "keyword", "index": true},
			"y": {"type": "text", "index": "true"},
			"o1": {"type": "object", "dynamic": "true"},
			"o2": {"type": "object", "dynamic": true}
		}
	}`
	mm := &Mapping{}
	err := jsoniter.UnmarshalFromString(s, mm)
	assert.Nil(t, err)
	assert.Equal(t, &Mapping{
		Dynamic: StringDynamic("true"),
		Properties: map[string]*Property{
			"x":  {Type: KEYWORD_FIELD_TYPE, Index: BoolValue(true)},
			"y":  {Type: TEXT_FIELD_TYPE, Index: StringValue("true")},
			"o1": {Type: OBJECT_FIELD_TYPE, Mapping: Mapping{Dynamic: StringDynamic("true")}},
			"o2": {Type: OBJECT_FIELD_TYPE, Mapping: Mapping{Dynamic: BoolDynamic(true)}},
		},
	}, mm)

	s = `{
		"dynamic": 90,
		"properties": {
			"x": {"type": "keyword", "index": true},
			"y": {"type": "text", "index": "true"},
			"o1": {"type": "object", "dynamic": "true"},
			"o2": {"type": "object", "dynamic": true}
		}
	}`
	mm = &Mapping{}
	err = jsoniter.UnmarshalFromString(s, mm)
	assert.NotNil(t, err)

	mm = &Mapping{
		Dynamic: StringDynamic(""),
		Properties: map[string]*Property{
			"x": {Type: KEYWORD_FIELD_TYPE, Index: StringValue("")},
			"y": {Type: KEYWORD_FIELD_TYPE, Index: StringValue("false")},
			"z": {Type: KEYWORD_FIELD_TYPE, Index: BoolValue(true)},
			"a": {Type: KEYWORD_FIELD_TYPE, Index: BoolValue(false)},
			"b": {Type: KEYWORD_FIELD_TYPE},
		},
	}
	_, err = jsoniter.MarshalToString(mm)
	assert.Nil(t, err)
}

func TestBoolOrString(t *testing.T) {
	b := BoolValue(true)
	assert.Equal(t, true, b.GetBool())
	assert.Equal(t, "true", b.GetString())

	s := StringValue("true")
	assert.Equal(t, true, s.GetBool())
	assert.Equal(t, "true", s.GetString())
}

func TestDynamic(t *testing.T) {
	b := BoolDynamic(true)
	assert.Equal(t, DYNAMIC_MAPPING, b.GetMappingType())

	b = BoolDynamic(false)
	assert.Equal(t, STATIC_MAPPING, b.GetMappingType())

	s := StringDynamic("true")
	assert.Equal(t, DYNAMIC_MAPPING, s.GetMappingType())

	s = StringDynamic("false")
	assert.Equal(t, STATIC_MAPPING, s.GetMappingType())

	s = StringDynamic("")
	assert.Equal(t, DYNAMIC_MAPPING, s.GetMappingType())

	s = StringDynamic("strict")
	assert.Equal(t, STRICT_MAPPING, s.GetMappingType())

	s = StringDynamic("runtime")
	assert.Equal(t, RUNTIME_MAPPING, s.GetMappingType())
}
