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
	RANGE_DSL_TYPE
	MATCH_DSL_TYPE
	TERM_DSL_TYPE
	EXIST_DSL_TYPE
	MATCH_PHRASE_DSL_TYPE
	REGEXP_DSL_TYPE
	FUZZY_DSL_TYPE
)
