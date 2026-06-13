package convert

import (
	"net"
	"strconv"
	"strings"
	"time"

	mapping "github.com/zhuliquan/es-mapping"
)

var dateFormats = []string{
	"2006-01-02",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z07:00",
	"2006/01/02",
	"01/02/2006",
	"02/01/2006",
}

// InferFieldType infers field type based on the value content.
// Returns keyword as default for most cases.
func InferFieldType(value string) mapping.FieldType {
	value = strings.Trim(value, "\"'")

	if isBoolean(value) {
		return mapping.BOOLEAN_FIELD_TYPE
	}

	if isInteger(value) || isFloat(value) {
		return mapping.KEYWORD_FIELD_TYPE
	}

	if isDate(value) {
		return mapping.DATE_FIELD_TYPE
	}

	if isIP(value) {
		return mapping.IP_FIELD_TYPE
	}

	return mapping.KEYWORD_FIELD_TYPE
}

func isBoolean(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false"
}

func isInteger(value string) bool {
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil
}

func isFloat(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func isDate(value string) bool {
	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func isIP(value string) bool {
	if net.ParseIP(value) != nil {
		return true
	}
	if _, _, err := net.ParseCIDR(value); err == nil {
		return true
	}
	return false
}

// CreateDefaultProperty creates a default Property with inferred type.
func CreateDefaultProperty(inferredType mapping.FieldType) *mapping.Property {
	return &mapping.Property{
		Type: inferredType,
	}
}
