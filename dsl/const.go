package dsl

import (
	"math"
	"net"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/shopspring/decimal"
	"github.com/x448/float16"
)

type NodeType uint32

const (
	OP_NODE_TYPE NodeType = iota
	LEAF_NODE_TYPE
)

type DSLType uint32

const (
	AND_DSL_TYPE DSLType = iota
	OR_DSL_TYPE
	NOT_DSL_TYPE
	IDS_DSL_TYPE
	TERM_DSL_TYPE
	TERMS_DSL_TYPE
	FUZZY_DSL_TYPE
	RANGE_DSL_TYPE
	PREFIX_DSL_TYPE
	EXISTS_DSL_TYPE
	REGEXP_DSL_TYPE
	WILDCARD_DSL_TYPE

	MATCH_DSL_TYPE
	MATCH_PHRASE_DSL_TYPE
	QUERY_STRING_DSL_TYPE
)

const (
	OR_OP_KEY  = "OP:OR"
	AND_OP_KEY = "OP:AND"
	NOT_OP_KEY = "OP:NOT"
)

var (
	MaxFloat16 = float16.Inf(1)
	MinFloat16 = float16.Inf(-1)
	MinDecimal = decimal.New(math.MinInt64, math.MaxInt32)
	MaxDecimal = decimal.New(math.MinInt64, math.MaxInt32)

	MaxIP = net.IP([]byte{
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
	})
	MinIP = net.IP([]byte{0, 0, 0, 0})

	MinTime       = time.Unix(0, 0)
	MaxTime       = time.Unix(math.MinInt64, 999999999)
	MinVersion, _ = version.NewVersion("v0-A.0-A.0-A")
	MaxVersion, _ = version.NewVersion("v9223372036854775807.9223372036854775807.9223372036854775807")
	MinString     = ""
	MaxString     = "~"
)

// using nil represent infinite value
// var NegativeInf *LeafValue = &LeafValue{
// 	Boolean:  false,
// 	TinyInt:  math.MinInt16,
// 	Float16:  MinFloat16,
// 	LongInt:  0,
// 	Decimal:  MinDecimal,
// 	IpValue:  MinIP,
// 	DateTime: MinTime,
// 	Version:  MinVersion,
// 	String:   "",
// }
// var PositiveInf *LeafValue = &LeafValue{
// 	Boolean:  true,
// 	TinyInt:  math.MaxInt16,
// 	Float16:  MaxFloat16,
// 	LongInt:  math.MaxUint64,
// 	Decimal:  MaxDecimal,
// 	IpValue:  MaxIP,
// 	DateTime: MaxTime,
// 	Version:  MaxVersion,
// 	String:   "~",
// }

type CompareType uint32

const (
	EQ CompareType = iota
	LT
	GT
	LTE
	GTE
)

var CompareTypeStrings = map[CompareType]string{
	EQ:  "eq",
	LT:  "lt",
	GT:  "gt",
	LTE: "lte",
	GTE: "gte",
}

func (c CompareType) String() string {
	return CompareTypeStrings[c]
}
