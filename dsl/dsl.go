package dsl

import (
	"encoding/json"
	"net"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/shopspring/decimal"
	"github.com/x448/float16"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

type DSL map[string]interface{}

func (d DSL) String() string {
	v, _ := json.Marshal(d)
	return string(v)
}

// type LeafValue interface{}

// type Boolean bool
// type TinyInt byte
// type ShortInt int16
// type Integer int32
// type Unsigned uint64

type LeafValue struct {
	Boolean  bool
	TinyInt  int16
	Float16  float16.Float16
	LongInt  uint64
	Decimal  decimal.Decimal
	Version  *version.Version
	IpValue  net.IP
	DateTime time.Time
	String   string
}

func (v *LeafValue) Value(typ mapping.FieldType) interface{} {
	switch typ {
	case mapping.BOOLEAN_FIELD_TYPE:
		return v.Boolean
	case mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE:
		return v.TinyInt
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		return v.LongInt
	case mapping.HALF_FLOAT_FIELD_TYPE:
		return v.Float16
	case mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE, mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE:
		return v.Decimal.BigInt().Int64()
	case mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE, mapping.DOUBLE_FIELD_TYPE, mapping.DOUBLE_RANGE_FIELD_TYPE:
		f, _ := v.Decimal.BigFloat().Float64()
		return f
	case mapping.SCALED_FLOAT_FIELD_TYPE:
		return v.Decimal.BigFloat()
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		return v.IpValue.String()
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
		return v.DateTime.Unix()
	case mapping.VERSION_FIELD_TYPE:
		return v.Version.String()
	default:
		return v.String

	}
}

var EmptyDSL = DSL{}
