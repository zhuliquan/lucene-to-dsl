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

type EqNode struct {
	Field string
	Value interface{}
}

type BndNode struct {
	Field string
	Bound *bnd.Bound
}

type TermNode struct {
	EqNode
	Boost float64
}

func (n *TermNode) GetDSLType() DSLType {
	return TERM_DSL_TYPE
}

func (n *TermNode) GetNodeType() NodeType {
	return LEAF_NODE_TYPE
}

func (n *TermNode) ToDSL() DSL {
	return DSL{"term": DSL{n.Field: DSL{"value": n.Value, "boost": n.Boost}}}
}

type RegexpNode struct {
	EqNode
}

func (n *RegexpNode) GetDSLType() DSLType {
	return REGEXP_DSL_TYPE
}

func (n *RegexpNode) GetNodeType() NodeType {
	return LEAF_NODE_TYPE
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

func (n *FuzzyNode) GetNodeType() DSLType {
	return NOT_DSL_TYPE
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
