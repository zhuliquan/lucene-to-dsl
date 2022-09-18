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
	switch x.(type) {
	case int:
		return uint64(x.(int))
	case uint:
		return uint64(x.(uint))
	default:
		return x.(uint64)
	}
}

func castInt(x LeafValue) int64 {
	switch x.(type) {
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
	var cmp = CompareAny(node.LeftValue, node.RightValue, node.Type)
	if cmp > 0 || (cmp == 0 && (node.LeftCmpSym == GT || node.RightCmpSym == LT)) {
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
