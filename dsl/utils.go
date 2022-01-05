package dsl

import (
	"bytes"
	"math"
	"net"
	"sort"
	"time"
)

// union join two string slice
func UnionJoinStrLst(al, bl []string) []string {
	sort.Strings(al)
	sort.Strings(bl)
	var cl = make([]string, 0, len(al)+len(bl))
	var i, j, na, nb = 0, 0, len(al), len(bl)

	for i < na || j < nb {
		if i == na || (j < nb && al[i] > bl[j]) {
			cl = append(cl, bl[j])
			j += 1
		} else {
			cl = append(cl, al[i])
			i += 1
		}
	}
	return UniqStrLst(cl)
}

func IntersectStrLst(al, bl []string) []string {
	sort.Strings(al)
	sort.Strings(bl)
	var cl = make([]string, 0, len(al)+len(bl))
	var i, j, na, nb = 0, 0, len(al), len(bl)

	for i < na && j < nb {
		if al[i] == bl[j] {
			cl = append(cl, al[i])
			i += 1
			j += 1
		} else if al[i] > bl[j] {
			j += 1
		} else {
			i += 1
		}
	}
	return UniqStrLst(cl)
}

// uniq a sort string slice
func UniqStrLst(a []string) []string {
	if len(a) == 0 || len(a) == 1 {
		return a
	}
	var r = []string{}
	for i, n := 0, len(a); i < n; i++ {
		if i == 0 {
			r = append(r, a[i])
		} else if a[i] != a[i-1] {
			r = append(r, a[i])
		}
	}
	return r
}

// negative  mean a < b
// positive  mean a > b
// zero      mean a == b
// using nil represent inf
func CompareAny(a, b interface{}, typ NodeValue) int {
	if a != nil && b == nil {
		return -1
	} else if a == nil && b != nil {
		return 1
	} else if a == nil && b == nil {
		return 0
	}

	switch typ {
	case INT_VALUE:
		return a.(int) - b.(int)
	case FLOAT_VALUE:
		var af = a.(float64)
		var bf = b.(float64)
		if math.Abs(af-bf) < 1E-6 {
			return 0
		} else if af < bf {
			return -1
		} else {
			return 1
		}
	case DATE_VALUE:
		var at = a.(time.Time)
		var bt = b.(time.Time)
		if at.UnixNano() == bt.UnixNano() {
			return 0
		} else if at.Before(bt) {
			return -1
		} else {
			return 1
		}

	case IP_VALUE:
		var ai = a.(net.IP)
		var bi = b.(net.IP)
		return bytes.Compare(ai, bi)

	case KEYWORD_VALUE, PHRASE_VALUE:
		var as = a.(string)
		var bs = b.(string)
		if as > bs {
			return 1
		} else if as < bs {
			return -1
		} else {
			return 0
		}
	default:
		return 0
	}
}
