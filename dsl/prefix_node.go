package dsl

import "strings"

type PrefixNode struct {
	kvNode
	rewriteNode
}

func NewPrefixNode(kvNode *kvNode, opts ...func(AstNode)) *PrefixNode {
	var n = &PrefixNode{
		kvNode:      *kvNode,
		rewriteNode: rewriteNode{rewrite: CONSTANT_SCORE},
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
		return prefixUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return prefixUnionJoinTermsNode(n, o.(*TermsNode))
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
		return prefixNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return prefixNodeIntersectTermsNode(n, o.(*TermsNode))
	case PREFIX_DSL_TYPE:
		return prefixNodeIntersectPrefixNode(n, o.(*PrefixNode))
	default:
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *PrefixNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
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

func prefixUnionJoinTermNode(n *PrefixNode, o *TermNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var valueO = n.value.(string)
	if strings.HasPrefix(valueO, prefixN) {
		return n, nil
	} else {
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func prefixUnionJoinTermsNode(n *PrefixNode, o *TermsNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var excludes = []LeafValue{}
	for _, term := range o.terms {
		if !strings.HasPrefix(term.(string), prefixN) {
			excludes = append(excludes, term)
		}
	}
	return astNodeUnionJoinTermsNode(n, o, excludes)
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

func prefixNodeIntersectTermNode(n *PrefixNode, o *TermNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var valueO = n.value.(string)
	if strings.HasPrefix(valueO, prefixN) {
		return o, nil
	} else {
		return lfNodeIntersectLfNode(n, o)
	}
}

func prefixNodeIntersectTermsNode(n *PrefixNode, o *TermsNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var excludes = []LeafValue{}
	for _, term := range o.terms {
		if !strings.HasPrefix(term.(string), prefixN) {
			excludes = append(excludes, term)
		}
	}
	return astNodeIntersectTermsNode(n, o, excludes)
}

func prefixNodeIntersectPrefixNode(n, o *PrefixNode) (AstNode, error) {
	var prefixN = n.value.(string)
	var prefixO = o.value.(string)
	if strings.HasPrefix(prefixN, prefixO) {
		return n, nil
	} else if strings.HasPrefix(prefixO, prefixN) {
		return o, nil
	} else {
		return lfNodeIntersectLfNode(n, o)
	}
}
