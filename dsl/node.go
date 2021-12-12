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

type OpNode struct{}

func (n *OpNode) GetNodeType() NodeType {
	return OP_NODE_TYPE
}

type OrDSLNode struct {
	OpNode
	Nodes []DSLNode
}

func (n *OrDSLNode) GetDSLType() DSLType {
	return OR_DSL_TYPE
}

func (n *OrDSLNode) ToDSL() DSL {
	var res = []DSL{}
	for _, node := range n.Nodes {
		res = append(res, node.ToDSL())
	}
	if len(res) == 1 {
		return res[0]
	} else {
		return DSL{"bool": DSL{"should": res}}
	}
}

type AndDSLNode struct {
	OpNode
	MustNodes   []DSLNode
	FilterNodes []DSLNode
}

func (n *AndDSLNode) GetDSLType() DSLType {
	return OR_DSL_TYPE
}

func (n *AndDSLNode) ToDSL() DSL {
	var FRes = []DSL{}
	var MRes = []DSL{}
	for _, node := range n.MustNodes {
		MRes = append(MRes, node.ToDSL())
	}
	for _, node := range n.FilterNodes {
		FRes = append(FRes, node.ToDSL())
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
	Nodes []DSLNode
}

func (n *NotDSLNode) GetDSLType() DSLType {
	return NOT_DSL_TYPE
}

func (n *NotDSLNode) ToDSL() DSL {
	var res = []DSL{}
	for _, node := range n.Nodes {
		res = append(res, node.ToDSL())
	}
	if len(res) == 1 {
		return DSL{"bool": DSL{"must_not": res[0]}}
	} else {
		return DSL{"bool": DSL{"must_not": res[1]}}
	}
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
	Field  string
	Bound  *bnd.Bound
	Boost  float64
	Format string
}

func (n *RangeNode) GetDSLType() DSLType {
	return RANGE_DSL_TYPE
}

func (n *RangeNode) ToDSL() DSL {
	if n == nil || n.Bound == nil {
		return nil

		// (-inf, +inf) mean field exist
	} else if n.Bound.LeftValue.IsInf() && n.Bound.RightValue.IsInf() {
		return DSL{"exists": DSL{"field": n.Field}}
	}

	var (
		res  DSL
		infL = n.Bound.LeftValue.IsInf()
		infR = n.Bound.RightValue.IsInf()
	)

	switch n.Bound.GetBoundType() {
	case bnd.LEFT_INCLUDE_RIGHT_INCLUDE:
		res = DSL{"gte": n.Bound.LeftValue.Value(), "lte": n.Bound.RightValue.Value()}
	case bnd.LEFT_INCLUDE_RIGHT_EXCLUDE:
		if infR {
			res = DSL{"gte": n.Bound.LeftValue.Value()}
		} else {
			res = DSL{"gte": n.Bound.LeftValue.Value(), "lt": n.Bound.RightValue.Value()}
		}
	case bnd.LEFT_EXCLUDE_RIGHT_INCLUDE:
		if infL {
			res = DSL{"lte": n.Bound.RightValue.Value()}
		} else {
			res = DSL{"gt": n.Bound.LeftValue.Value(), "lte": n.Bound.RightValue.Value()}
		}
	case bnd.LEFT_EXCLUDE_RIGHT_EXCLUDE:
		if infL && !infR {
			res = DSL{"lt": n.Bound.RightValue.Value()}
		} else if !infL && infR {
			res = DSL{"gt": n.Bound.LeftValue.Value()}
		} else {
			res = DSL{"gt": n.Bound.LeftValue.Value(), "lt": n.Bound.RightValue.Value()}
		}
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

func (n *MatchNode) ToDSL() DSL {
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
