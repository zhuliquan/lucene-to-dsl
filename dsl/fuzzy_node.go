package dsl

import (
	"fmt"
	"reflect"
	"strconv"
)

type FuzzyNode struct {
	KvNode
	LowFuzziness  int
	HighFuzziness int
}

func (n *FuzzyNode) DslType() DslType {
	return FUZZY_DSL_TYPE
}

func (n *FuzzyNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return fuzzyNodeUnionJoinAstNode[EXISTS_DSL_TYPE](o, n)
	case FUZZY_DSL_TYPE:
		return fuzzyNodeUnionJoinAstNode[FUZZY_DSL_TYPE](o, n)
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	}
}

var fuzzyNodeUnionJoinAstNode = map[DslType]func(AstNode, AstNode) (AstNode, error){
	EXISTS_DSL_TYPE: func(o, n AstNode) (AstNode, error) {
		return o.UnionJoin(n)
	},
	FUZZY_DSL_TYPE: func(n1, n2 AstNode) (AstNode, error) {
		var n = n2.(*FuzzyNode)
		var t = n1.(*FuzzyNode)
		if !reflect.DeepEqual(t.Value, n.Value) {
			return nil, fmt.Errorf("failed to union join %s to %s, err: value is conflict", t.ToDSL(), n.ToDSL())
		}

		var low, hig int
		if t.HighFuzziness == 0 && n.HighFuzziness == 0 {
			if t.LowFuzziness < n.LowFuzziness {
				return &FuzzyNode{KvNode: n.KvNode, LowFuzziness: t.LowFuzziness, HighFuzziness: n.LowFuzziness}, nil
			} else if t.LowFuzziness > n.LowFuzziness {
				return &FuzzyNode{KvNode: n.KvNode, LowFuzziness: n.LowFuzziness, HighFuzziness: t.LowFuzziness}, nil
			} else {
				return n, nil
			}
		} else if t.HighFuzziness != 0 && n.HighFuzziness == 0 {
			if t.LowFuzziness < n.LowFuzziness {
				low = t.LowFuzziness
				hig = n.LowFuzziness
			} else {
				low = n.LowFuzziness
				hig = t.LowFuzziness
			}

			if hig < t.HighFuzziness {
				hig = t.HighFuzziness
			}
		} else if t.HighFuzziness == 0 && n.HighFuzziness != 0 {
			if t.LowFuzziness < n.LowFuzziness {
				low = t.LowFuzziness
				hig = n.LowFuzziness
			} else {
				low = n.LowFuzziness
				hig = t.LowFuzziness
			}

			if hig < n.HighFuzziness {
				hig = n.HighFuzziness
			}
		} else {
			if t.LowFuzziness < n.LowFuzziness {
				low = t.LowFuzziness
			} else {
				low = n.LowFuzziness
			}

			if t.HighFuzziness < n.HighFuzziness {
				hig = n.HighFuzziness
			} else {
				hig = t.HighFuzziness
			}
		}

		if low != hig {
			return &FuzzyNode{KvNode: n.KvNode, LowFuzziness: low, HighFuzziness: hig}, nil
		} else {
			return &FuzzyNode{KvNode: n.KvNode, LowFuzziness: low}, nil
		}
	},
}

func (n *FuzzyNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case FUZZY_DSL_TYPE:
		var t = o.(*FuzzyNode)
		if !reflect.DeepEqual(t.Value, n.Value) || t.LowFuzziness != n.LowFuzziness || t.HighFuzziness != n.HighFuzziness {
			return nil, fmt.Errorf("failed to union join %s to %s, err: value is conflict", t.ToDSL(), n.ToDSL())
		} else {
			return n, nil
		}
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	}
}

func (n *FuzzyNode) Inverse() (AstNode, error) {
	return nil, fmt.Errorf("failed to inverse fuzzy node")
}

func (n *FuzzyNode) ToDSL() DSL {
	var fuzziness string
	if n.LowFuzziness != 0 && n.HighFuzziness != 0 {
		fuzziness = fmt.Sprintf("AUTO:%d,%d", n.LowFuzziness, n.HighFuzziness)
	} else if n.LowFuzziness != 0 {
		fuzziness = strconv.Itoa(n.LowFuzziness)
	} else if n.HighFuzziness != 0 {
		fuzziness = strconv.Itoa(n.HighFuzziness)
	} else {
		return EmptyDSL
	}
	return DSL{"fuzzy": DSL{n.Field: DSL{"value": n.Value, "fuzziness": fuzziness}}}
}
