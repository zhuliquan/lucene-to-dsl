package dsl

import (
	"regexp"
)

type RegexpNode struct {
	kvNode
	rewriteNode
	statesNode
	patternNode
	flags RegexpFlagType
}

func WithFlags(flags RegexpFlagType) func(AstNode) {
	return func(n AstNode) {
		if r, ok := n.(*RegexpNode); ok {
			r.flags = flags
		}
	}
}

func NewRegexpNode(kvNode *kvNode, pattern *regexp.Regexp, opts ...func(AstNode)) *RegexpNode {
	var n = &RegexpNode{
		kvNode:      *kvNode,
		rewriteNode: rewriteNode{rewrite: CONSTANT_SCORE},
		statesNode:  statesNode{maxDeterminizedStates: 10000},
		patternNode: patternNode{matcher: pattern},
		flags:       ALL_FLAG,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *RegexpNode) DslType() DslType {
	return REGEXP_DSL_TYPE
}

func (n *RegexpNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE, BOOL_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		return patternNodeUnionJoinTermNode(n, o.(*TermNode))
	case REGEXP_DSL_TYPE:
		return valueNodeUnionJoinValueNode(n, o)
	default:
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func (n *RegexpNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE, BOOL_DSL_TYPE:
		return o.InterSect(n)
	case TERM_DSL_TYPE:
		return patternNodeIntersectTermNode(n, o.(*TermNode))
	case REGEXP_DSL_TYPE:
		return valueNodeIntersectValueNode(n, o)
	default:
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	}
}

func (n *RegexpNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}

func (n *RegexpNode) ToDSL() DSL {
	return DSL{
		REGEXP_KEY: DSL{
			n.field: DSL{
				VALUE_KEY:   n.toPrintValue(),
				REWRITE_KEY: n.getRewrite(),
				FLAGS_KEY:   n.flags,

				MAX_DETERMINIZED_STATES_KEY: n.getMaxDeterminizedStates(),
			},
		},
	}
}
