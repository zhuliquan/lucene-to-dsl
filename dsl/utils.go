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
func lfNodeUnionJoinLfNode(a, b AstNode) (AstNode, error) {
	return &OrNode{
		MinimumShouldMatch: 1,
		Nodes: map[string][]AstNode{
			a.NodeKey(): {a, b},
		},
	}, nil
}

// intersect two leaf node
func lfNodeIntersectLfNode(a, b AstNode) (AstNode, error) {
	af := a.(FilterCtxNode)
	bf := b.(FilterCtxNode)
	var mustNodes = map[string][]AstNode{
		a.NodeKey(): {},
	}
	var filterNodes = map[string][]AstNode{
		a.NodeKey(): {},
	}
	if af.getFilterCtx() {
		filterNodes[a.NodeKey()] = append(
			filterNodes[a.NodeKey()], a,
		)
	} else {
		mustNodes[a.NodeKey()] = append(
			mustNodes[a.NodeKey()], a,
		)
	}
	if bf.getFilterCtx() {
		filterNodes[a.NodeKey()] = append(
			filterNodes[a.NodeKey()], b,
		)
	} else {
		mustNodes[a.NodeKey()] = append(
			mustNodes[a.NodeKey()], b,
		)
	}

	var andNode = &AndNode{}
	if len(mustNodes[a.NodeKey()]) != 0 {
		andNode.MustNodes = mustNodes
	}
	if len(filterNodes[a.NodeKey()]) != 0 {
		andNode.FilterNodes = filterNodes
	}
	return andNode, nil
}

func flattenNodes(nodesMap map[string][]AstNode) interface{} {
	var dslRes = []DSL{}
	for _, nodes := range nodesMap {
		for _, node := range nodes {
			dslRes = append(dslRes, node.ToDSL())
		}
	}
	if len(dslRes) == 0 {
		return nil
	} else if len(dslRes) == 1 {
		return dslRes[0]
	} else {
		return dslRes
	}
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

func astNodeUnionJoinTermsNode(n AstNode, o *TermsNode, excludes []LeafValue) (AstNode, error) {
	if len(excludes) == 0 {
		return n, nil
	} else if len(excludes) == 1 {
		return lfNodeUnionJoinLfNode(n, &TermNode{
			kvNode: kvNode{
				fieldNode: o.fieldNode,
				valueNode: valueNode{valueType: o.valueType, value: excludes[0]},
			},
			boostNode: o.boostNode,
		})
	} else {
		return lfNodeUnionJoinLfNode(n,
			&TermsNode{
				fieldNode: o.fieldNode,
				valueType: o.valueType,
				boostNode: o.boostNode,
				terms:     excludes,
			},
		)
	}
}

func astNodeIntersectTermsNode(n AstNode, o *TermsNode, excludes []LeafValue) (AstNode, error) {
	if len(excludes) == 0 {
		return n, nil
	} else if len(excludes) == 1 {
		return lfNodeIntersectLfNode(n, &TermNode{
			kvNode: kvNode{
				fieldNode: o.fieldNode,
				valueNode: valueNode{valueType: o.valueType, value: excludes[0]},
			},
			boostNode: o.boostNode,
		})
	} else {
		return lfNodeIntersectLfNode(n,
			&TermsNode{
				fieldNode: o.fieldNode,
				valueType: o.valueType,
				boostNode: o.boostNode,
				terms:     excludes,
			},
		)
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

func termsToPrintValue(terms []LeafValue, t mapping.FieldType) interface{} {
	x := []interface{}{}
	for _, term := range terms {
		x = append(x, leafValueToPrintValue(term, t))
	}
	return x
}

func compareBoost(a, b BoostNode) int {
	return CompareAny(a.getBoost(), b.getBoost(), mapping.DOUBLE_FIELD_TYPE)
}

// wildcard match text and pattern
func wildcardMatch(text []rune, pattern []rune) bool {
	var n, m = len(text), len(pattern)
	var dp = make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]bool, m+1)
	}

	dp[0][0] = true
	for i := 1; i <= n; i++ {
		dp[i][0] = false
	}
	for j := 1; j <= m; j++ {
		if pattern[j-1] == '*' {
			dp[0][j] = dp[0][j-1]
		}
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if pattern[j-1] == '?' || pattern[j-1] == text[i-1] {
				dp[i][j] = dp[i-1][j-1]
			} else if pattern[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else {
				dp[i][j] = false
			}
		}
	}
	return dp[n][m]
}

func patternNodeUnionJoinTermNode(n PatternMatcher, o *TermNode) (AstNode, error) {
	if n.Match([]byte(o.value.(string))) {
		return n.(AstNode), nil
	} else {
		return lfNodeUnionJoinLfNode(n.(AstNode), o)
	}
}

func patternNodeUnionJoinTermsNode(n PatternMatcher, o *TermsNode) (AstNode, error) {
	var excludes = []LeafValue{}
	for _, term := range o.terms {
		if !n.Match([]byte(term.(string))) {
			excludes = append(excludes, term)
		}
	}
	return astNodeUnionJoinTermsNode(n.(AstNode), o, excludes)
}

func patternNodeIntersectTermNode(n PatternMatcher, o *TermNode) (AstNode, error) {
	if n.Match([]byte(o.value.(string))) {
		return o, nil
	} else if n.(ArrayTypeNode).isArrayType() {
		return lfNodeIntersectLfNode(n.(AstNode), o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.(AstNode).ToDSL(), o.ToDSL())
	}
}

func patternNodeIntersectTermsNode(n PatternMatcher, o *TermsNode) (AstNode, error) {
	if n.(ArrayTypeNode).isArrayType() {
		var excludes = []LeafValue{}
		for _, term := range o.terms {
			if !n.Match([]byte(term.(string))) {
				excludes = append(excludes, term)
			}
		}
		return astNodeIntersectTermsNode(n.(AstNode), o, excludes)
	} else {
		var includes = []LeafValue{}
		for _, term := range o.terms {
			if n.Match([]byte(term.(string))) {
				includes = append(includes, term)
			}
		}
		if len(includes) == 0 {
			return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.(AstNode).ToDSL(), o.ToDSL())
		} else if len(includes) == 1 {
			return &TermNode{
				kvNode: kvNode{
					fieldNode: o.fieldNode,
					valueNode: valueNode{
						valueType: o.valueType,
						value:     includes[0],
					},
				},
				boostNode: o.boostNode,
			}, nil
		} else {
			return &TermsNode{
				fieldNode: o.fieldNode,
				boostNode: o.boostNode,
				valueType: o.valueType,
				terms:     includes,
			}, nil
		}
	}
}

func valueNodeUnionJoinValueNode(n, o AstNode) (AstNode, error) {
	nn := n.(ValueNode)
	on := n.(ValueNode)
	if CompareAny(nn.getValue(), on.getValue(), nn.getVType().mType) == 0 {
		return n, nil
	} else {
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func valueNodeIntersectValueNode(n, o AstNode) (AstNode, error) {
	nn := n.(ValueNode)
	on := n.(ValueNode)
	if CompareAny(nn.getValue(), on.getValue(), nn.getVType().mType) == 0 {
		return n, nil
	} else if nn.getVType().aType {
		return lfNodeIntersectLfNode(n, o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
	}
}
