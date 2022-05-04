package dsl

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/shopspring/decimal"
	"github.com/x448/float16"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func valueLstToStrLst(vl []LeafValue) []string {
	var rl = make([]string, len(vl), len(vl))
	for i, v := range vl {
		rl[i] = v.(string)
	}
	return rl
}

func strLstToValueLst(vl []string) []LeafValue {
	var rl = make([]LeafValue, len(vl), len(vl))
	for i, v := range vl {
		rl[i] = v
	}
	return rl
}

// union join two string slice
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

// negative  mean a < b
// positive  mean a > b
// zero      mean a == b
// using nil represent inf
func CompareAny(a, b LeafValue, typ mapping.FieldType) int {
	var ret = 0
	switch typ {
	case mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE,
		mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE,
		mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE:
		ret = int(a.(int64) - b.(int64))
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		ret = int(a.(uint64) - b.(uint64))
	case mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE,
		mapping.DOUBLE_FIELD_TYPE, mapping.DOUBLE_RANGE_FIELD_TYPE:
		var af = a.(float64)
		var bf = b.(float64)
		if math.Abs(af-bf) <= eps {
			return 0
		} else if af < bf {
			return -1
		} else {
			return 1
		}
	case mapping.HALF_FLOAT_FIELD_TYPE:
		var af = a.(float16.Float16).Float32()
		var bf = b.(float16.Float16).Float32()
		if math.Abs(float64(af-bf)) <= eps {
			return 0
		} else if af < bf {
			return -1
		} else {
			return 1
		}
	case mapping.SCALED_FLOAT_FIELD_TYPE:
		var ad = a.(decimal.Decimal)
		var bd = b.(decimal.Decimal)
		return ad.Cmp(bd)
	case mapping.KEYWORD_FIELD_TYPE, mapping.TEXT_FIELD_TYPE,
		mapping.WILDCARD_FIELD_TYPE, mapping.CONSTANT_KEYWORD_FIELD_TYPE:
		var as = a.(string)
		var bs = b.(string)
		if as > bs {
			return 1
		} else if as < bs {
			return -1
		} else {
			return 0
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
		var at = a.(time.Time)
		var bt = b.(time.Time)
		if at.UnixNano() == bt.UnixNano() {
			return 0
		} else if at.Before(bt) {
			return -1
		} else {
			return 1
		}
	case mapping.VERSION_FIELD_TYPE:
		var av = a.(*version.Version)
		var bv = b.(*version.Version)
		return av.Compare(bv)
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		var ai = []byte(a.(net.IP))
		var bi = []byte(b.(net.IP))
		return bytes.Compare(ai, bi)
	default:
		return 0
	}
	if ret > 0 {
		fmt.Println(a, b)
		return 1
	} else if ret < 0 {
		fmt.Println(a, b)
		return -1
	} else {
		return 0
	}
}

func CheckValidRangeNode(node *RangeNode) error {
	var cmp = CompareAny(node.LeftValue, node.RightValue, node.ValueType)
	if cmp > 0 || (cmp == 0 && (node.LeftCmpSym == GT || node.RightCmpSym == LT)) {
		return fmt.Errorf("range is conflict")
	}
	return nil
}
