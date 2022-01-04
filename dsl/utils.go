package dsl

import (
	"sort"
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
