package dsl

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// 定义dsl的ast node
type DSLNode interface {
	GetNodeType() NodeType
	GetDSLType() DSLType
	UnionJoin(DSLNode) (DSLNode, error)
	InterSect(DSLNode) (DSLNode, error)
	Inverse() (DSLNode, error)
	GetId() string
	ToDSL() DSL
}

type OpNode struct{}

func (n *OpNode) GetNodeType() NodeType {
	return OP_NODE_TYPE
}

type OrDSLNode struct {
	OpNode
	MinimumShouldMatch int
	Nodes              map[string][]DSLNode
}

func (n *OrDSLNode) GetDSLType() DSLType {
	return OR_DSL_TYPE
}

func (n *OrDSLNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if node == nil && n == nil {
		return nil, ErrUnionJoinNilNode
	} else if node == nil && n != nil {
		return n, nil
	} else if node != nil && n == nil {
		return node, nil
	}
	var t = node.(*OrDSLNode)
	for key, curNodes := range t.Nodes {
		if preNodes, ok := n.Nodes[key]; ok {
			if key == AND_OP_KEY || key == NOT_OP_KEY {
				n.Nodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].UnionJoin(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.Nodes, key)
					n.Nodes[key] = []DSLNode{newNode}
				}
			}

		} else {
			n.Nodes[key] = curNodes
		}
	}
	return n, nil
}

func (n *OrDSLNode) InterSect(DSLNode) (DSLNode, error) {
	return nil, nil
}

func (n *OrDSLNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return &NotDSLNode{
		Nodes: n.Nodes,
	}, nil
}

func (n *OrDSLNode) GetId() string { return OR_OP_KEY }

func (n *OrDSLNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	var res = []DSL{}
	for _, nodes := range n.Nodes {
		for _, node := range nodes {
			res = append(res, node.ToDSL())
		}
	}
	if len(res) == 1 {
		return res[0]
	} else {
		var shouldMatch = 1
		if n.MinimumShouldMatch != 0 {
			shouldMatch = n.MinimumShouldMatch
		}
		return DSL{"bool": DSL{"should": res}, "minimum_should_match": shouldMatch}
	}
}

type AndDSLNode struct {
	OpNode
	MustNodes   map[string][]DSLNode
	FilterNodes map[string][]DSLNode
}

func (n *AndDSLNode) GetDSLType() DSLType {
	return OR_DSL_TYPE
}

func (n *AndDSLNode) UnionJoin(DSLNode) (DSLNode, error) {
	return nil, nil
}

func (n *AndDSLNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	var t = node.(*AndDSLNode)
	for key, curNodes := range t.MustNodes {
		if preNodes, ok := n.MustNodes[key]; ok {
			if key == OR_OP_KEY {
				n.MustNodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.MustNodes, key)
					n.MustNodes[key] = []DSLNode{newNode}
				}
			}

		} else {
			n.MustNodes[key] = curNodes
		}
	}

	for key, curNodes := range t.FilterNodes {
		if preNodes, ok := n.FilterNodes[key]; ok {
			if key == OR_OP_KEY {
				n.FilterNodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.FilterNodes, key)
					n.FilterNodes[key] = []DSLNode{newNode}
				}
			}
		} else {
			n.FilterNodes[key] = curNodes
		}
	}

	return n, nil
}

func (n *AndDSLNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	var resNodes = make(map[string][]DSLNode)
	for key, nodes := range n.MustNodes {
		resNodes[key] = nodes
	}
	for key, nodes := range n.FilterNodes {
		resNodes[key] = nodes
	}
	return &OrDSLNode{Nodes: resNodes, MinimumShouldMatch: -1}, nil
}

func (n *AndDSLNode) GetId() string { return AND_OP_KEY }

func (n *AndDSLNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	var FRes = []DSL{}
	var MRes = []DSL{}
	for _, nodes := range n.MustNodes {
		for _, node := range nodes {
			MRes = append(MRes, node.ToDSL())
		}
	}
	for _, nodes := range n.FilterNodes {
		for _, node := range nodes {
			FRes = append(FRes, node.ToDSL())
		}
	}

	if len(FRes) == 1 && len(n.MustNodes) == 0 {
		return DSL{"bool": DSL{"filter": FRes[0]}}
	} else if len(FRes) == 1 && len(n.MustNodes) == 1 {
		return DSL{"bool": DSL{"must": MRes[0], "filter": FRes[0]}}
	} else if len(FRes) == 1 && len(n.MustNodes) > 1 {
		return DSL{"bool": DSL{"must": MRes, "filter": FRes[0]}}
	} else if len(FRes) == 0 && len(n.MustNodes) == 1 {
		return MRes[0]
	} else if len(FRes) == 0 && len(n.MustNodes) > 1 {
		return DSL{"bool": DSL{"must": MRes}}
	} else if len(FRes) > 1 && len(n.MustNodes) == 0 {
		return DSL{"bool": DSL{"filter": FRes}}
	} else if len(FRes) > 1 && len(n.MustNodes) == 1 {
		return DSL{"bool": DSL{"must": MRes[0], "filter": FRes}}
	} else {
		return DSL{"bool": DSL{"must": MRes, "filter": FRes}}
	}
}

type NotDSLNode struct {
	OpNode
	Nodes map[string][]DSLNode
}

func (n *NotDSLNode) GetDSLType() DSLType {
	return NOT_DSL_TYPE
}

func (n *NotDSLNode) GetId() string { return NOT_OP_KEY }

func (n *NotDSLNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	var res = []DSL{}
	for _, nodes := range n.Nodes {
		for _, node := range nodes {
			res = append(res, node.ToDSL())
		}
	}
	if len(res) == 1 {
		return DSL{"bool": DSL{"must_not": res[0]}}
	} else {
		return DSL{"bool": DSL{"must_not": res}}
	}
}

func (n *NotDSLNode) UnionJoin(node DSLNode) (DSLNode, error) {
	return nil, nil
}

func (n *NotDSLNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	var t = node.(*NotDSLNode)
	for key, curNodes := range t.Nodes {
		if preNodes, ok := n.Nodes[key]; ok {
			if key == OR_OP_KEY {
				n.Nodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.Nodes, key)
					n.Nodes[key] = []DSLNode{newNode}
				}
			}

		} else {
			n.Nodes[key] = curNodes
		}
	}
	return n, nil
}

// 全部都不是的反例是至少有一个
func (n *NotDSLNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return &OrDSLNode{Nodes: n.Nodes, MinimumShouldMatch: 1}, nil
}

type LeafNode struct{}

func (n *LeafNode) GetNodeType() NodeType { return LEAF_NODE_TYPE }

type ExistsNode struct {
	LeafNode
	Field string
}

func (n *ExistsNode) GetDSLType() DSLType { return EXISTS_DSL_TYPE }

// if union same field node, you can return exist node, for example {"exists": {"field" : "x"}} union {"match": {"x": "foo bar"}}
// "exists": {"field": "x"} > "match": {"x": "foo bar"}
func (n *ExistsNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	return n, nil
}

func (n *ExistsNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	return node, nil
}

func (n *ExistsNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return &NotDSLNode{
		Nodes: map[string][]DSLNode{
			n.Field: {n},
		},
	}, nil
}

func (n *ExistsNode) GetId() string { return "LEAF:" + n.Field }

func (n *ExistsNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"exists": DSL{"field": n.Field}}
}

type IdsNode struct {
	LeafNode
	Values []string
}

func (n *IdsNode) GetDSLType() DSLType { return IDS_DSL_TYPE }

func (n *IdsNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	if node.GetDSLType() == IDS_DSL_TYPE {
		var t = node.(*IdsNode)
		n.Values = UnionJoinStrLst(n.Values, t.Values)
		return n, nil
	} else {
		return nil, fmt.Errorf("failed to union join %v and %v, err: term type if conflict", n, node)
	}
}

func (n *IdsNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	if node.GetDSLType() == IDS_DSL_TYPE {
		var t = node.(*IdsNode)
		n.Values = IntersectStrLst(n.Values, t.Values)
		return n, nil
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: term type if conflict", n, node)
	}
}

func (n *IdsNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("ids node can't inverse own")
}

func (n *IdsNode) GetId() string { return "LEAF:" + "_id" }

func (n *IdsNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"ids": DSL{"values": n.Values}}
}

type EqNode struct {
	LeafNode
	Field string
	Type  DSLTermType
	Value *DSLTermValue
}

func (n *EqNode) GetId() string { return "LEAF:" + n.Field }

type TermNode struct {
	EqNode
	Boost float64
}

func (n *TermNode) GetDSLType() DSLType {
	return TERM_DSL_TYPE
}

func (n *TermNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case TERM_DSL_TYPE:
		var t = node.(*TermNode)
		if math.Abs(n.Boost-t.Boost) <= 1E-6 {
			if CompareAny(n.Value, t.Value, n.Type) == 0 {
				return n, nil
			} else {
				return &TermsNode{
					Field:  n.Field,
					Type:   n.Type,
					Values: []*DSLTermValue{n.Value, t.Value},
					Boost:  n.Boost,
				}, nil
			}
		} else {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), t.ToDSL())
		}
	case TERMS_DSL_TYPE:
		var t = node.(*TermsNode)
		if math.Abs(n.Boost-t.Boost) <= 1E-6 {
			t.Values = append(t.Values, n.Value)
			return t, nil
		} else {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), t.ToDSL())
		}
	case RANGE_DSL_TYPE:
		// put logic of compare and collision into range node
		return node.(*RangeNode).UnionJoin(n)
	case QUERY_STRING_DSL_TYPE:
		var t = node.(*QueryStringNode)
		if math.Abs(n.Boost-t.Boost) <= 1E-6 {
			// t.Value = fmt.Sprintf("%s OR %v", t.Value.StringTerm, n.Value)
			return t, nil
		} else {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), t.ToDSL())
		}
	case IDS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is conflict", n.ToDSL(), node.ToDSL())

	default:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is unknown", n.ToDSL(), node.ToDSL())
	}
}

func (n *TermNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case TERM_DSL_TYPE:
		var t = node.(*TermNode)
		if math.Abs(n.Boost-t.Boost) <= 1E-6 {
			if reflect.DeepEqual(n.Value, t.Value) {
				return node, nil
			} else {
				return nil, nil
				// return &QueryStringNode{EqNode: EqNode{Field: n.Field, Value: fmt.Sprintf("%v AND %v", n.Value, t.Value)}, Boost: n.Boost}, nil
			}
		} else {
			return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), t.ToDSL())
		}
	case RANGE_DSL_TYPE:
		// put logic of compare and collision into range node
		return node.(*RangeNode).UnionJoin(n)
	case IDS_DSL_TYPE, TERMS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), node.ToDSL())
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is unknown", n.ToDSL(), node.ToDSL())
	}
}

func (n *TermNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return &NotDSLNode{
		Nodes: map[string][]DSLNode{
			n.Field: {n},
		},
	}, nil
}

func (n *TermNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"term": DSL{n.Field: DSL{"value": n.Value, "boost": n.Boost}}}
}

type TermsNode struct {
	LeafNode
	Field  string
	Type   DSLTermType
	Values []*DSLTermValue
	Boost  float64
}

func (n *TermsNode) GetDSLType() DSLType {
	return TERMS_DSL_TYPE
}

func (n *TermsNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case TERM_DSL_TYPE:
		return node.UnionJoin(n)
	case TERMS_DSL_TYPE:
		var t = node.(*TermsNode)
		if math.Abs(n.Boost-t.Boost) <= 1E-6 {
			t.Values = append(t.Values, n.Values...)
			return t, nil
		} else {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), t.ToDSL())
		}
	case RANGE_DSL_TYPE:
		// put logic of compare and collision into range node
		return node.(*RangeNode).UnionJoin(n)

	case QUERY_STRING_DSL_TYPE:
		var t = node.(*QueryStringNode)
		if math.Abs(n.Boost-t.Boost) <= 1E-6 {
			var s = ""
			for _, val := range n.Values {
				s += fmt.Sprintf(" OR %s", val)
			}
			// t.Value = fmt.Sprintf("%s%s", t.Value, s)
			return t, nil
		} else {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), t.ToDSL())
		}
	case IDS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is conflict", n.ToDSL(), node.ToDSL())

	default:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is unknown", n.ToDSL(), node.ToDSL())

	}
}

func (n *TermsNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case RANGE_DSL_TYPE:
		// put logic of compare and collision into range node
		return node.(*RangeNode).UnionJoin(n)

	case TERM_DSL_TYPE:
		return nil, nil
	case TERMS_DSL_TYPE:
		return nil, nil
	case IDS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), node.ToDSL())
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is unknown", n.ToDSL(), node.ToDSL())
	}
}

func (n *TermsNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	var nodes = []DSLNode{}
	for _, val := range n.Values {
		nodes = append(nodes, &TermNode{EqNode: EqNode{Field: n.Field, Value: val}, Boost: n.Boost})
	}
	return &NotDSLNode{Nodes: map[string][]DSLNode{n.Field: nodes}}, nil
}

func (n *TermsNode) GetId() string { return "LEAF:" + n.Field }

func (n *TermsNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"terms": DSL{n.Field: n.Values, "boost": n.Boost}}
}

type RegexpNode struct {
	EqNode
}

func (n *RegexpNode) GetDSLType() DSLType {
	return REGEXP_DSL_TYPE
}

func (n *RegexpNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to union join regexp node")
	}
}

func (n *RegexpNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to intersect regexp node")
	}
}

func (n *RegexpNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("failed to inverse regexp node")
}

func (n *RegexpNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"regexp": DSL{n.Field: DSL{"value": n.Value}}}
}

type FuzzyNode struct {
	EqNode
	LowFuzziness  int
	HighFuzziness int
}

func (n *FuzzyNode) GetDSLType() DSLType {
	return FUZZY_DSL_TYPE
}

func (n *FuzzyNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case FUZZY_DSL_TYPE:
		var t = node.(*FuzzyNode)
		if !reflect.DeepEqual(t.Value, n.Value) {
			return nil, fmt.Errorf("failed to union join %s to %s, err: value is conflict", t.ToDSL(), n.ToDSL())
		}

		var low, high int
		if t.HighFuzziness == 0 && n.HighFuzziness == 0 {
			if t.LowFuzziness < n.LowFuzziness {
				return &FuzzyNode{EqNode: n.EqNode, LowFuzziness: t.LowFuzziness, HighFuzziness: n.LowFuzziness}, nil
			} else if t.LowFuzziness > n.LowFuzziness {
				return &FuzzyNode{EqNode: n.EqNode, LowFuzziness: n.LowFuzziness, HighFuzziness: t.LowFuzziness}, nil
			} else {
				return n, nil
			}
		} else if t.HighFuzziness != 0 && n.HighFuzziness == 0 {
			if t.LowFuzziness < n.LowFuzziness {
				low = t.LowFuzziness
				high = n.LowFuzziness
			} else {
				low = n.LowFuzziness
				high = t.LowFuzziness
			}

			if high < t.HighFuzziness {
				high = t.HighFuzziness
			}
		} else if t.HighFuzziness == 0 && n.HighFuzziness != 0 {
			if t.LowFuzziness < n.LowFuzziness {
				low = t.LowFuzziness
				high = n.LowFuzziness
			} else {
				low = n.LowFuzziness
				high = t.LowFuzziness
			}

			if high < n.HighFuzziness {
				high = n.HighFuzziness
			}
		} else {
			if t.LowFuzziness < n.LowFuzziness {
				low = t.LowFuzziness
			} else {
				low = n.LowFuzziness
			}

			if t.HighFuzziness < n.HighFuzziness {
				high = n.HighFuzziness
			} else {
				high = t.HighFuzziness
			}
		}

		if low != high {
			return &FuzzyNode{EqNode: n.EqNode, LowFuzziness: low, HighFuzziness: high}, nil
		} else {
			return &FuzzyNode{EqNode: n.EqNode, LowFuzziness: low}, nil
		}
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), node.ToDSL())
	}
}

func (n *FuzzyNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case FUZZY_DSL_TYPE:
		var t = node.(*FuzzyNode)
		if !reflect.DeepEqual(t.Value, n.Value) || t.LowFuzziness != n.LowFuzziness || t.HighFuzziness != n.HighFuzziness {
			return nil, fmt.Errorf("failed to union join %s to %s, err: value is conflict", t.ToDSL(), n.ToDSL())
		} else {
			return n, nil
		}
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), node.ToDSL())
	}
}

func (n *FuzzyNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("failed to inverse prefix node")
}

func (n *FuzzyNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	var fuzziness string
	if n.LowFuzziness != 0 && n.HighFuzziness != 0 {
		fuzziness = fmt.Sprintf("AUTO:%d,%d", n.LowFuzziness, n.HighFuzziness)
	} else if n.LowFuzziness != 0 {
		fuzziness = strconv.Itoa(n.LowFuzziness)
	} else if n.HighFuzziness != 0 {
		fuzziness = strconv.Itoa(n.HighFuzziness)
	} else {
		return DSL{}
	}
	return DSL{"fuzzy": DSL{n.Field: DSL{"value": n.Value, "fuzziness": fuzziness}}}
}

type PrefixNode struct {
	EqNode
}

func (n *PrefixNode) GetDSLType() DSLType {
	return PREFIX_DSL_TYPE
}

func (n *PrefixNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to union join prefix node")
	}
}

func (n *PrefixNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to intersect prefix node")
	}
}

func (n *PrefixNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("failed to inverse prefix node")
}

func (n *PrefixNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"prefix": DSL{n.Field: DSL{"values": n.Value}}}
}

type RangeNode struct {
	LeafNode
	Field       string
	ValueType   DSLTermType
	LeftValue   *DSLTermValue
	RightValue  *DSLTermValue
	LeftCmpSym  CompareType
	RightCmpSym CompareType
	Boost       float64
	Format      string
}

func (n *RangeNode) GetDSLType() DSLType {
	return RANGE_DSL_TYPE
}

func (n *RangeNode) GetId() string { return "LEAF:" + n.Field }

// TODO: very hard
func (n *RangeNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	case TERMS_DSL_TYPE:
		return nil, nil
	case TERM_DSL_TYPE:
		return nil, nil

	case RANGE_DSL_TYPE:
		// var otherNode = node.(*RangeNode)
		// if otherNode.LeftCmpSym == LT

	}
	return nil, nil
}

// TODO: very hard
func (n *RangeNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}

	return nil, nil
}

func (n *RangeNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, nil
}

func (n *RangeNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
		// (-inf, +inf) mean field exist
	} else if n.LeftValue == InfValue && n.RightValue == InfValue {
		return DSL{"exists": DSL{"field": n.Field}}
	}
	var res = DSL{}
	if n.LeftValue != InfValue {
		res[n.LeftCmpSym.String()] = n.LeftValue
	}
	if n.RightValue != InfValue {
		res[n.RightCmpSym.String()] = n.RightValue
	}

	res["relation"] = "WITHIN"
	res["boost"] = n.Boost
	if len(n.Format) != 0 {
		res["format"] = n.Format
	}
	return DSL{"range": DSL{n.Field: res}}
}

type WildCardNode struct {
	LeafNode
	EqNode
	Boost float64
}

func (n *WildCardNode) GetDSLType() DSLType {
	return WILDCARD_DSL_TYPE
}

func (n *WildCardNode) UnionJoin(node DSLNode) (DSLNode, error) {
	if n == nil && node == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && node != nil {
		return node, nil
	} else if n != nil && node == nil {
		return n, nil
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to union join wildcard node")
	}
}

func (n *WildCardNode) InterSect(node DSLNode) (DSLNode, error) {
	if n == nil || node == nil {
		return nil, ErrIntersectNilNode
	}
	switch node.GetDSLType() {
	case EXISTS_DSL_TYPE:
		return node.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to intersect wildcard node")
	}
}

func (n *WildCardNode) Inverse() (DSLNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("failed to inverse wildcard node")
}

func (n *WildCardNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"wildcard": DSL{n.Field: DSL{"values": n.Value, "boost": n.Boost}}}

}

type MatchNode struct {
	EqNode
	Boost float64
}

func (n *MatchNode) GetDSLType() DSLType {
	return MATCH_DSL_TYPE
}

func (n *MatchNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"match": DSL{n.Field: DSL{"value": n.Value, "boost": n.Boost}}}
}

type MatchPhraseNode struct {
	EqNode
}

func (n *MatchPhraseNode) GetDSLType() DSLType {
	return MATCH_PHRASE_DSL_TYPE
}

func (n *MatchPhraseNode) ToDSL() DSL {
	return DSL{"match_phrase": DSL{n.Field: n.Value}}
}

type QueryStringNode struct {
	EqNode
	Boost float64
}

func (n *QueryStringNode) UnionJoin(node DSLNode) (DSLNode, error) {
	return nil, nil
}

func (n *QueryStringNode) InterSect(node DSLNode) (DSLNode, error) {
	return nil, nil
}

func (n *QueryStringNode) Inverse() (DSLNode, error) {
	return nil, nil
}

func (n *QueryStringNode) GetDSLType() DSLType {
	return QUERY_STRING_DSL_TYPE
}

func (n *QueryStringNode) ToDSL() DSL {
	return DSL{"query_string": DSL{"query": n.Value, "default_field": n.Field, "boost": n.Boost}}
}
