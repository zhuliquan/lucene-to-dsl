package convert

import (
	"testing"

	mapping "github.com/zhuliquan/es-mapping"
)

func TestInferFieldType(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected mapping.FieldType
	}{
		{
			name:     "boolean true",
			value:    "true",
			expected: mapping.BOOLEAN_FIELD_TYPE,
		},
		{
			name:     "boolean false",
			value:    "false",
			expected: mapping.BOOLEAN_FIELD_TYPE,
		},
		{
			name:     "boolean True (case insensitive)",
			value:    "True",
			expected: mapping.BOOLEAN_FIELD_TYPE,
		},
		{
			name:     "integer",
			value:    "123",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
		{
			name:     "negative integer",
			value:    "-456",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
		{
			name:     "float",
			value:    "3.14",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
		{
			name:     "negative float",
			value:    "-2.5",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
		{
			name:     "date YYYY-MM-DD",
			value:    "2021-01-01",
			expected: mapping.DATE_FIELD_TYPE,
		},
		{
			name:     "date with time",
			value:    "2021-01-01T12:00:00",
			expected: mapping.DATE_FIELD_TYPE,
		},
		{
			name:     "date with timezone",
			value:    "2021-01-01T12:00:00Z",
			expected: mapping.DATE_FIELD_TYPE,
		},
		{
			name:     "IPv4 address",
			value:    "192.168.1.1",
			expected: mapping.IP_FIELD_TYPE,
		},
		{
			name:     "IPv6 address",
			value:    "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: mapping.IP_FIELD_TYPE,
		},
		{
			name:     "CIDR notation",
			value:    "192.168.0.0/24",
			expected: mapping.IP_FIELD_TYPE,
		},
		{
			name:     "plain string",
			value:    "hello",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
		{
			name:     "mixed string",
			value:    "hello123",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
		{
			name:     "string with quotes",
			value:    "\"hello\"",
			expected: mapping.KEYWORD_FIELD_TYPE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InferFieldType(tt.value)
			if result != tt.expected {
				t.Errorf("InferFieldType(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestCreateDefaultProperty(t *testing.T) {
	prop := CreateDefaultProperty(mapping.KEYWORD_FIELD_TYPE)
	if prop.Type != mapping.KEYWORD_FIELD_TYPE {
		t.Errorf("CreateDefaultProperty() type = %v, want %v", prop.Type, mapping.KEYWORD_FIELD_TYPE)
	}

	prop2 := CreateDefaultProperty(mapping.DATE_FIELD_TYPE)
	if prop2.Type != mapping.DATE_FIELD_TYPE {
		t.Errorf("CreateDefaultProperty() type = %v, want %v", prop2.Type, mapping.DATE_FIELD_TYPE)
	}
}
