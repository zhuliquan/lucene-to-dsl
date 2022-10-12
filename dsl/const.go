package dsl

import (
	"math"
	"net"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/x448/float16"
	"github.com/zhuliquan/scaled_float"
)

type AstType uint32

const (
	EMPTY_NODE_TYPE AstType = iota
	OP_NODE_TYPE
	LEAF_NODE_TYPE
)

type DslType uint32

const (
	EMPTY_DSL_TYPE DslType = iota
	AND_DSL_TYPE
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
	eps        = 1E-8
	MinUint    = uint64(0)
	MinFloat16 = float16.Fromfloat32(-65504)
	MaxFloat16 = float16.Fromfloat32(65504)

	MinScaledFloat = scaled_float.NegativeInf
	MaxScaledFloat = scaled_float.PositiveInf

	MaxIP = net.IP([]byte{
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
		math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
	})
	MinIP = net.IP([]byte{0, 0, 0, 0})

	MinTime       = time.Unix(math.MinInt64, 0)
	MaxTime       = time.Unix(math.MaxInt64, 999999999)
	MinVersion, _ = version.NewVersion("v0-A.0-A.0-A")
	MaxVersion, _ = version.NewVersion("v9223372036854775807.9223372036854775807.9223372036854775807")
	MinString     = ""
	MaxString     = string([]rune{math.MaxInt32})
)

var MinInt = map[int]int64{
	8:  int64(math.MinInt8),
	16: int64(math.MinInt16),
	32: int64(math.MinInt32),
	64: int64(math.MinInt64),
}
var MaxInt = map[int]int64{
	8:  int64(math.MaxInt8),
	16: int64(math.MaxInt16),
	32: int64(math.MaxInt32),
	64: int64(math.MaxInt64),
}
var MaxUint = map[int]uint64{
	8:  uint64(math.MaxUint8),
	16: uint64(math.MaxUint16),
	32: uint64(math.MaxUint32),
	64: uint64(math.MaxUint64),
}

var MinFloat = map[int]interface{}{
	16:  MinFloat16,
	32:  -math.MaxFloat32,
	64:  -math.MaxFloat64,
	128: MinScaledFloat,
}

var MaxFloat = map[int]interface{}{
	16:  MaxFloat16,
	32:  math.MaxFloat32,
	64:  math.MaxFloat64,
	128: MaxScaledFloat,
}

type CompareType uint32

const (
	EQ CompareType = iota
	LT
	GT
	LTE
	GTE
)

var compareTypeStrings = map[CompareType]string{
	EQ:  "eq",
	LT:  "lt",
	GT:  "gt",
	LTE: "lte",
	GTE: "gte",
}

func (c CompareType) String() string {
	return compareTypeStrings[c]
}

const (
	// Uses the constant_score_boolean method for fewer matching terms.
	// Otherwise, this method finds all matching terms in sequence and returns matching documents using a bit set.
	CONSTANT_SCORE = "constant_score" // default

	// Assigns each document a relevance score equal to the boost parameter.
	// This method changes the original query to a bool query.
	// This bool query contains a should clause and term query for each matching term.
	// This method can cause the final bool query to exceed the clause limit in the `indices.query.bool.max_clause_count` setting.
	// If the query exceeds this limit, Elasticsearch returns an error.
	CONSTANT_SCORE_BOOLEAN = "constant_score_boolean"

	// Calculates a relevance score for each matching document.
	// This method changes the original query to a bool query.
	// This bool query contains a should clause and term query for each matching term.
	// This method can cause the final bool query to exceed the clause limit in the indices.query.bool.max_clause_count setting.
	// If the query exceeds this limit, Elasticsearch returns an error.
	SCORING_BOOLEAN = "scoring_boolean"

	// Calculates a relevance score for each matching document as if all terms had the same frequency.
	// This frequency is the maximum frequency of all matching terms.
	// This method changes the original query to a bool query.
	// This bool query contains a should clause and term query for each matching term.
	// The final bool query only includes term queries for the top N scoring terms.
	// You can use this method to avoid exceeding the clause limit in the `indices.query.bool.max_clause_count` setting.
	TOP_TERMS_BLENDED_FREQS_N = "top_terms_blended_freqs_N"

	// Assigns each matching document a relevance score equal to the boost parameter.
	// This method changes the original query to a bool query.
	// This bool query contains a should clause and term query for each matching term.
	// The final bool query only includes term queries for the top N terms.
	// You can use this method to avoid exceeding the clause limit in the `indices.query.bool.max_clause_count` setting.
	TOP_TERMS_BOOST_N = "top_terms_boost_N"

	// Calculates a relevance score for each matching document.
	// This method changes the original query to a bool query.
	// This bool query contains a should clause and term query for each matching term.
	// The final bool query only includes term queries for the top N scoring terms.
	// You can use this method to avoid exceeding the clause limit in the `indices.query.bool.max_clause_count` setting.
	TOP_TERMS_N = "top_terms_N"
)

const (
	// Matches documents with a range field value that intersects the query’s range.
	INTERSECTS = "INTERSECTS" // default
	// Matches documents with a range field value that entirely contains the query’s range.
	CONTAINS = "CONTAINS"
	// Matches documents with a range field value entirely within the query’s range.
	WITHIN = "WITHIN"
)
