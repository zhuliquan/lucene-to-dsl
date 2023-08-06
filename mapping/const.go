package mapping

type FieldType string

var luceneSupportFieldType = map[FieldType]bool{
	ALIAS_FIELD_TYPE:            true,
	BINARY_FIELD_TYPE:           true,
	KEYWORD_FIELD_TYPE:          true,
	CONSTANT_KEYWORD_FIELD_TYPE: true,
	WILDCARD_FIELD_TYPE:         true,
	TEXT_FIELD_TYPE:             true,
	MATCH_ONLY_TEXT_FIELD_TYPE:  true,
	BOOLEAN_FIELD_TYPE:          true,
	BYTE_FIELD_TYPE:             true,
	SHORT_FIELD_TYPE:            true,
	INTEGER_FIELD_TYPE:          true,
	LONG_FIELD_TYPE:             true,
	UNSIGNED_LONG_FIELD_TYPE:    true,
	HALF_FLOAT_FIELD_TYPE:       true,
	FLOAT_FIELD_TYPE:            true,
	DOUBLE_FIELD_TYPE:           true,
	SCALED_FLOAT_FIELD_TYPE:     true,
	IP_FIELD_TYPE:               true,
	DATE_FIELD_TYPE:             true,
	IP_RANGE_FIELD_TYPE:         true,
	DATE_RANGE_FIELD_TYPE:       true,
	DATE_NANOS_FIELD_TYPE:       true,
	INTEGER_RANGE_FIELD_TYPE:    true,
	LONG_RANGE_FIELD_TYPE:       true,
	FLOAT_RANGE_FIELD_TYPE:      true,
	DOUBLE_RANGE_FIELD_TYPE:     true,
	OBJECT_FIELD_TYPE:           true,
	FLATTENED_FIELD_TYPE:        true,
	NESTED_FIELD_TYPE:           true,
	JOIN_FIELD_TYPE:             true,
}

const (
	UNKNOWN_FIELD_TYPE          FieldType = ""
	ALIAS_FIELD_TYPE            FieldType = "alias"            // alias
	BINARY_FIELD_TYPE           FieldType = "binary"           // base64 string
	KEYWORD_FIELD_TYPE          FieldType = "keyword"          // keyword
	CONSTANT_KEYWORD_FIELD_TYPE FieldType = "constant_keyword" // constant keyword
	WILDCARD_FIELD_TYPE         FieldType = "wildcard"         // wildcard
	TEXT_FIELD_TYPE             FieldType = "text"             // text
	MATCH_ONLY_TEXT_FIELD_TYPE  FieldType = "match_only_text"  // match_only_text
	VERSION_FIELD_TYPE          FieldType = "version"          // version， like 1.1.2
	BOOLEAN_FIELD_TYPE          FieldType = "boolean"          // true / false
	BYTE_FIELD_TYPE             FieldType = "byte"             // signed int8
	SHORT_FIELD_TYPE            FieldType = "short"            // signed int16
	INTEGER_FIELD_TYPE          FieldType = "integer"          // signed int32
	LONG_FIELD_TYPE             FieldType = "long"             // signed int64
	UNSIGNED_LONG_FIELD_TYPE    FieldType = "unsigned_long"    // unsigned int64
	HALF_FLOAT_FIELD_TYPE       FieldType = "half_float"       // float 16
	FLOAT_FIELD_TYPE            FieldType = "float"            // float 32
	DOUBLE_FIELD_TYPE           FieldType = "double"           // float 64
	SCALED_FLOAT_FIELD_TYPE     FieldType = "scaled_float"     // scaled float
	IP_FIELD_TYPE               FieldType = "ip"               // ipv4 / ipv6
	DATE_FIELD_TYPE             FieldType = "date"             // date
	DATE_NANOS_FIELD_TYPE       FieldType = "date_nanos"       // date_nanos
	IP_RANGE_FIELD_TYPE         FieldType = "ip_range"         // ip range
	DATE_RANGE_FIELD_TYPE       FieldType = "date_range"       // date range
	INTEGER_RANGE_FIELD_TYPE    FieldType = "integer_range"    // int32 range
	LONG_RANGE_FIELD_TYPE       FieldType = "long_range"       // int64 range
	FLOAT_RANGE_FIELD_TYPE      FieldType = "float_range"      // float32 range
	DOUBLE_RANGE_FIELD_TYPE     FieldType = "double_range"     // float64 range
	// properties 嵌套结构
	OBJECT_FIELD_TYPE    FieldType = "object"
	FLATTENED_FIELD_TYPE FieldType = "flattened" // flattened field
	NESTED_FIELD_TYPE    FieldType = "nested"    // nested field
	JOIN_FIELD_TYPE      FieldType = "join"      // join field

	// doesn't support by lucene
	DENSE_VECTOR_FIELD_TYPE       FieldType = "dense_vector"
	SPARSE_VECTOR_FIELD_TYPE      FieldType = "sparse_vector"
	RANK_FEATURE_FIELD_TYPE       FieldType = "rank_feature"
	RANK_FEATURES_FIELD_TYPE      FieldType = "rank_features"
	GEO_POINT_FIELD_TYPE          FieldType = "geo_point"
	GEO_SHAPE_FIELD_TYPE          FieldType = "geo_shape"
	POINT_FIELD_TYPE              FieldType = "point"
	SHAPE_FIELD_TYPE              FieldType = "shape"
	ANNOTATED_TEXT_FIELD_TYPE     FieldType = "annotated-text"
	COMPLETION_FIELD_TYPE         FieldType = "completion"
	SEARCH_AS_YOU_TYPE_FIELD_TYPE FieldType = "search_as_you_type"
	TOKEN_COUNT_FIELD_TYPE        FieldType = "token_count"

	// aggregate field type
	HISTOGRAM_FIELD_TYPE                FieldType = "histogram"               // used by
	AGGREGATE_METRICS_DOUBLE_FIELD_TYPE FieldType = "aggregate_metric_double" // used by (exists/range/term/terms) query
)

type IndexOptions string

const (
	UNKNOWN_INDEX_OPTIONS     IndexOptions = ""
	DOCS_INDEX_OPTIONS        IndexOptions = "docs"
	FREQUENCIES_INDEX_OPTIONS IndexOptions = "freqs"
	POSITIONS_INDEX_OPTIONS   IndexOptions = "positions"
	OFFSETS_INDEX_OPTIONS     IndexOptions = "offsets"
)

type MetricsType string

const (
	UNKNOWN_METRICS_TYPE     MetricsType = ""
	MAX_METRICS_TYPE         MetricsType = "max"
	MIN_METRICS_TYPE         MetricsType = "min"
	SUM_METRICS_TYPE         MetricsType = "sum"
	AVG_METRICS_TYPE         MetricsType = "avg"
	VALUE_COUNT_METRICS_TYPE MetricsType = "value_count"
)

type MetaUnitType string

const (
	UNKNOWN_META_UNIT_TYPE MetaUnitType = ""

	// float numeric
	PERCENT MetaUnitType = "%"

	// Whenever durations need to be specified,
	// e.g. for a timeout parameter, the duration must specify the unit,
	// like 2d for 2 days. The supported units are:
	YEAR        MetaUnitType = "y"
	MONTH       MetaUnitType = "M"
	Week        MetaUnitType = "w"
	DAY         MetaUnitType = "d"
	LOWER_HOUR  MetaUnitType = "h"
	UPPER_HOUR  MetaUnitType = "H"
	MINUTE      MetaUnitType = "m"
	SECOND      MetaUnitType = "s"
	MILLISECOND MetaUnitType = "ms"
	MICROSECOND MetaUnitType = "micros"
	NANOSECOND  MetaUnitType = "nanos"

	// Whenever the byte size of data needs to be specified,
	// e.g. when setting a buffer size parameter,
	// the value must specify the unit, like 10kb for 10 kilobytes.
	// Note that these units use powers of 1024, so 1kb means 1024 bytes. The supported units are:
	BYTE     MetaUnitType = "b"
	KILOBYTE MetaUnitType = "kb"
	MEGABYTE MetaUnitType = "mb"
	GIGABYTE MetaUnitType = "gb"
	TERABYTE MetaUnitType = "tb"
	PETABYTE MetaUnitType = "pb"

	// Unit-less quantities means that they don’t have a "unit" like "bytes" or "Hertz" or "meter" or "long tonne".
	// If one of these quantities is large we’ll print it out like 10m for 10,000,000 or 7k for 7,000.
	// We’ll still print 87 when we mean 87 though. These are the supported multipliers:
	BILO MetaUnitType = "k"
	MEGA MetaUnitType = "m"
	GIGA MetaUnitType = "g"
	TERA MetaUnitType = "t"
	PETA MetaUnitType = "p"

	// Wherever distances need to be specified, such as the distance parameter in the Geo-distance),
	// the default unit is meters if none is specified. Distances can be specified in other units,
	// such as "1km" or "2mi" (2 miles). The full list of units is listed below:
	MILE               MetaUnitType = "mi"
	MILE_FULL          MetaUnitType = "miles"
	YARD               MetaUnitType = "yd"
	YARD_FULL          MetaUnitType = "yards"
	FEET               MetaUnitType = "ft"
	FEET_FULL          MetaUnitType = "feet"
	INCH               MetaUnitType = "in"
	INCH_FULL          MetaUnitType = "inch"
	KILOMETER          MetaUnitType = "km"
	KILOMETER_FULL     MetaUnitType = "kilometers"
	METER              MetaUnitType = "m"
	METER_FULL         MetaUnitType = "meters"
	CENTIMETER         MetaUnitType = "cm"
	CENTIMETER_FULL    MetaUnitType = "centimeters"
	MILLIMETER         MetaUnitType = "mm"
	MILLIMETER_FULL    MetaUnitType = "millimeters"
	NAUTICAL_MILE      MetaUnitType = "NM"
	NAUTICAL_MILE_1    MetaUnitType = "nmi"
	NAUTICAL_MILE_FULL MetaUnitType = "nauticalmiles"
)

type MetaMetricsType string

const (
	UNKNOWN_META_METRICS_TYPE MetaMetricsType = ""

	GAUGE   MetaMetricsType = "gauge"
	COUNTER MetaMetricsType = "counter"
	SUMMARY MetaMetricsType = "summary"
)

type MappingType uint8

type Dynamic interface {
	GetMappingType() MappingType
}

const (
	// New fields are added to the mapping (default).
	DYNAMIC_MAPPING MappingType = 1 // true
	// New fields are ignored. These fields will not be indexed or searchable,
	// but will still appear in the _source field of returned hits.
	// These fields will not be added to the mapping, and new fields must be added explicitly.
	STATIC_MAPPING MappingType = 2 // false
	// If new fields are detected, an exception is thrown and the document is rejected.
	// New fields must be explicitly added to the mapping.
	STRICT_MAPPING MappingType = 3 // strict
	// New fields are added to the mapping as runtime fields.
	// These fields are not indexed, and are loaded from _source at query time.
	RUNTIME_MAPPING MappingType = 4 // runtime
)

var MappingTypeString = map[MappingType]string{
	DYNAMIC_MAPPING: "true",
	STATIC_MAPPING:  "false",
	STRICT_MAPPING:  "strict",
	RUNTIME_MAPPING: "runtime",
}

type BoolDynamic bool 

func (b BoolDynamic) GetMappingType() MappingType {
	if b {
		return DYNAMIC_MAPPING
	}
	return STATIC_MAPPING
}

type StringDynamic string


func (s StringDynamic) GetMappingType() MappingType {
	switch s {
	case "true":
		return DYNAMIC_MAPPING
	case "false":
		return STATIC_MAPPING
	case "strict":
		return STRICT_MAPPING
	case "runtime":
		return RUNTIME_MAPPING
	default:
		return DYNAMIC_MAPPING				
	}
}

type Similarity string

const (
	UNKNOWN_SIMILARITY  Similarity = ""
	BM25_SIMILARITY     Similarity = "BM25"
	CLASSSIC_SIMILARITY Similarity = "classic"
	BOOLEAN_SIMILARITY  Similarity = "boolean"
)

type TermVector string

const (
	UNKNOWN_TERM_VECTOR TermVector = ""

	// No term vectors are stored. (default)
	NO_TERM_VECTOR TermVector = "no"

	// Just the terms in the field are stored.
	YES_TERM_VECTOR TermVector = "yes"

	// Terms and positions are stored.
	WITH_POSITIONS TermVector = "with_positions"

	// Terms and character offsets are stored.
	WITH_OFFSETS TermVector = "with_offsets"

	// Terms, positions, and character offsets are stored.
	WITH_POSITIONS_OFFSETS TermVector = "with_positions_offsets"

	// Terms, positions, and payloads are stored.
	WITH_POSITIONS_PAYLOADS TermVector = "with_positions_payloads"

	// Terms, positions, offsets and payloads are stored.
	WITH_POSITIONS_OFFSETS_PAYLOADS TermVector = "with_positions_offsets_payloads"
)
