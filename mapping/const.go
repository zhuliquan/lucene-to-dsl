package mapping

type FieldType int32

const (
	BINARY_FIELD_TYPE         FieldType = iota // base64
	BOOLEAN_FIELD_TYPE                         // true / false
	BYTE_FIELD_TYPE                            // int8
	SHORT_FIELD_TYPE                           // int16
	INTEGER_FIELD_TYPE                         // int32
	LONG_FIELD_TYPE                            // int64
	HALF_FLOAT_FIELD_TYPE                      // float 16
	FLOAT_FIELD_TYPE                           // float 32
	DOUBLE_FIELD_TYPE                          // float 64
	SCALED_FLOAT_FIELD_TYPE                    // scaled float
	IP_FIELD_TYPE                              // ipv4 / ipv6
	DATE_FIELD_TYPE                            // date field
	INTERGER_RANGE_FIELD_TYPE                  // int32
	LONG_RANGE_FIELD_TYPE                      // int64
	FLOAT_RANGE_FIELD_TYPE                     // float32
	DOUBLE_RANGE_FIELD_TYPE                    // float64
	DATE_RANGE_FIELD_TYPE                      // date
	ALIAS_FIELD_TYPE                           // alias for exists field
	OBJECT_FIELD_TYPE                          // json strcut

	FLATTENED_FIELD_TYPE // flattened field
	NESTED_FIELD_TYPE    // nested field
	JOIN_FIELD_TYPE      // join field
)
