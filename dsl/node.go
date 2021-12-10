package dsl

import (
	"fmt"
	"strconv"

	bnd "github.com/zhuliquan/lucene-to-dsl/internal/bound"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

// 定义dsl的ast node
type DSLNode interface {
	GetNodeType() NodeType
	GetDSLType() DSLType
	ToDSL() DSL
}

type OrDSLNode struct {
	Nodes map[string]DSLNode
}

type AndDSLNode struct {
	Nodes map[string]DSLNode
}

type NotDSLNode struct {
	Nodes map[string]DSLNode
}

type LeafNode struct{}

func (n *LeafNode) GetNodeType() NodeType { return LEAF_NODE_TYPE }

type ExistsNode struct {
	LeafNode
	Field string
}

func (n *ExistsNode) GetDSLType() DSLType {
	return EXISTS_DSL_TYPE
}

func (n *ExistsNode) ToDSL() DSL {
	if n == nil {
		return nil
	}
	return DSL{"exists": DSL{"field": n.Field}}
}

type IdsNode struct {
	LeafNode
	Values []string
}

func (n *IdsNode) GetDSLType() DSLType {
	return IDS_DSL_TYPE
}

func (n *IdsNode) ToDSL() DSL {
	if n == nil {
		return nil
	}
	return DSL{"ids": DSL{"values": n.Values}}
}

type EqNode struct {
	LeafNode
	Field string
	Value interface{}
}

type TermNode struct {
	LeafNode
	EqNode
	Boost float64
}

func (n *TermNode) GetDSLType() DSLType {
	return TERM_DSL_TYPE
}

func (n *TermNode) ToDSL() DSL {
	if n == nil {
		return nil
	}
	return DSL{"term": DSL{n.Field: DSL{"value": n.Value, "boost": n.Boost}}}
}

type TermsNode struct {
	LeafNode
	Field  string
	Values []string
	Boost  float64
}

func (n *TermsNode) GetDSLType() DSLType {
	return TERMS_DSL_TYPE
}

func (n *TermsNode) ToDSL() DSL {
	if n == nil {
		return nil
	}
	return DSL{"terms": DSL{n.Field: n.Values, "boost": n.Boost}}
}

type RegexpNode struct {
	EqNode
}

func (n *RegexpNode) GetDSLType() DSLType {
	return REGEXP_DSL_TYPE
}

func (n *RegexpNode) ToDSL() DSL {
	if n == nil {
		return nil
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

func (n *FuzzyNode) ToDSL() DSL {
	if n == nil {
		return nil
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

func (n *PrefixNode) ToDSL() DSL {
	if n == nil {
		return nil
	}
	return DSL{"prefix": DSL{n.Field: DSL{"values": n.Value}}}
}

type RangeNode struct {
	LeafNode
	Field   string
	Bound   *bnd.Bound
	Mapping *mapping.FieldMapping
}

func (n *RangeNode) GetDSLType() DSLType {
	return RANGE_DSL_TYPE
}

// func (n *RangeNode) ToDSL() DSL {
// 	if n == nil {
// 		return nil
// 	}
// 	return n.Bound.ToDSL(n.Field, n.Mapping)
// }

type WildCardNode struct {
	LeafNode
	EqNode
	Boost float64
}

func (n *WildCardNode) GetDSLType() DSLType {
	return WILDCARD_DSL_TYPE
}

func (n *WildCardNode) ToDSL() DSL {
	if n == nil {
		return nil
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

// func (n *MatchNode) ToDSL() DSL {
// 	return DSL{"match": }

// }