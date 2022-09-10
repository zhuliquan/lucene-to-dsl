package mapping

import (
	"fmt"
	"strings"
)

func checkTypeSupportLucene(typ FieldType) bool {
	_, ok := luceneSupportFieldType[typ]
	return ok
}

func matchFieldPath(matchPath []string, targetPath []string, index int) bool {
	if len(matchPath) > len(targetPath)-index {
		return false
	}
	for i := 0; i < len(matchPath); i++ {
		if matchPath[i] != targetPath[index+i] {
			return false
		}
	}
	return true
}

func _getProperty(mpp map[string]*Property, index int, paths []string) (*Property, error) {
	for cf, cp := range mpp {
		if cp.Type == ALIAS_FIELD_TYPE {
			continue
		}
		if !matchFieldPath(strings.Split(cf, "."), paths, index) {
			continue
		}

		if index == len(paths)-1 {
			switch cp.Type {
			case OBJECT_FIELD_TYPE, NESTED_FIELD_TYPE:
				return nil, fmt.Errorf("field: %s is not fully field path", strings.Join(paths, "."))
			default:
				return cp, nil
			}
		}

		if len(cp.Properties) != 0 {
			if p, err := _getProperty(cp.Properties, index+1, paths); err != nil {
				if strings.HasSuffix(err.Error(), "is not fully field path") {
					return nil, err
				}
			} else {
				return p, nil
			}
		}
	}
	return nil, fmt.Errorf("don't found field: %s in mapping", strings.Join(paths, "."))
}

func getProperty(m *PropertyMapping, target string) (*Property, error) {
	return _getProperty(m.fieldMapping.Properties, 0, strings.Split(target, "."))
}

func flattenAlias(pt map[string]*Property, pf string, am map[string]string, pp *PropertyMapping) error {
	for cf, cp := range pt {
		var fd string
		if pf != "" {
			fd = fmt.Sprintf("%s.%s", pf, cf)
		} else {
			fd = cf
		}
		switch cp.Type {
		case ALIAS_FIELD_TYPE:
			if cp.Path == "" {
				return fmt.Errorf("field: %s is alias, but not path parameter", fd)
			} else if cp.Path == fd {
				return fmt.Errorf("field: %s is alias, but path is same", fd)
			}
			if property, err := getProperty(pp, cp.Path); err != nil {
				return fmt.Errorf("filed: %s is alias, but can't find property for path: %s", fd, cp.Path)
			} else {
				pp.propertyCache[cp.Path] = property
			}
			am[fd] = cp.Path
		default:
			if len(cp.Properties) != 0 {
				if err := flattenAlias(cp.Properties, fd, am, pp); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func extractFieldAliasMap(pm *PropertyMapping) (map[string]string, error) {
	var (
		amm = map[string]string{}
		ppt = pm.fieldMapping.Properties
	)
	if err := flattenAlias(ppt, "", amm, pm); err != nil {
		return nil, err
	} else {
		return amm, nil
	}
}

func _fillDefaultParameter(pt map[string]*Property, pmt MappingType) {
	for cf := range pt {
		if len(pt[cf].Properties) != 0 {
			if pt[cf].MappingType == "" {
				pt[cf].MappingType = pmt
			}
			if pt[cf].Type == "" {
				pt[cf].Type = OBJECT_FIELD_TYPE
			}
			_fillDefaultParameter(
				pt[cf].Properties,
				pt[cf].MappingType,
			)
		}
	}
}

// updating mapping type recursively
func fillDefaultParameter(pm *PropertyMapping) {
	if pm.fieldMapping.MappingType == "" {
		pm.fieldMapping.MappingType = DYNAMIC_MAPPING
	}
	_fillDefaultParameter(
		pm.fieldMapping.Properties,
		pm.fieldMapping.MappingType,
	)
}
