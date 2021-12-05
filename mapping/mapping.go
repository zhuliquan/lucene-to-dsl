package mapping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Meta struct {
	Unit        MetaUnitType    `json:"unit,omitempty"`
	MetricsType MetaMetricsType `json:"metrics_type,omitempty"`
}

type Source struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Includes []string `json:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
}

type All struct {
	Enable bool `json:"anabled,omitempty"`
	Store  bool `json:"store,omitempty"`
}

type FieldMapping struct {
	Mapping
	Type      FieldType `json:"type,omitempty"`
	Boost     float64   `json:"boost,omitempty"`
	Index     bool      `json:"index,omitempty"`      // Should the field be searchable? Accepts true (default) and false.
	Store     bool      `json:"store,omitempty"`      // whether save field value, so it can return this field from original document.
	NullValue string    `json:"null_value,omitempty"` // using for check whether field exist.
	Path      string    `json:"path,omitempty"`       // using be alias type

	Meta                *Meta        `json:"meta,omitempty"`          // meta data for numeric type, such as long, float
	IndexOptions        IndexOptions `json:"index_options,omitempty"` // index value will be saved
	DocValues           bool         `json:"doc_values,omitempty"`    // take advantage for field aggregating and field sorting
	IncludeInAll        bool         `json:"include_in_all,omitempty"`
	DepthLimit          int          `json:"depth_limit,omitempty"`
	EagerGlobalOrdinals bool         `json:"eager_global_ordinals,omitempty,"`
	IgnoreAbove         string       `json:"ignore_above,omitempty"`
	IgnoreMalformed     bool         `json:"ignore_malformed,omitempty"`
	Format              string       `json:"format,omitempty"`

	// 自定义
	UpperCase bool `json:"upper_case,omitempty"`
	LowerCase bool `json:"lower_case,omitempty"`
}

type Mapping struct {
	MappingType  MappingType              `json:"dynamic,omitempty"`
	IncludeInAll bool                     `json:"include_in_all,omitempty"`
	Source       *Source                  `json:"_source,omitempty"`
	All          *All                     `json:"_all,omitempty"`
	Properties   map[string]*FieldMapping `json:"properties,omitempty"`
}

func (m *Mapping) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func LoadMapping(mappingPath string) (*Mapping, error) {
	if data, err := ioutil.ReadFile(mappingPath); err != nil {
		return nil, fmt.Errorf("failed to read mapping file: %s, err: %+v", mappingPath, err)
	} else {
		var mappingData = &Mapping{}
		if err := json.Unmarshal(data, mappingData); err != nil {
			return nil, fmt.Errorf("failed to parser mapping data: %s, err: %+v", data, err)
		} else {
			return mappingData, nil
		}
	}
}
