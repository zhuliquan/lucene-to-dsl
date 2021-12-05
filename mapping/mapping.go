package mapping

type Meta struct {
	Unit        MetaUnitType    `json:"unit"`
	MetricsType MetaMetricsType `json:"metrics_type"`
}

type Source struct {
	Enabled  bool     `json:"enabled"`
	Includes []string `json:"includes"`
	Excludes []string `json:"excludes"`
}

type All struct {
	Enable bool `json:"anabled"`
	Store  bool `json:"store"`
}

type FieldMapping struct {
	Mapping
	Type      FieldType `json:"type"`
	Boost     float64   `json:"boost"`
	Index     bool      `json:"index,omitempty"`      // Should the field be searchable? Accepts true (default) and false.
	Store     bool      `json:"store"`                // whether save field value, so it can return this field from original document.
	NullValue string    `json:"null_value,omitempty"` // using for check whether field exist.
	Path      string    `json:"path,omitempty"`       // using be alias type

	Meta                *Meta        `json:"meta,omitempty"`       // meta data for numeric type, such as long, float
	IndexOptions        IndexOptions `json:"index_options"`        // index value will be saved
	DocValues           bool         `json:"doc_values,omitempty"` // take advantage for field aggregating and field sorting
	IncludeInAll        bool         `json:"include_in_all,omitempty"`
	DepthLimit          int          `json:"depth_limit"`
	EagerGlobalOrdinals bool         `json:"eager_global_ordinals"`
	IgnoreAbove         string       `json:"ignore_above"`
	IgnoreMalformed     bool         `json:"ignore_malformed"`

	// 自定义
	UpperCase bool `json:"upper_case,omitempty"`
	LowerCase bool `json:"lower_case,omitempty"`
}

type Mapping struct {
	MappingType  MappingType              `json:"dynamic"`
	IncludeInAll bool                     `json:"include_in_all"`
	Source       *Source                  `json:"_source"`
	All          *All                     `json:"_all"`
	Properties   map[string]*FieldMapping `json:"properties"`
}
