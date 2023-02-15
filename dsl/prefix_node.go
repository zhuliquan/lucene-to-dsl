package dsl

import (
	"fmt"
	"strings"

	"github.com/zhuliquan/lucene-to-dsl/utils"
)

type PrefixNode struct {
	kvNode
	rewriteNode
	patternNode
}

func NewPrefixNode(kvNode *kvNode, pattern utils.PatternMatcher, opts ...func(AstNode)) *PrefixNode {
	var n = &PrefixNode{
		kvNode:      *kvNode,
		rewriteNode: rewriteNode{rewrite: CONSTANT_SCORE},
		patternNode: patternNode{matcher: pattern},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *PrefixNode) DslType() DslType {
	return PREFIX_DSL_TYPE
}

func (n *PrefixNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		return patternNodeUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return patternNodeUnionJoinTermsNode(n, o.(*TermsNode))
	case PREFIX_DSL_TYPE:
		return prefixNodeUnionJoinPrefixNode(n, o.(*PrefixNode))
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *PrefixNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	case TERM_DSL_TYPE:
		return patternNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return patternNodeIntersectTermsNode(n, o.(*TermsNode))
	case PREFIX_DSL_TYPE:
		return prefixNodeIntersectPrefixNode(n, o.(*PrefixNode))
	default:
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *PrefixNode) Inverse() (AstNode, error) {
	return NewBoolNode(n, NOT), nil
}

func (n *PrefixNode) ToDSL() DSL {
	return DSL{
		PREFIX_KEY: DSL{
			n.field: DSL{
				VALUE_KEY:   n.toPrintValue(),
				REWRITE_KEY: n.getRewrite(),
			},
		},
	}
}

func prefixNodeUnionJoinPrefixNode(n, o *PrefixNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var prefixO = o.value.(string)
	if strings.HasPrefix(prefixN, prefixO) {
		return o, nil
	} else if strings.HasPrefix(prefixO, prefixN) {
		return n, nil
	} else {
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func prefixNodeIntersectPrefixNode(n, o *PrefixNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var prefixO = o.value.(string)
	if strings.HasPrefix(prefixN, prefixO) {
		return n, nil
	} else if strings.HasPrefix(prefixO, prefixN) {
		return o, nil
	} else if n.isArrayType() {
		return lfNodeIntersectLfNode(n, o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: prefix value is conflict", n.ToDSL(), o.ToDSL())
	}
}
