package dsl

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
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
func CompareAny(a, b *LeafValue, typ mapping.FieldType) int {
	var ret = 0
	switch typ {
	case mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE:
		fmt.Println("enter tiny int")
		ret = int(a.TinyInt - b.TinyInt)
	case mapping.HALF_FLOAT_FIELD_TYPE:
		if float32(a.Float16-b.Float16) < 1E-6 {
			return 0
		} else if a.Float16 < b.Float16 {
			return -1
		} else {
			return 1
		}
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		ret = int(a.LongInt - b.LongInt)
	case mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE,
		mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE,
		mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE,
		mapping.DOUBLE_FIELD_TYPE, mapping.DOUBLE_RANGE_FIELD_TYPE, mapping.SCALED_FLOAT_FIELD_TYPE:
		return a.Decimal.Cmp(b.Decimal)
	case mapping.KEYWORD_FIELD_TYPE, mapping.TEXT_FIELD_TYPE, mapping.WILDCARD_FIELD_TYPE, mapping.CONSTANT_KEYWORD_FIELD_TYPE:
		var as = a.String
		var bs = b.String
		if as > bs {
			return 1
		} else if as < bs {
			return -1
		} else {
			return 0
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
		var at = a.DateTime
		var bt = b.DateTime
		if at.UnixNano() == bt.UnixNano() {
			return 0
		} else if at.Before(bt) {
			return -1
		} else {
			return 1
		}
	case mapping.VERSION_FIELD_TYPE:
		var av = a.Version
		var bv = b.Version
		return av.Compare(bv)
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		var ai = []byte(a.IpValue)
		var bi = []byte(b.IpValue)
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

func compareInt64(a, b int64, c CompareType) int64 {
	switch c {
	case LT:
		return ltInt64(a, b)
	case GT:
		return gtInt64(a, b)
	case LTE:
		return lteInt64(a, b)
	case GTE:
		return gteInt64(a, b)
	default:
		return a
	}
}

func ltInt64(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteInt64(a, b int64) int64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtInt64(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteInt64(a, b int64) int64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func compareUInt64(a, b uint64, c CompareType) uint64 {
	switch c {
	case LT:
		return ltUInt64(a, b)
	case GT:
		return gtUInt64(a, b)
	case LTE:
		return lteUInt64(a, b)
	case GTE:
		return gteUInt64(a, b)
	default:
		return a
	}
}

func ltUInt64(a, b uint64) uint64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteUInt64(a, b uint64) uint64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtUInt64(a, b uint64) uint64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteUInt64(a, b uint64) uint64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func compareFloat64(a, b float64, c CompareType) float64 {
	switch c {
	case LT:
		return ltFloat64(a, b)
	case GT:
		return gtFloat64(a, b)
	case LTE:
		return lteFloat64(a, b)
	case GTE:
		return gteFloat64(a, b)
	default:
		return a
	}
}

func ltFloat64(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteFloat64(a, b float64) float64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtFloat64(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteFloat64(a, b float64) float64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func compareIp(a, b net.IP, c CompareType) net.IP {
	switch c {
	case LT:
		return ltIP(a, b)
	case GT:
		return gtIP(a, b)
	case LTE:
		return lteIP(a, b)
	case GTE:
		return gteIP(a, b)
	default:
		return a
	}
}

func ltIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res < 0 {
		return a
	} else {
		return b
	}
}

func lteIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res <= 0 {
		return a
	} else {
		return b
	}
}

func gtIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res > 0 {
		return a
	} else {
		return b
	}
}

func gteIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res >= 0 {
		return a
	} else {
		return b
	}
}

func compareDate(a, b time.Time, c CompareType) time.Time {
	switch c {
	case LT:
		return ltDate(a, b)
	case GT:
		return gtDate(a, b)
	case LTE:
		return lteDate(a, b)
	case GTE:
		return gteDate(a, b)
	default:
		return a
	}
}

func ltDate(a, b time.Time) time.Time {
	if a.UnixNano() < b.UnixNano() {
		return a
	} else {
		return b
	}
}

func lteDate(a, b time.Time) time.Time {
	if a.UnixNano() <= b.UnixNano() {
		return a
	} else {
		return b
	}
}

func gtDate(a, b time.Time) time.Time {
	if a.UnixNano() > b.UnixNano() {
		return a
	} else {
		return b
	}
}

func gteDate(a, b time.Time) time.Time {
	if a.UnixNano() >= b.UnixNano() {
		return a
	} else {
		return b
	}
}

func compareString(a, b string, c CompareType) string {
	switch c {
	case LT:
		return ltString(a, b)
	case GT:
		return gtString(a, b)
	case LTE:
		return lteString(a, b)
	case GTE:
		return gteString(a, b)
	default:
		return a
	}
}

func ltString(a, b string) string {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteString(a, b string) string {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtString(a, b string) string {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteString(a, b string) string {
	if a >= b {
		return a
	} else {
		return b
	}
}
