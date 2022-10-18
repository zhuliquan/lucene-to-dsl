package dsl

import (
	"regexp"
)

type RegexpNode struct {
	kvNode
	rewriteNode
	statesNode
	patternNode
	flags string
}

func NewRegexNode(kvNode *kvNode, pattern *regexp.Regexp, opts ...func(AstNode)) *RegexpNode {
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
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		return patternNodeUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return patternNodeUnionJoinTermsNode(n, o.(*TermsNode))
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *RegexpNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	case TERM_DSL_TYPE:
		return patternNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return patternNodeIntersectTermsNode(n, o.(*TermsNode))
	default:
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *RegexpNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
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
