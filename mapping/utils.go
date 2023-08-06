package mapping

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/zhuliquan/lucene-to-dsl/utils"
)

type BoolOrString interface {
	GetBool() bool
	GetString() string
}

type BoolValue bool

func (b BoolValue) GetBool() bool {
	return bool(b)
}

func (b BoolValue) GetString() string {
	return strconv.FormatBool(bool(b))
}

type StringValue string

func (s StringValue) GetBool() bool {
	b, _ := strconv.ParseBool(string(s))
	return b
}

func (s StringValue) GetString() string {
	return string(s)
}

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

func checkWildcard(path []string) bool {
	wildcard := false
	for _, p := range path {
		if strings.Contains(p, "*") || strings.Contains(p, "?") {
			return true
		}
	}
	return wildcard
}

func addProperty(res map[string]*Property, key string, prop *Property) error {
	if pre, ok := res[key]; ok {
		if pre.Type != prop.Type {
			return fmt.Errorf("field: %s have conflict type with %s and %s", key, pre.Type, prop.Type)
		}
	} else {
		res[key] = prop
	}
	return nil
}

func _getProperty(mpp map[string]*Property, index int, matchedPath, patternPath []string, res map[string]*Property) error {
	for cf, cp := range mpp {
		if cp.Type == ALIAS_FIELD_TYPE { // skip alias
			continue
		}

		currFieldPath := strings.Split(cf, ".")
		if !matchFieldPath(currFieldPath, patternPath, index) {
			continue
		}

		idxInc := len(currFieldPath)
		tempMatchedPath := make([]string, len(matchedPath))
		copy(tempMatchedPath, matchedPath)
		tempMatchedPath = append(tempMatchedPath, cf)
		if index+idxInc == len(patternPath) {
			switch cp.Type {
			case OBJECT_FIELD_TYPE, NESTED_FIELD_TYPE, FLATTENED_FIELD_TYPE:
				// expect not terminated here, so do nothing
			default:
				key := strings.Join(tempMatchedPath, ".")
				err := addProperty(res, key, cp)
				if err != nil {
					return err
				}
			}
			continue
		}

		if cp.Type == FLATTENED_FIELD_TYPE { // in flattened type, every leaf node is keyword type
			// flattened type doesn't support have wildcard
			if !checkWildcard(patternPath[index+idxInc:]) {
				key := strings.Join(patternPath, ".")
				err := addProperty(res, key, &Property{Type: KEYWORD_FIELD_TYPE})
				if err != nil {
					return err
				}
			}
			continue
		}

		for _, subProperties := range []map[string]*Property{
			cp.Properties, cp.Fields,
		} {
			err := _getProperty(subProperties, index+idxInc, tempMatchedPath, patternPath, res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getProperty(m *PropertyMapping, target string) (map[string]*Property, error) {
	patternPath := strings.Split(target, ".")
	var res = map[string]*Property{}
	err := _getProperty(m.fieldMapping.Properties, 0, []string{}, patternPath, res)
	return res, err
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
			if property, err := getProperty(pp, cp.Path); len(property) == 0 {
				return fmt.Errorf("filed: %s is alias, but can't find property for path: %s", fd, cp.Path)
			} else if err != nil {
				return err
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

func _fillDefaultParameter(pt map[string]*Property, pmt Dynamic) {
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
			if pt[cf].Dynamic == nil {
				pt[cf].Dynamic = pmt
			}
			if pt[cf].Type == "" {
				pt[cf].Type = OBJECT_FIELD_TYPE
			}
			_fillDefaultParameter(
				pt[cf].Properties,
				pt[cf].Dynamic,
			)
		}
	}
}

// updating mapping type recursively
func fillDefaultParameter(pm *PropertyMapping) {
	if pm.fieldMapping.Dynamic == nil {
		pm.fieldMapping.Dynamic = BoolDynamic(true)
	}
	_fillDefaultParameter(
		pm.fieldMapping.Properties,
		pm.fieldMapping.Dynamic,
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
