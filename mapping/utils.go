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

func completePropertyMapping(pm *PropertyMapping) error {
	return nil
}

// func flattenMapping(m *PropertyMapping) (map[string]*Property, error) {
// 	var mp = map[string]*Property{}
// 	var err = func(_m *Mapping, _mp map[string]*Property) error {
// 		for f, p := range m.Properties {
// 			switch p.Type {
// 			case ALIAS_FIELD_TYPE:
// 				if p.Path == "" {
// 					return fmt.Errorf("field: %s is alias, but not path parameter", f)
// 				} else {
// 					_aliasMap[f] = p.Path
// 				}
// 			case OBJECT_FIELD_TYPE, NESTED_FIELD_TYPE, FLATTENED_FIELD_TYPE:
// 				// if len(fields) == len(fs) {
// 				// 	return nil, fmt.Errorf("don't found")
// 				// } else {
// 				// 	if p.Properties == nil {
// 				// 		return p, nil
// 				// 	} else {
// 				// 		return GetProperty(strings.Join(fields[len(fs):], "."))
// 				// 	}
// 				// }
// 			default:
// 				// if len(p.Fields) != nil {
// 				// 	for f, p := range p.Fields {

// 				// 	}

// 				// }
// 			}
// 		}
// 	}(m, mp)
// 	return mp, err
// }
