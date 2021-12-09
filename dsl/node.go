package dsl

import (
	"fmt"
	"strconv"

	bnd "github.com/zhuliquan/lucene-to-dsl/internal/bound"
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
	return DSL{"terms": DSL{n.Field: n.Values, "boost": n.Boost}}
}

type RegexpNode struct {
	EqNode
}

func (n *RegexpNode) GetDSLType() DSLType {
	return REGEXP_DSL_TYPE
}

func (n *RegexpNode) ToDSL() DSL {
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
	return DSL{"prefix": DSL{n.Field: DSL{"values": n.Value}}}
}

type RangeNode struct {
	LeafNode
	Field  string
	Format string
	Bound  *bnd.Bound
}

func (n *RangeNode) GetDSLType() DSLType {
	return RANGE_DSL_TYPE
}

// func (n *RangeNode) ToDSL() DSLType {
// 	return DSL{"range": DSL{n.Field: DSL{""}}}
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
	return DSL{"wildcard": DSL{n.Field: DSL{"values": n.Value, "boost": n.Boost}}}

}

type MatchNode struct {
	EqNode
	Boost float64
}

func (n *MatchNode) GetNodeType() NodeType {
	return LEAF_NODE_TYPE
}

func (n *MatchNode) GetDSLType() DSLType {
	return MATCH_DSL_TYPE
}

// func (n *MatchNode) ToDSL() DSL {
// 	return DSL{"match": }

// }
