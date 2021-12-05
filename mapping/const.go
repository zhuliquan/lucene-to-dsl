package mapping

type FieldType string

const (
	UNKNOWN_FIELD_TYPE          FieldType = ""
	BINARY_FIELD_TYPE           FieldType = "binary"           // base64 string
	KEYWORD_FIELD_TYPE          FieldType = "keyword"          // keyword
	CONSTANT_KEYWORD_FIELD_TYPE FieldType = "constant_keyword" // constant keyword
	WILDCARD_FIELD_TYPE         FieldType = "wildcard"         // wildcard
	TEXT_FIELD_TYPE             FieldType = "text"             // text
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
	IP_RANGE_FIELD_TYPE         FieldType = "ip_range"         // ip range
	DATE_RANGE_FIELD_TYPE       FieldType = "date_range"       // date range
	INTERGER_RANGE_FIELD_TYPE   FieldType = "integer_range"    // int32 range
	LONG_RANGE_FIELD_TYPE       FieldType = "long_range"       // int64 range
	FLOAT_RANGE_FIELD_TYPE      FieldType = "float_range"      // float32 range
	DOUBLE_RANGE_FIELD_TYPE     FieldType = "double_range"     // float64 range
	ALIAS_FIELD_TYPE            FieldType = "alias"            // alias for exists field
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
	MAX_METRICS_TYPE         MetricsType = "max"
	MIN_METRICS_TYPE         MetricsType = "min"
	SUM_METRICS_TYPE         MetricsType = "sum"
	AVG_METRICS_TYPE         MetricsType = "avg"
	VALUE_COUNT_METRICS_TYPE MetricsType = "value_count"
)
