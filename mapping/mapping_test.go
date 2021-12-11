package mapping

import (
	"reflect"
	"testing"
)

func TestLoadMapping(t *testing.T) {
	type args struct {
		mappingPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Mapping
		wantErr bool
	}{
		{
			name: "test_load_keyword_mapping",
			args: args{mappingPath: "./test_mapping_file/keyword_mapping.json"},
			want: &Mapping{
				Source: &Source{Enabled: true},
				Properties: map[string]*FieldMapping{
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
			name: "test_load_alias_mapping",
			args: args{mappingPath: "./test_mapping_file/alias_mapping.json"},
			want: &Mapping{
				Properties: map[string]*FieldMapping{
					"distance": {Type: LONG_FIELD_TYPE},
					"route_length_miles": {
						Type: ALIAS_FIELD_TYPE,
						Path: "distance",
					},
					"transit_mode": {
						Type: KEYWORD_FIELD_TYPE,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_load_object_mapping",
			args: args{mappingPath: "./test_mapping_file/object_mapping.json"},
			want: &Mapping{
				Properties: map[string]*FieldMapping{
					"region": {Type: KEYWORD_FIELD_TYPE},
					"manager": {
						Mapping: Mapping{
							Properties: map[string]*FieldMapping{
								"age": {Type: INTEGER_FIELD_TYPE},
								"name": {
									Mapping: Mapping{
										Properties: map[string]*FieldMapping{
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
			want: &Mapping{
				Properties: map[string]*FieldMapping{
					"title":  {Type: TEXT_FIELD_TYPE},
					"labels": {Type: FLATTENED_FIELD_TYPE},
				},
			},
			wantErr: false,
		},
		{
			name: "test_load_nested_mapping",
			args: args{mappingPath: "./test_mapping_file/nested_mapping.json"},
			want: &Mapping{
				Properties: map[string]*FieldMapping{
					"user": {Type: NESTED_FIELD_TYPE},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadMapping(tt.args.mappingPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}
