package mapping

import (
	"fmt"
	"math"
	"strings"

	"github.com/zhuliquan/lucene-to-dsl/utils"
)

func checkTypeSupportLucene(typ FieldType) bool {
	_, ok := luceneSupportFieldType[typ]
	return ok
}

func matchFieldPath(partialPath []string, patternPath []string, index int) bool {
	if len(partialPath) > len(patternPath)-index {
		return false
	}
	for i := 0; i < len(partialPath); i++ {
		if !utils.WildcardMatch([]rune(partialPath[i]), []rune(
			strings.ReplaceAll(strings.ReplaceAll(patternPath[index+i], "\\*", "*"), "\\?", "?"))) {
			return false
		}
	}
	return true
}

func checkRestPathFlattenOk(path []string) bool {
	for _, p := range path {
		// flatten path can't contain any wildcard char.
		if strings.Contains(p, "*") || strings.Contains(p, "?") {
			return false
		}
	}
	return true
}

func _getProperty(mpp map[string]*Property, index int, matchedPath, patternPath []string, wildcard bool) map[string]*Property {
	res := map[string]*Property{}
	for cf, cp := range mpp {

		if !wildcard && len(res) > 0 { // normal field (no wildcard) find property need return
			break
		}

		if cp.Type == ALIAS_FIELD_TYPE {
			continue
		}

		partialPath := strings.Split(cf, ".")
		if !matchFieldPath(partialPath, patternPath, index) {
			continue
		}

		idxInc := len(partialPath)
		matchingPath := append(matchedPath, cf)
		if index+idxInc == len(patternPath) {
			switch cp.Type {
			case OBJECT_FIELD_TYPE, NESTED_FIELD_TYPE, FLATTENED_FIELD_TYPE:
				// do nothing
			default:
				res[strings.Join(matchingPath, ".")] = cp
			}
			continue
		}

		if cp.Type == FLATTENED_FIELD_TYPE { // support flattened type
			// pattern path must be specific, can't be fuzzy (i.g. \*.x, x.\*)
			if checkRestPathFlattenOk(patternPath) {
				res[strings.Join(patternPath, ".")] = cp
			}
			continue
		}

		for _, subProperties := range []map[string]*Property{
			cp.Properties, cp.Fields,
		} {
			for p, subRes := range _getProperty(subProperties, index+idxInc, matchingPath, patternPath, wildcard) {
				res[p] = subRes
			}
		}

		if cp.Type == OBJECT_FIELD_TYPE || cp.Type == NESTED_FIELD_TYPE {
			// object / nested
			res[strings.Join(patternPath, ".")] = cp
		}
	}
	return res
}

func getProperty(m *PropertyMapping, target string) map[string]*Property {
	patternPath := strings.Split(target, ".")
	isWildcardField := !checkRestPathFlattenOk(patternPath)
	return _getProperty(m.fieldMapping.Properties, 0, []string{}, patternPath, isWildcardField)
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
			if property := getProperty(pp, cp.Path); len(property) == 0 {
				return fmt.Errorf("filed: %s is alias, but can't find property for path: %s", fd, cp.Path)
			} else {
				pp.propertyCache[cp.Path] = property[cp.Path]
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
		if (pt[cf].Type == DATE_FIELD_TYPE || pt[cf].Type == DATE_RANGE_FIELD_TYPE) &&
			pt[cf].Format == "" {
			pt[cf].Format = "strict_date_optional_time||epoch_millis"
		}
		if pt[cf].Type == DATE_NANOS_FIELD_TYPE && pt[cf].Format == "" {
			pt[cf].Format = "strict_date_optional_time_nanos||epoch_millis"
		}
		if pt[cf].Type == SCALED_FLOAT_FIELD_TYPE && math.Abs(pt[cf].ScalingFactor-0.0) <= 1e-8 {
			pt[cf].ScalingFactor = 1.0
		}
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

func CheckNumberType(t FieldType) bool {
	return CheckIntType(t) || CheckUIntType(t) || CheckFloatType(t)

}

func CheckIntType(t FieldType) bool {
	return t == BYTE_FIELD_TYPE || t == SHORT_FIELD_TYPE ||
		t == INTEGER_FIELD_TYPE || t == INTEGER_RANGE_FIELD_TYPE ||
		t == LONG_FIELD_TYPE || t == LONG_RANGE_FIELD_TYPE
}

func CheckUIntType(t FieldType) bool {
	return t == UNSIGNED_LONG_FIELD_TYPE
}

func CheckFloatType(t FieldType) bool {
	return t == FLOAT_FIELD_TYPE || t == FLOAT_RANGE_FIELD_TYPE ||
		t == DOUBLE_FIELD_TYPE || t == DOUBLE_RANGE_FIELD_TYPE ||
		t == HALF_FLOAT_FIELD_TYPE || t == SCALED_FLOAT_FIELD_TYPE
}

func CheckDateType(t FieldType) bool {
	return t == DATE_FIELD_TYPE || t == DATE_NANOS_FIELD_TYPE || t == DATE_RANGE_FIELD_TYPE
}

func CheckVersionType(t FieldType) bool {
	return t == VERSION_FIELD_TYPE
}

func CheckIPType(t FieldType) bool {
	return t == IP_FIELD_TYPE || t == IP_RANGE_FIELD_TYPE
}

func CheckStringType(t FieldType) bool {
	return CheckKeywordType(t) || CheckTextType(t)
}

func CheckKeywordType(t FieldType) bool {
	return t == KEYWORD_FIELD_TYPE ||
		t == CONSTANT_KEYWORD_FIELD_TYPE ||
		t == WILDCARD_FIELD_TYPE
}

func CheckTextType(t FieldType) bool {
	return t == TEXT_FIELD_TYPE ||
		t == MATCH_ONLY_TEXT_FIELD_TYPE
}
