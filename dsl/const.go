package dsl

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

type DSLTermType uint32

const (
	KEYWORD_VALUE DSLTermType = iota
	PHRASE_VALUE
	INT_VALUE
	FLOAT_VALUE
	IP_VALUE
	IP_CIDR_VALUE
	DATE_VALUE
)

// using nil represent infinite value
var InfValue *DSLTermValue = nil
