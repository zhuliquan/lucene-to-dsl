package mapping

func checkTypeSupportLucene(typ FieldType) bool {
	_, ok := luceneSupportFieldType[typ]
	return ok
}

func strLstHasPrefix(va, vb []string) bool {
	if len(vb) > len(va) {
		return false
	}
	for i, n := 0, len(vb); i < n; i++ {
		if va[i] != vb[i] {
			return false
		}
	}
	return true
}
