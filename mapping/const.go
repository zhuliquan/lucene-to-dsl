package mapping

type FieldType string

var LuceneSupportFieldType = []FieldType{
	BINARY_FIELD_TYPE,
	KEYWORD_FIELD_TYPE,
	CONSTANT_KEYWORD_FIELD_TYPE,
	WILDCARD_FIELD_TYPE,
	TEXT_FIELD_TYPE,
	BOOLEAN_FIELD_TYPE,
	BYTE_FIELD_TYPE,
	SHORT_FIELD_TYPE,
	INTEGER_FIELD_TYPE,
	LONG_FIELD_TYPE,
	UNSIGNED_LONG_FIELD_TYPE,
	HALF_FLOAT_FIELD_TYPE,
	FLOAT_FIELD_TYPE,
	DOUBLE_FIELD_TYPE,
	SCALED_FLOAT_FIELD_TYPE,
	IP_FIELD_TYPE,
	DATE_FIELD_TYPE,
	IP_RANGE_FIELD_TYPE,
	DATE_RANGE_FIELD_TYPE,
	INTERGER_RANGE_FIELD_TYPE,
	LONG_RANGE_FIELD_TYPE,
	FLOAT_RANGE_FIELD_TYPE,
	DOUBLE_RANGE_FIELD_TYPE,
	ALIAS_FIELD_TYPE,
	OBJECT_FIELD_TYPE,
	FLATTENED_FIELD_TYPE,
	NESTED_FIELD_TYPE,
	JOIN_FIELD_TYPE,
}

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

type MappingType string

const (
	// New fields are added to the mapping (default).
	DYNAMIC_MAPPING MappingType = "true"
	// New fields are ignored. These fields will not be indexed or searchable,
	// but will still appear in the _source field of returned hits.
	// These fields will not be added to the mapping, and new fields must be added explicitly.
	STATIC_MAPPING MappingType = "false"
	// If new fields are detected, an exception is thrown and the document is rejected.
	// New fields must be explicitly added to the mapping.
	STRICT_MAPPING MappingType = "strict"
	// New fields are added to the mapping as runtime fields.
	// These fields are not indexed, and are loaded from _source at query time.
	RUNTIME_MAPPING MappingType = "runtime"
)

type TimeFormat string

const (
	// A formatter for the number of milliseconds since the epoch.
	// Note, that this timestamp is subject to the limits of a Java Long.MIN_VALUE and Long.MAX_VALUE.
	EPOCH_MILLIS TimeFormat = "epoch_millis"

	// A formatter for the number of seconds since the epoch.
	// Note, that this timestamp is subject to the limits of a Java Long.MIN_VALUE and Long.
	// MAX_VALUE divided by 1000 (the number of milliseconds in a second).
	// date_optional_time or strict_date_optional_time
	// A generic ISO datetime parser, where the date must include the year at a minimum,
	// and the time (separated by T), is optional. Examples: yyyy-MM-dd'T'HH:mm:ss.SSSZ or yyyy-MM-dd.
	EPOCH_SECOND TimeFormat = "epoch_second"

	// A generic ISO datetime parser, where the date must include the year at a minimum, and the time (separated by T), is optional. The fraction of a second part has a nanosecond resolution. Examples: yyyy-MM-dd'T'HH:mm:ss.SSSSSSZ or yyyy-MM-dd.
	STRICT_DATE_OPTIONAL_TIME_NANOS TimeFormat = "strict_date_optional_time_nanos"

	// A basic formatter for a full date as four digit year, two digit month of year, and two digit day of month: yyyyMMdd.
	BASIC_DATE TimeFormat = "basic_date"

	// A basic formatter that combines a basic date and time, separated by a T: yyyyMMdd'T'HHmmss.SSSZ.
	BASIC_DATE_TIME TimeFormat = "basic_date_time"

	// A basic formatter that combines a basic date and time without millis, separated by a T: yyyyMMdd'T'HHmmssZ.
	BASIC_DATE_TIME_NO_MILLIS TimeFormat = "basic_date_time_no_millis"

	// A formatter for a full ordinal date, using a four digit year and three digit dayOfYear: yyyyDDD.
	BASIC_ORDINAL_DATE TimeFormat = "basic_ordinal_date"

	// A formatter for a full ordinal date and time, using a four digit year and three digit dayOfYear: yyyyDDD'T'HHmmss.SSSZ.
	BASIC_ORDINAL_DATE_TIME TimeFormat = "basic_ordinal_date_time"

	// A formatter for a full ordinal date and time without millis, using a four digit year and three digit dayOfYear: yyyyDDD'T'HHmmssZ.
	BASIC_ORDINAL_DATE_TIME_NO_MILLIS TimeFormat = "basic_ordinal_date_time_no_millis"

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit millis, and time zone offset: HHmmss.SSSZ.
	BASIC_TIME TimeFormat = "basic_time"

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset: HHmmssZ.
	BASIC_TIME_NO_MILLIS TimeFormat = "basic_time_no_millis"

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit millis, and time zone off set prefixed by T: 'T'HHmmss.SSSZ.
	BASIC_T_TIME TimeFormat = "basic_t_time"

	// A basic formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset prefixed by T: 'T'HHmmssZ.
	BASIC_T_TIME_NO_MILLIS TimeFormat = "basic_t_time_no_millis"

	// A basic formatter for a full date as four digit weekyear, two digit week of weekyear, and one digit day of week: xxxx'W'wwe.
	BASIC_WEEK_DATE        TimeFormat = "basic_week_date"
	STRICT_BASIC_WEEK_DATE TimeFormat = "strict_basic_week_date"

	// A basic formatter that combines a basic weekyear date and time, separated by a T: xxxx'W'wwe'T'HHmmss.SSSZ.
	BASIC_WEEK_DATE_TIME        TimeFormat = "basic_week_date_time"
	STRICT_BASIC_WEEK_DATE_TIME TimeFormat = "strict_basic_week_date_time"

	// A basic formatter that combines a basic weekyear date and time without millis, separated by a T: xxxx'W'wwe'T'HHmmssZ.
	BASIC_WEEK_DATE_TIME_NO_MILLIS        TimeFormat = "basic_week_date_time_no_millis"
	STRICT_BASIC_WEEK_DATE_TIME_NO_MILLIS TimeFormat = "strict_basic_week_date_time_no_millis"
	// A formatter for a full date as four digit year, two digit month of year, and two digit day of month: yyyy-MM-dd.
	DATE        TimeFormat = "date"
	STRICT_DATE TimeFormat = "strict_date"
	// A formatter that combines a full date and two digit hour of day: yyyy-MM-dd'T'HH.
	DATE_HOUR        TimeFormat = "date_hour"
	STRICT_DATE_HOUR TimeFormat = "strict_date_hour"

	// A formatter that combines a full date, two digit hour of day, and two digit minute of hour: yyyy-MM-dd'T'HH:mm.
	DATE_HOUR_MINUTE        TimeFormat = "date_hour_minute"
	STRICT_DATE_HOUR_MINUTE TimeFormat = "strict_date_hour_minute"

	// A formatter that combines a full date, two digit hour of day, two digit minute of hour, and two digit second of minute: yyyy-MM-dd'T'HH:mm:ss.
	DATE_HOUR_MINUTE_SECOND        TimeFormat = "date_hour_minute_second"
	STRICT_DATE_HOUR_MINUTE_SECOND TimeFormat = "strict_date_hour_minute_second"

	// A formatter that combines a full date, two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: yyyy-MM-dd'T'HH:mm:ss.SSS.
	DATE_HOUR_MINUTE_SECOND_FRACTION        TimeFormat = "date_hour_minute_second_fraction"
	STRICT_DATE_HOUR_MINUTE_SECOND_FRACTION TimeFormat = "strict_date_hour_minute_second_fraction"

	// A formatter that combines a full date, two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: yyyy-MM-dd'T'HH:mm:ss.SSS.
	DATE_HOUR_MINUTE_SECOND_MILLIS        TimeFormat = "date_hour_minute_second_millis"
	STRICT_DATE_HOUR_MINUTE_SECOND_MILLIS TimeFormat = "strict_date_hour_minute_second_millis"

	// A formatter that combines a full date and time, separated by a T: yyyy-MM-dd'T'HH:mm:ss.SSSZZ.
	DATE_TIME        TimeFormat = "date_time"
	STRICT_DATE_TIME TimeFormat = "strict_date_time"

	// A formatter that combines a full date and time without millis, separated by a T: yyyy-MM-dd'T'HH:mm:ssZZ.
	DATE_TIME_NO_MILLIS        TimeFormat = "date_time_no_millis"
	STRICT_DATE_TIME_NO_MILLIS TimeFormat = "strict_date_time_no_millis"

	// A formatter for a two digit hour of day: HH
	HOUR        TimeFormat = "hour"
	STRICT_HOUR TimeFormat = "strict_hour"

	// A formatter for a two digit hour of day and two digit minute of hour: HH:mm.
	HOUR_MINUTE        TimeFormat = "hour_minute"
	STRICT_HOUR_MINUTE TimeFormat = "strict_hour_minute"

	// A formatter for a two digit hour of day, two digit minute of hour, and two digit second of minute: HH:mm:ss.
	HOUR_MINUTE_SECOND        TimeFormat = "hour_minute_second"
	STRICT_HOUR_MINUTE_SECOND TimeFormat = "strict_hour_minute_second"

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: HH:mm:ss.SSS.
	HOUR_MINUTE_SECOND_FRACTION        TimeFormat = "hour_minute_second_fraction"
	STRICT_HOUR_MINUTE_SECOND_FRACTION TimeFormat = "strict_hour_minute_second_fraction"

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and three digit fraction of second: HH:mm:ss.SSS.
	HOUR_MINUTE_SECOND_MILLIS        TimeFormat = "hour_minute_second_millis"
	STRICT_HOUR_MINUTE_SECOND_MILLIS TimeFormat = "strict_hour_minute_second_millis"

	// A formatter for a full ordinal date, using a four digit year and three digit dayOfYear: yyyy-DDD.
	ORDINAL_DATE        TimeFormat = "ordinal_date"
	STRICT_ORDINAL_DATE TimeFormat = "strict_ordinal_date"

	// A formatter for a full ordinal date and time, using a four digit year and three digit dayOfYear: yyyy-DDD'T'HH:mm:ss.SSSZZ.
	ORDINAL_DATE_TIME        TimeFormat = "ordinal_date_time"
	STRICT_ORDINAL_DATE_TIME TimeFormat = "strict_ordinal_date_time"

	// A formatter for a full ordinal date and time without millis, using a four digit year and three digit dayOfYear: yyyy-DDD'T'HH:mm:ssZZ.
	ORDINAL_DATE_TIME_NO_MILLIS        TimeFormat = "ordinal_date_time_no_millis"
	STRICT_ORDINAL_DATE_TIME_NO_MILLIS TimeFormat = "strict_ordinal_date_time_no_millis"

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit fraction of second, and time zone offset: HH:mm:ss.SSSZZ.
	TIME        TimeFormat = "time"
	STRICT_TIME TimeFormat = "strict_time"

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset: HH:mm:ssZZ.
	TIME_NO_MILLIS        TimeFormat = "time_no_millis"
	STRICT_TIME_NO_MILLIS TimeFormat = "strict_time_no_millis"

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, three digit fraction of second, and time zone offset prefixed by T: 'T'HH:mm:ss.SSSZZ.
	T_TIME        TimeFormat = "t_time"
	STRICT_T_TIME TimeFormat = "strict_t_time"

	// A formatter for a two digit hour of day, two digit minute of hour, two digit second of minute, and time zone offset prefixed by T: 'T'HH:mm:ssZZ.
	T_TIME_NO_MILLIS        TimeFormat = "t_time_no_millis"
	STRICT_T_TIME_NO_MILLIS TimeFormat = "t_time_no_millis"

	// A formatter for a full date as four digit weekyear, two digit week of weekyear, and one digit day of week: xxxx-'W'ww-e.
	WEEK_DATE        TimeFormat = "week_date"
	STRICT_WEEK_DATE TimeFormat = "strict_week_date"

	// A formatter that combines a full weekyear date and time, separated by a T: xxxx-'W'ww-e'T'HH:mm:ss.SSSZZ.
	WEEK_DATE_TIME        TimeFormat = "strict_week_date"
	STRICT_WEEK_DATE_TIME TimeFormat = "strict_week_date_time"

	// A formatter that combines a full weekyear date and time without millis, separated by a T: xxxx-'W'ww-e'T'HH:mm:ssZZ.
	WEEK_DATE_TIME_NO_MILLIS        TimeFormat = "strict_week_date_time"
	STRICT_WEEK_DATE_TIME_NO_MILLIS TimeFormat = "strict_week_date_time_no_millis"

	// A formatter for a four digit weekyear: xxxx.
	WEEKYEAR        TimeFormat = "weekyear"
	STRICT_WEEKYEAR TimeFormat = "strict_weekyear"

	// A formatter for a four digit weekyear and two digit week of weekyear: xxxx-'W'ww.
	WEEKYEAR_WEEK        TimeFormat = "weekyear_week"
	STRICT_WEEKYEAR_WEEK TimeFormat = "strict_weekyear_week"

	// A formatter for a four digit weekyear, two digit week of weekyear, and one digit day of week: xxxx-'W'ww-e.
	WEEKYEAR_WEEK_DAY        TimeFormat = "weekyear_week_day"
	STRICT_WEEKYEAR_WEEK_DAY TimeFormat = "strict_weekyear_week_day"

	// A formatter for a four digit year and two digit month of year: yyyy-MM.
	YEAR_MONTH        TimeFormat = "year_month"
	STRICT_YEAR_MONTH TimeFormat = "strict_year_month"

	// A formatter for a four digit year: yyyy.
	YEAR_FORMAT TimeFormat = "year"
	STRICT_YEAR TimeFormat = "strict_year"

	// A formatter for a four digit year, two digit month of year, and two digit day of month: yyyy-MM-dd.
	YEAR_MONTH_DAY        TimeFormat = "year_month_day"
	STRICT_YEAR_MONTH_DAY TimeFormat = "strict_year_month_day"
)

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
