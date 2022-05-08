package mapping

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
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
			name: "test_load_keyword_mapping",
			args: args{mappingPath: "./test_mapping_file/keyword_mapping.json"},
			wantMapping: &Mapping{
				Source: &Source{Enabled: true},
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
				Properties: map[string]*Property{
					"region": {Type: KEYWORD_FIELD_TYPE},
					"manager": {
						Mapping: Mapping{
							Properties: map[string]*Property{
								"age": {Type: INTEGER_FIELD_TYPE},
								"name": {
									Mapping: Mapping{
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
				Properties: map[string]*Property{
					"user": {Type: NESTED_FIELD_TYPE},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.mappingPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got._mapping, tt.wantMapping) {
				t.Errorf("LoadMapping() = %v, want %v", got._mapping, tt.wantMapping)
			}
		})
	}
}
