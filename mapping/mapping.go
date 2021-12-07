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
	// Strings longer than the ignore_above setting will not be indexed or stored.
	// For arrays of strings, ignore_above will be applied for each array element
	// separately and string elements longer than ignore_above will not be indexed or stored.
	// All strings/array elements will still be present in the _source field,
	// if the latter is enabled which is the default in Elasticsearch.
	IgnoreAbove string `json:"ignore_above,omitempty"`

	// Sometimes you don’t have much control over the data that you receive.
	// One user may send a login field that is a date, and another sends a login field that is an email address.
	// Trying to index the wrong data type into a field throws an exception by default,
	// and rejects the whole document. The ignore_malformed parameter,
	// if set to true, allows the exception to be ignored. The malformed field is not indexed,
	// but other fields in the document are processed normally.
	IgnoreMalformed bool `json:"ignore_malformed,omitempty"`

	// In JSON documents, dates are represented as strings.
	// Elasticsearch uses a set of preconfigured formats to recognize and parse these strings into a long value representing milliseconds-since-the-epoch in UTC.
	// Besides the built-in formats, your own custom formats can be specified using the familiar yyyy/MM/dd syntax:
	Format TimeFormat `json:"format,omitempty"`

	// 	Elasticsearch tries to index all of the fields you give it,
	// but sometimes you want to just store the field without indexing it.
	// For instance, imagine that you are using Elasticsearch as a web session store.
	// You may want to index the session ID and last update time,
	//  but you don’t need to query or run aggregations on the session data itself.
	// The enabled setting, which can be applied only to the top-level mapping definition and to object fields,
	// causes Elasticsearch to skip parsing of the contents of the field entirely. The JSON can still be retrieved from the _source field, but it is not searchable or stored in any other way:
	Enabled bool `json:"enabled"`

	Coerce bool `json:"coerce,omitempty"`

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
