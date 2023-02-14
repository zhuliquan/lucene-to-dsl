package dsl

// fuzzy node represent fuzzy query
type FuzzyNode struct {
	kvNode
	rewriteNode
	expandsNode
	fuzziness      string
	prefixLength   int
	transpositions bool
}

func WithPrefixLength(prefixLength int) func(AstNode) {
	return func(n AstNode) {
		if v, ok := n.(*FuzzyNode); ok {
			v.prefixLength = prefixLength
		}
	}
}

func WithFuzziness(fuzziness string) func(AstNode) {
	return func(n AstNode) {
		if v, ok := n.(*FuzzyNode); ok {
			v.fuzziness = fuzziness
		}
	}
}

func WithTranspositions(transpositions bool) func(AstNode) {
	return func(n AstNode) {
		if v, ok := n.(*FuzzyNode); ok {
			v.transpositions = transpositions
		}
	}
}

func NewFuzzyNode(kvNode *kvNode, opts ...func(AstNode)) *FuzzyNode {
	var n = &FuzzyNode{
		kvNode:         *kvNode,
		rewriteNode:    rewriteNode{rewrite: CONSTANT_SCORE},
		expandsNode:    expandsNode{maxExpands: 50},
		fuzziness:      "AUTO",
		prefixLength:   0,
		transpositions: true,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *FuzzyNode) DslType() DslType {
	return FUZZY_DSL_TYPE
}

func (n *FuzzyNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *FuzzyNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	default:
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *FuzzyNode) Inverse() (AstNode, error) {
	return NewBoolNode(n, NOT), nil
}

func (n *FuzzyNode) ToDSL() DSL {
	return DSL{
		FUZZY_KEY: DSL{
			n.field: DSL{
				VALUE_KEY:          n.toPrintValue(),
				REWRITE_KEY:        n.rewrite,
				FUZZINESS_KEY:      n.fuzziness,
				PREFIX_LENGTH_KEY:  n.prefixLength,
				MAX_EXPANSIONS_KEY: n.getMaxExpands(),
				TRANSPOSITIONS_KEY: n.transpositions,
			},
		},
	}
}
