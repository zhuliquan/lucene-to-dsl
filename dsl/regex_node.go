package dsl

import (
	"fmt"
	"regexp"
)

type RegexpNode struct {
	kvNode
	boostNode
	rewriteNode
	statesNode

	pattern *regexp.Regexp
}

func NewRegexNode(kvNode *kvNode, pattern *regexp.Regexp, opts ...func(AstNode)) *RegexpNode {
	var n = &RegexpNode{
		kvNode:      *kvNode,
		boostNode:   boostNode{boost: 1.0},
		rewriteNode: rewriteNode{rewrite: CONSTANT_SCORE},
		statesNode:  statesNode{maxDeterminizedStates: 10000},
		pattern:     pattern,
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
		return regexNodeUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return regexNodeUnionJoinTermsNode(n, o.(*TermsNode))
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *RegexpNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	case TERM_DSL_TYPE:
		return regexNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return regexNodeIntersectTermsNode(n, o.(*TermsNode))
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

func regexNodeUnionJoinTermNode(n *RegexpNode, t *TermNode) (AstNode, error) {
	if n.pattern.Match([]byte(t.value.(string))) {
		return n, nil
	} else {
		return lfNodeUnionJoinLfNode(n, t)
	}
}

func regexNodeUnionJoinTermsNode(n *RegexpNode, t *TermsNode) (AstNode, error) {
	var excludes = []LeafValue{}
	for _, term := range t.terms {
		if !n.pattern.Match([]byte(term.(string))) {
			excludes = append(excludes, term)
		}
	}
	return astNodeUnionJoinTermsNode(n, t, excludes)

}

func regexNodeIntersectTermNode(n *RegexpNode, o *TermNode) (AstNode, error) {
	if n.isArrayType() {
		return lfNodeIntersectLfNode(n, o)
	} else if n.pattern.Match([]byte(o.value.(string))) {
		return o, nil
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
	}
}

func regexNodeIntersectTermsNode(n *RegexpNode, o *TermsNode) (AstNode, error) {
	if n.isArrayType() {
		var excludes = []LeafValue{}
		for _, term := range o.terms {
			if !n.pattern.Match([]byte(term.(string))) {
				excludes = append(excludes, term)
			}
		}
		return astNodeIntersectTermsNode(n, o, excludes)
	} else {
		var includes = []LeafValue{}
		for _, term := range o.terms {
			if n.pattern.Match([]byte(term.(string))) {
				includes = append(includes, term)
			}
		}
		if len(includes) == 0 {
			return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
		} else if len(includes) == 1 {
			return &TermNode{
				kvNode: kvNode{
					fieldNode: n.fieldNode,
					valueNode: valueNode{
						valueType: n.valueType,
						value:     includes[0],
					},
				},
				boostNode: n.boostNode,
			}, nil
		} else {
			return &TermsNode{
				fieldNode: n.fieldNode,
				boostNode: n.boostNode,
				valueType: n.valueType,
				terms:     includes,
			}, nil
		}
	}
}

func (n *RegexpNode) ToDSL() DSL {
	return DSL{
		REGEXP_KEY: DSL{
			n.field: DSL{
				VALUE_KEY:   n.toPrintValue(),
				BOOST_KEY:   n.getBoost(),
				REWRITE_KEY: n.getRewrite(),

				MAX_DETERMINIZED_STATES_KEY: n.getMaxDeterminizedStates(),
			},
		},
	}
}
