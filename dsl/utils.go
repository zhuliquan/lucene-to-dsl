package dsl

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/x448/float16"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	"github.com/zhuliquan/scaled_float"
)

func FindAny(vl []LeafValue, v LeafValue, typ mapping.FieldType) int {
	idx := -1
	for i, x := range vl {
		if CompareAny(x, v, typ) == 0 {
			idx = i
			break
		}
	}
	return idx
}

func BinaryFindAny(vl []LeafValue, v LeafValue, typ mapping.FieldType) int {
	l, r := -1, len(vl)
	for r != l+1 {
		m := (l + r) >> 1
		f := CompareAny(vl[m], v, typ)
		if f < 0 {
			l = m
		} else {
			r = m
		}
	}
	if r >= len(vl) || CompareAny(vl[r], v, typ) != 0 {
		return -1
	}
	return r
}

func ValueLstToStrLst(vl []LeafValue) []string {
	var rl = make([]string, len(vl))
	for i, v := range vl {
		rl[i] = v.(string)
	}
	return rl
}

func StrLstToValueLst(vl []string) []LeafValue {
	var rl = make([]LeafValue, len(vl))
	for i, v := range vl {
		rl[i] = v
	}
	return rl
}

// union join two leaf value slice
func UnionJoinValueLst(al, bl []LeafValue, typ mapping.FieldType) []LeafValue {
	sort.Slice(al, func(i, j int) bool { return CompareAny(al[i], al[j], typ) < 0 })
	sort.Slice(bl, func(i, j int) bool { return CompareAny(bl[i], bl[j], typ) < 0 })
	var cl = make([]LeafValue, 0, len(al)+len(bl))
	var i, j, na, nb = 0, 0, len(al), len(bl)

	for i < na || j < nb {
		if i == na || (j < nb && CompareAny(al[i], bl[j], typ) > 0) {
			cl = append(cl, bl[j])
			j += 1
		} else {
			cl = append(cl, al[i])
			i += 1
		}
	}
	return UniqValueLst(cl, typ)
}

// intersect two leaf value slice
func IntersectValueLst(al, bl []LeafValue, typ mapping.FieldType) []LeafValue {
	sort.Slice(al, func(i, j int) bool { return CompareAny(al[i], al[j], typ) < 0 })
	sort.Slice(bl, func(i, j int) bool { return CompareAny(bl[i], bl[j], typ) < 0 })
	var cl = make([]LeafValue, 0, len(al)+len(bl))
	var i, j, na, nb = 0, 0, len(al), len(bl)

	for i < na && j < nb {
		if al[i] == bl[j] {
			cl = append(cl, al[i])
			i += 1
			j += 1
		} else if CompareAny(al[i], bl[j], typ) > 0 {
			j += 1
		} else {
			i += 1
		}
	}
	return UniqValueLst(cl, typ)
}

// difference two leaf value slice
func DifferenceValueList(al, bl []LeafValue, typ mapping.FieldType) []LeafValue {
	sort.Slice(al, func(i, j int) bool { return CompareAny(al[i], al[j], typ) < 0 })
	sort.Slice(bl, func(i, j int) bool { return CompareAny(bl[i], bl[j], typ) < 0 })
	var cl = make([]LeafValue, 0, len(al)+len(bl))
	for i, na := 0, len(al); i < na; i++ {
		if BinaryFindAny(bl, al[i], typ) == -1 {
			cl = append(cl, al[i])
		}
	}
	return UniqValueLst(cl, typ)
}

// uniq a sort string slice
func UniqValueLst(a []LeafValue, typ mapping.FieldType) []LeafValue {
	if len(a) == 0 || len(a) == 1 {
		return a
	}
	var r = []LeafValue{}
	for i, n := 0, len(a); i < n; i++ {
		if i == 0 {
			r = append(r, a[i])
		} else if CompareAny(a[i], a[i-1], typ) != 0 {
			r = append(r, a[i])
		}
	}
	return r
}

// -1  mean a < b
// +1  mean a > b
// 0   mean a == b
func CompareAny(a, b LeafValue, typ mapping.FieldType) int {
	var compare = compareFunc[typ]
	if compare != nil {
		return compare(a, b)
	} else {
		return 0
	}
}

var compareFunc = map[mapping.FieldType]func(LeafValue, LeafValue) int{
	mapping.IP_FIELD_TYPE:               compareIp,
	mapping.IP_RANGE_FIELD_TYPE:         compareIp,
	mapping.DATE_FIELD_TYPE:             compareDate,
	mapping.DATE_RANGE_FIELD_TYPE:       compareDate,
	mapping.DATE_NANOS_FIELD_TYPE:       compareDate,
	mapping.BYTE_FIELD_TYPE:             compareInt,
	mapping.SHORT_FIELD_TYPE:            compareInt,
	mapping.INTEGER_FIELD_TYPE:          compareInt,
	mapping.INTEGER_RANGE_FIELD_TYPE:    compareInt,
	mapping.LONG_FIELD_TYPE:             compareInt,
	mapping.LONG_RANGE_FIELD_TYPE:       compareInt,
	mapping.UNSIGNED_LONG_FIELD_TYPE:    compareUint,
	mapping.FLOAT_FIELD_TYPE:            compareFloat,
	mapping.FLOAT_RANGE_FIELD_TYPE:      compareFloat,
	mapping.DOUBLE_FIELD_TYPE:           compareFloat,
	mapping.DOUBLE_RANGE_FIELD_TYPE:     compareFloat,
	mapping.HALF_FLOAT_FIELD_TYPE:       compareFloat16,
	mapping.SCALED_FLOAT_FIELD_TYPE:     compareDecimal,
	mapping.KEYWORD_FIELD_TYPE:          compareString,
	mapping.TEXT_FIELD_TYPE:             compareString,
	mapping.WILDCARD_FIELD_TYPE:         compareString,
	mapping.CONSTANT_KEYWORD_FIELD_TYPE: compareString,
	mapping.VERSION_FIELD_TYPE:          compareVersion,
}

func compareIp(a, b LeafValue) int {
	var ai = []byte(a.(net.IP))
	var bi = []byte(b.(net.IP))
	return bytes.Compare(ai, bi)
}

func castUInt(x LeafValue) uint64 {
	switch x.(interface{}).(type) {
	case int:
		return uint64(x.(int))
	case uint:
		return uint64(x.(uint))
	default:
		return x.(uint64)
	}
}

func castInt(x LeafValue) int64 {
	switch x.(interface{}).(type) {
	case int:
		return int64(x.(int))
	case uint:
		return int64(x.(uint))
	default:
		return x.(int64)
	}
}

func compareInt(a, b LeafValue) int {
	var ai = castInt(a)
	var bi = castInt(b)
	if ai < bi {
		return -1
	} else if ai > bi {
		return 1
	} else {
		return 0
	}
}

func compareDate(a, b LeafValue) int {
	var at = a.(time.Time)
	var bt = b.(time.Time)
	if at.UnixNano() == bt.UnixNano() {
		return 0
	} else if at.Before(bt) {
		return -1
	} else {
		return 1
	}
}

func compareUint(a, b LeafValue) int {
	var au = castUInt(a)
	var bu = castUInt(b)
	if au < bu {
		return -1
	} else if au > bu {
		return 1
	} else {
		return 0
	}
}

func compareFloat(a, b LeafValue) int {
	var af = a.(float64)
	var bf = b.(float64)
	if math.Abs(af-bf) <= eps {
		return 0
	} else if af < bf {
		return -1
	} else {
		return 1
	}
}

func compareFloat16(a, b LeafValue) int {
	var af = a.(float16.Float16).Float32()
	var bf = b.(float16.Float16).Float32()
	if math.Abs(float64(af-bf)) <= eps {
		return 0
	} else if af < bf {
		return -1
	} else {
		return 1
	}
}

func compareDecimal(a, b LeafValue) int {
	var ad = a.(*scaled_float.ScaledFloat)
	var bd = b.(*scaled_float.ScaledFloat)
	return ad.Compare(bd)
}

func compareString(a, b LeafValue) int {
	var as = a.(string)
	var bs = b.(string)
	if as > bs {
		return 1
	} else if as < bs {
		return -1
	} else {
		return 0
	}
}

func compareVersion(a, b LeafValue) int {
	var av = a.(*version.Version)
	var bv = b.(*version.Version)
	return av.Compare(bv)
}

func CheckValidRangeNode(node *RangeNode) error {
	var cmp = CompareAny(node.lValue, node.rValue, node.mType)
	if cmp > 0 || (cmp == 0 && (node.lCmpSym == GT || node.rCmpSym == LT)) {
		return fmt.Errorf("range is conflict")
	}
	return nil
}

func isMinInf(a LeafValue, t mapping.FieldType) bool {
	return CompareAny(a, minInf[t], t) == 0
}

func isMaxInf(a LeafValue, t mapping.FieldType) bool {
	return CompareAny(a, maxInf[t], t) == 0
}

var minInf = map[mapping.FieldType]LeafValue{
	mapping.BYTE_FIELD_TYPE:             MinInt[8],
	mapping.SHORT_FIELD_TYPE:            MinInt[16],
	mapping.INTEGER_FIELD_TYPE:          MinInt[32],
	mapping.INTEGER_RANGE_FIELD_TYPE:    MinInt[32],
	mapping.LONG_FIELD_TYPE:             MinInt[64],
	mapping.LONG_RANGE_FIELD_TYPE:       MinInt[64],
	mapping.UNSIGNED_LONG_FIELD_TYPE:    MinUint,
	mapping.HALF_FLOAT_FIELD_TYPE:       MinFloat[16],
	mapping.FLOAT_FIELD_TYPE:            MinFloat[32],
	mapping.FLOAT_RANGE_FIELD_TYPE:      MinFloat[32],
	mapping.DOUBLE_FIELD_TYPE:           MinFloat[64],
	mapping.DOUBLE_RANGE_FIELD_TYPE:     MinFloat[64],
	mapping.SCALED_FLOAT_FIELD_TYPE:     MinFloat[128],
	mapping.IP_FIELD_TYPE:               MinIP,
	mapping.IP_RANGE_FIELD_TYPE:         MinIP,
	mapping.DATE_FIELD_TYPE:             MinTime,
	mapping.DATE_NANOS_FIELD_TYPE:       MinTime,
	mapping.DATE_RANGE_FIELD_TYPE:       MinTime,
	mapping.VERSION_FIELD_TYPE:          MinVersion,
	mapping.KEYWORD_FIELD_TYPE:          MinString,
	mapping.TEXT_FIELD_TYPE:             MinString,
	mapping.WILDCARD_FIELD_TYPE:         MinString,
	mapping.CONSTANT_KEYWORD_FIELD_TYPE: MinString,
}

var maxInf = map[mapping.FieldType]LeafValue{
	mapping.BYTE_FIELD_TYPE:             MaxInt[8],
	mapping.SHORT_FIELD_TYPE:            MaxInt[16],
	mapping.INTEGER_FIELD_TYPE:          MaxInt[32],
	mapping.INTEGER_RANGE_FIELD_TYPE:    MaxInt[32],
	mapping.LONG_FIELD_TYPE:             MaxInt[64],
	mapping.LONG_RANGE_FIELD_TYPE:       MaxInt[64],
	mapping.UNSIGNED_LONG_FIELD_TYPE:    MaxUint[64],
	mapping.HALF_FLOAT_FIELD_TYPE:       MaxFloat[16],
	mapping.FLOAT_FIELD_TYPE:            MaxFloat[32],
	mapping.FLOAT_RANGE_FIELD_TYPE:      MaxFloat[32],
	mapping.DOUBLE_FIELD_TYPE:           MaxFloat[64],
	mapping.DOUBLE_RANGE_FIELD_TYPE:     MaxFloat[64],
	mapping.SCALED_FLOAT_FIELD_TYPE:     MaxFloat[128],
	mapping.IP_FIELD_TYPE:               MaxIP,
	mapping.IP_RANGE_FIELD_TYPE:         MaxIP,
	mapping.DATE_FIELD_TYPE:             MaxTime,
	mapping.DATE_NANOS_FIELD_TYPE:       MaxTime,
	mapping.DATE_RANGE_FIELD_TYPE:       MaxTime,
	mapping.VERSION_FIELD_TYPE:          MaxVersion,
	mapping.KEYWORD_FIELD_TYPE:          MaxString,
	mapping.TEXT_FIELD_TYPE:             MaxString,
	mapping.WILDCARD_FIELD_TYPE:         MaxString,
	mapping.CONSTANT_KEYWORD_FIELD_TYPE: MaxString,
}

// union join two leaf node
func lfNodeUnionJoinLfNode(key string, a, b AstNode) (AstNode, error) {
	orNode := newDefaultBoolNode(OR)
	orNode.Should = map[string][]AstNode{key: {a, b}}
	return orNode, nil
}

// intersect two leaf node
func lfNodeIntersectLfNode(key string, a, b AstNode) (AstNode, error) {
	af := a.(FilterCtxNode)
	bf := b.(FilterCtxNode)
	var mustNodes = map[string][]AstNode{key: {}}
	var filterNodes = map[string][]AstNode{key: {}}
	if af.getFilterCtx() {
		filterNodes[key] = append(filterNodes[key], a)
	} else {
		mustNodes[key] = append(mustNodes[key], a)
	}
	if bf.getFilterCtx() {
		filterNodes[key] = append(filterNodes[key], b)
	} else {
		mustNodes[key] = append(mustNodes[key], b)
	}

	var andNode = newDefaultBoolNode(AND)
	if len(mustNodes[key]) != 0 {
		andNode.Must = mustNodes
	}
	if len(filterNodes[key]) != 0 {
		andNode.Filter = filterNodes
	}
	return andNode, nil
}

func flattenNodes(nodesMap map[string][]AstNode) []AstNode {
	var nodes = []AstNode{}
	for _, nodeList := range nodesMap {
		nodes = append(nodes, nodeList...)
	}
	return nodes
}

func ReduceAstNode(x AstNode) AstNode {
	if x.AstType() == OP_NODE_TYPE {
		n := x.(*BoolNode)
		switch n.opType {
		case AND:
			if len(n.Must) == 1 && len(n.Filter) == 0 {
				nodes := flattenNodes(n.Must)
				if len(nodes) == 1 {
					return nodes[0]
				}
			}
			return n
		case OR:
			if len(n.Should) == 1 {
				nodes := flattenNodes(n.Should)
				if len(nodes) == 1 {
					return nodes[0]
				}
			}
			return n
		default:
			return n
		}
	} else {
		return x
	}
}

func reduceAstNodes(nodes []AstNode, mergeMethodName string, mergeMethodFunc MergeMethodFunc) ([]AstNode, error) {
	for before, first := nodes, true; ; {
		var join bool
		var node AstNode
		var rest []AstNode
		var lo, up = 0, len(before) - 1
		if first {
			lo = len(before) - 1
		}
		for k := up; k >= lo; k-- {
			node, rest = restAstNodes(before, k)
			for i, n1 := range rest {
				if n2, err := mergeMethodFunc(n1, node); err == nil {
					if n2.AstType() != OP_NODE_TYPE { // merge two nodes into a single node as soon as possible
						rest[i] = n2
						join = true
						goto check
					}
				} else {
					return nil, fmt.Errorf("failed to %s node: %v, err: %+v", mergeMethodName, node, err)
				}
			}
		}
	check:
		if !join {
			return before, nil
		} else {
			first = false
			before = rest // loop find any other node which can be merge with n0
		}
	}
}

// split node[i] from nodes
func restAstNodes(nodes []AstNode, index int) (AstNode, []AstNode) {
	if len(nodes) == 0 {
		return nil, nil
	}

	var n = len(nodes)
	if index < 0 {
		return nodes[0], nodes[1:]
	} else if index < n {
		var res1 = nodes[index]
		// NOTE: not use `append(nodes[:index],nodes[index+1:]...)`, because it can modify nodes
		var res2 = make([]AstNode, 0, n-1)
		res2 = append(res2, nodes[:index]...)
		res2 = append(res2, nodes[index+1:]...)
		return res1, res2
	} else {
		return nodes[n-1], nodes[:n-1]
	}
}

func nodesToDSLList(nodes []AstNode) []DSL {
	var dslList = []DSL{}
	for _, node := range nodes {
		dslList = append(dslList, node.ToDSL())
	}
	return dslList
}

func reduceDSLList(dslList []DSL) interface{} {
	if len(dslList) == 0 {
		return nil
	} else if len(dslList) == 1 {
		return dslList[0]
	} else {
		return dslList
	}
}

func inverseNode(node AstNode) AstNode {
	var boolNode = newDefaultBoolNode(NOT)
	boolNode.MustNot = map[string][]AstNode{
		node.NodeKey(): {node},
	}
	return boolNode
}

func minEditDistance(termWord1, termWord2 string) int {
	var (
		l1 int16 = int16(len(termWord1))
		l2 int16 = int16(len(termWord2))

		i int16 = 0
		j int16 = 0
	)
	var dp = make([][]int16, l1+1)
	for i = 0; i <= l1; i++ {
		dp[i] = make([]int16, l2+1)
	}

	for i = 0; i <= l1; i++ {
		dp[i][0] = i
	}
	for j = 0; j <= l2; j++ {
		dp[0][j] = j
	}

	for i = 1; i <= l1; i++ {
		for j = 1; j <= l2; j++ {
			if termWord1[i-1] == termWord2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = minInt16(dp[i-1][j-1], minInt16(dp[i-1][j], dp[i][j-1])) + 1
			}
		}
	}
	return int(dp[l1][l2])
}

func minInt16(a, b int16) int16 {
	if a < b {
		return a
	} else {
		return b
	}
}

func leafValueToPrintValue(x LeafValue, t mapping.FieldType) interface{} {
	if mapping.CheckDateType(t) {
		return x.(time.Time).UnixNano() / 1e6
	} else if mapping.CheckIPType(t) {
		return x.(net.IP).String()
	} else if mapping.CheckVersionType(t) {
		return x.(*version.Version).String()
	} else if t == mapping.HALF_FLOAT_FIELD_TYPE {
		return x.(float16.Float16).Float32()
	} else if t == mapping.SCALED_FLOAT_FIELD_TYPE {
		return x.(*scaled_float.ScaledFloat).RawFloat()
	} else {
		return x
	}
}

func compareBoost(a, b BoostNode) int {
	return CompareAny(a.getBoost(), b.getBoost(), mapping.DOUBLE_FIELD_TYPE)
}

func patternNodeUnionJoinTermNode(n PatternNode, o *TermNode) (AstNode, error) {
	if n.Match([]byte(o.value.(string))) {
		return n.(AstNode), nil
	} else {
		return lfNodeUnionJoinLfNode(o.NodeKey(), n.(AstNode), o)
	}
}

func patternNodeIntersectTermNode(n PatternNode, o *TermNode) (AstNode, error) {
	if n.Match([]byte(o.value.(string))) {
		return o, nil
	} else if n.(ArrayTypeNode).isArrayType() {
		return lfNodeIntersectLfNode(o.NodeKey(), n.(AstNode), o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.(AstNode).ToDSL(), o.ToDSL())
	}
}

func valueNodeUnionJoinValueNode(n, o AstNode) (AstNode, error) {
	nn := n.(ValueNode)
	on := o.(ValueNode)
	if CompareAny(nn.getValue(), on.getValue(), nn.getVType().mType) == 0 {
		return n, nil
	} else {
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func valueNodeIntersectValueNode(n, o AstNode) (AstNode, error) {
	nn := n.(ValueNode)
	on := o.(ValueNode)
	if CompareAny(nn.getValue(), on.getValue(), nn.getVType().mType) == 0 {
		return n, nil
	} else if nn.getVType().aType {
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
	}
}

func checkCommonDslType(dslType DslType) bool {
	return dslType == EXISTS_DSL_TYPE ||
		dslType == BOOL_DSL_TYPE ||
		dslType == MATCH_ALL_DSL_TYPE ||
		dslType == EMPTY_DSL_TYPE
}
