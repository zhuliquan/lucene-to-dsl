package lucene_to_dsl

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
	"github.com/zhuliquan/lucene-to-dsl/convert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
)

func mustDSL(jsonStr string) dsl.DSL {
	var raw interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		panic(fmt.Sprintf("invalid DSL JSON: %v", err))
	}
	return convertToDSL(raw).(dsl.DSL)
}

func convertToDSL(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		d := make(dsl.DSL)
		for k, v2 := range val {
			d[k] = convertToDSL(v2)
		}
		return d
	case []interface{}:
		for i, v2 := range val {
			val[i] = convertToDSL(v2)
		}
		return val
	default:
		return v
	}
}

func assertDSLEqual(t *testing.T, expected dsl.DSL, actual dsl.DSL) {
	t.Helper()
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	assert.Equal(t, string(expectedJSON), string(actualJSON))
}

var mappingJSON = []byte(`{
  "properties": {
    "status": {"type": "keyword"},
    "title": {"type": "text"},
    "count": {"type": "integer"},
    "price": {"type": "float"},
    "is_active": {"type": "boolean"},
    "created_at": {"type": "date"},
    "ip_address": {"type": "ip"},
    "tags": {"type": "keyword"},
    "description": {"type": "text"},
    "level": {"type": "byte"},
    "weight": {"type": "half_float"},
    "uuid": {"type": "wildcard"}
  }
}`)

func TestLuceneToDSL_WithMapping(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		// Basic term queries
		{"keyword_term", `status:active`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
		{"integer_term", `count:100`, mustDSL(`{"term":{"count":{"boost":1,"value":100}}}`), false},
		{"boolean_term", `is_active:true`, mustDSL(`{"term":{"is_active":{"boost":1,"value":true}}}`), false},
		{"float_term", `price:19.99`, mustDSL(`{"term":{"price":{"boost":1,"value":19.989999771118164}}}`), false},
		{"ip_term", `ip_address:192.168.1.1`, mustDSL(`{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}}`), false},
		{"byte_term", `level:5`, mustDSL(`{"term":{"level":{"boost":1,"value":5}}}`), false},
		{"half_float_term", `weight:1.5`, mustDSL(`{"term":{"weight":{"boost":1,"value":1.5}}}`), false},
		{"wildcard_field", `uuid:*abc*`, mustDSL(`{"wildcard":{"uuid":{"boost":1,"rewrite":"constant_score","value":"*abc*"}}}`), false},

		// Range queries
		{"closed_range", `count:[10 TO 100]`, mustDSL(`{"range":{"count":{"boost":1,"gte":10,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"open_range", `count:{10 TO 100}`, mustDSL(`{"range":{"count":{"boost":1,"gt":10,"lt":100,"relation":"INTERSECTS"}}}`), false},
		{"gt", `count:>10`, mustDSL(`{"range":{"count":{"boost":1,"gt":10,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},
		{"gte", `count:>=10`, mustDSL(`{"range":{"count":{"boost":1,"gte":10,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},
		{"lt", `count:<100`, mustDSL(`{"range":{"count":{"boost":1,"gt":-2147483648,"lt":100,"relation":"INTERSECTS"}}}`), false},
		{"lte", `count:<=100`, mustDSL(`{"range":{"count":{"boost":1,"gt":-2147483648,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"date_range", `created_at:[2021-01-01 TO 2021-12-31]`, mustDSL(`{"range":{"created_at":{"boost":1,"format":"epoch_millis","gte":1609459200000,"lte":1640908800000,"relation":"INTERSECTS"}}}`), false},
		{"ip_range", `ip_address:[192.168.0.0 TO 192.168.255.255]`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.0","lte":"192.168.255.255","relation":"INTERSECTS"}}}`), false},
		{"invalid_range", `count:[100 TO 10]`, nil, true},

		// String operations
		{"prefix", `status:act*`, mustDSL(`{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}`), false},
		{"wildcard", `status:act*ve`, mustDSL(`{"wildcard":{"status":{"boost":1,"rewrite":"constant_score","value":"act*ve"}}}`), false},
		{"fuzzy", `status:active~2`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
		{"regexp", `status:/act.*/`, mustDSL(`{"regexp":{"status":{"flags":"ALL","max_determinized_states":10000,"rewrite":"constant_score","value":"act.*"}}}`), false},
		{"match_phrase", `title:"hello world"`, mustDSL(`{"match_phrase":{"title":{"boost":1,"query":"hello world"}}}`), false},
		{"match_query", `title:hello`, mustDSL(`{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}`), false},

		// Boolean operations (same field type)
		{"or_op", `status:active OR status:pending`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"not_op", `NOT status:inactive`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"term":{"status":{"boost":1,"value":"inactive"}}}}}`), false},
		{"or_with_same_field", `status:active OR status:pending OR status:inactive`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}},{"term":{"status":{"boost":1,"value":"inactive"}}}]}}`), false},
		{"and_with_same_field", `status:active AND status:pending`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"not_with_range", `NOT count:<100`, mustDSL(`{"range":{"count":{"boost":1,"gte":100,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},

		// Boost
		{"boost_term", `status:active^2.0`, mustDSL(`{"term":{"status":{"boost":2,"value":"active"}}}`), false},
		{"boost_range", `count:[10 TO 100]^1.5`, mustDSL(`{"range":{"count":{"boost":1.5,"gte":10,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"boost_match_phrase", `title:"hello world"^1.5`, mustDSL(`{"match_phrase":{"title":{"boost":1.5,"query":"hello world"}}}`), false},
		{"boost_prefix", `status:act*^1.2`, mustDSL(`{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}`), false},

		// Special queries
		{"exists", `_exists_:status`, mustDSL(`{"exists":{"field":"status"}}`), false},
		{"exists_or", `_exists_:status OR _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},

		// Text field behavior
		{"text_single", `title:hello`, mustDSL(`{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}`), false},
		{"text_phrase", `title:"hello world"`, mustDSL(`{"match_phrase":{"title":{"boost":1,"query":"hello world"}}}`), false},
		{"text_boost", `title:hello^2.0`, mustDSL(`{"match":{"title":{"boost":2,"max_expansions":50,"query":"hello"}}}`), false},
		{"text_phrase_boost", `title:"hello world"^1.5`, mustDSL(`{"match_phrase":{"title":{"boost":1.5,"query":"hello world"}}}`), false},
		{"text_or", `title:hello OR title:world`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"match":{"title":{"boost":1,"max_expansions":50,"query":"world"}}}]}}`), false},

		// IP field behavior
		{"ip_exact", `ip_address:192.168.1.1`, mustDSL(`{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}}`), false},
		{"ip_cidr", `ip_address:192.168.0.0/24`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.1","lte":"192.168.0.254","relation":"INTERSECTS"}}}`), false},
		{"ip_range_query", `ip_address:[192.168.0.0 TO 192.168.255.255]`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.0","lte":"192.168.255.255","relation":"INTERSECTS"}}}`), false},

		// Edge cases - error cases
		{"empty_query", ``, nil, true},
		{"invalid_syntax", `status: AND count:`, nil, true},
		{"invalid_mapping", `status:active`, nil, true},
		{"field_not_in_mapping", `unknown_field:value`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []Option
			if tt.name == "invalid_mapping" {
				opts = []Option{WithMappingData([]byte(`{invalid json}`))}
			} else {
				opts = []Option{WithMappingData(mappingJSON)}
			}
			got, err := LuceneToDSL(tt.query, opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_WithoutMapping(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"boolean_true", `active:true`, mustDSL(`{"term":{"active":{"boost":1,"value":true}}}`), false},
		{"boolean_false", `active:false`, mustDSL(`{"term":{"active":{"boost":1,"value":false}}}`), false},
		{"integer", `count:123`, mustDSL(`{"term":{"count":{"boost":1,"value":"123"}}}`), false},
		{"negative_integer", `count:-456`, mustDSL(`{"term":{"count":{"boost":1,"value":"-456"}}}`), false},
		{"float", `price:3.14`, mustDSL(`{"term":{"price":{"boost":1,"value":"3.14"}}}`), false},
		{"date", `created_at:2021-01-01`, mustDSL(`{"range":{"created_at":{"boost":1,"format":"epoch_millis","gte":1609459200000,"lte":1640995199999,"relation":"INTERSECTS"}}}`), false},
		{"ipv4", `ip:192.168.1.1`, mustDSL(`{"term":{"ip":{"boost":1,"value":"192.168.1.1"}}}`), false},
		{"ipv4_cidr", `ip:192.168.0.0/24`, mustDSL(`{"range":{"ip":{"boost":1,"gte":"192.168.0.1","lte":"192.168.0.254","relation":"INTERSECTS"}}}`), false},
		{"keyword", `status:active`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
		{"and_op", `status:active AND count:>100`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":"100","lt":"\ufffd","relation":"INTERSECTS"}}}]}}`), false},
		{"or_op", `status:active OR status:pending`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"not_op", `NOT status:inactive`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"term":{"status":{"boost":1,"value":"inactive"}}}}}`), false},
		{"complex", `(status:active OR status:pending) AND count:>100`, mustDSL(`{"bool":{"minimum_should_match":1,"must":{"range":{"count":{"boost":1,"gt":"100","lt":"\ufffd","relation":"INTERSECTS"}}},"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"prefix", `status:act*`, mustDSL(`{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}`), false},
		{"wildcard", `status:act*ve`, mustDSL(`{"wildcard":{"status":{"boost":1,"rewrite":"constant_score","value":"act*ve"}}}`), false},
		{"fuzzy", `status:active~2`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
		{"regexp", `status:/act.*/`, mustDSL(`{"regexp":{"status":{"flags":"ALL","max_determinized_states":10000,"rewrite":"constant_score","value":"act.*"}}}`), false},
		{"exists", `_exists_:field`, mustDSL(`{"exists":{"field":"field"}}`), false},
		{"boost", `status:active^2.0`, mustDSL(`{"term":{"status":{"boost":2,"value":"active"}}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_WithCustomConvertFunc(t *testing.T) {
	customFuncs := map[string]convert.ConvertFunc{
		"title": func(val interface{}, props mapping.ExtProperties) (interface{}, error) {
			if str, ok := val.(string); ok {
				return fmt.Sprintf("[%s]", str), nil
			}
			return val, nil
		},
	}

	got, err := LuceneToDSL(
		`title:hello`,
		WithMappingData(mappingJSON),
		WithCustomConvertFunc(customFuncs),
	)
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func TestLuceneToDSL_Options(t *testing.T) {
	t.Run("with_mapping", func(t *testing.T) {
		got, err := LuceneToDSL(`status:active`, WithMappingData(mappingJSON))
		assert.NoError(t, err)
		assertDSLEqual(t, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), got)
	})

	t.Run("without_mapping", func(t *testing.T) {
		got, err := LuceneToDSL(`status:active`)
		assert.NoError(t, err)
		assertDSLEqual(t, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), got)
	})

	t.Run("with_custom_func", func(t *testing.T) {
		customFuncs := map[string]convert.ConvertFunc{
			"title": func(val interface{}, props mapping.ExtProperties) (interface{}, error) {
				return val, nil
			},
		}
		got, err := LuceneToDSL(
			`title:hello`,
			WithMappingData(mappingJSON),
			WithCustomConvertFunc(customFuncs),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestLuceneToDSL_FieldTypesWithMapping(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"keyword", `status:active`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
		{"text_single", `title:hello`, mustDSL(`{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}`), false},
		{"text_phrase", `title:"hello world"`, mustDSL(`{"match_phrase":{"title":{"boost":1,"query":"hello world"}}}`), false},
		{"integer", `count:100`, mustDSL(`{"term":{"count":{"boost":1,"value":100}}}`), false},
		{"integer_range", `count:[10 TO 100]`, mustDSL(`{"range":{"count":{"boost":1,"gte":10,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"float", `price:3.14`, mustDSL(`{"term":{"price":{"boost":1,"value":3.140000104904175}}}`), false},
		{"boolean", `is_active:true`, mustDSL(`{"term":{"is_active":{"boost":1,"value":true}}}`), false},
		{"date", `created_at:2021-01-01`, mustDSL(`{"range":{"created_at":{"boost":1,"format":"epoch_millis","gte":1609459200000,"lte":1640995199999,"relation":"INTERSECTS"}}}`), false},
		{"date_range", `created_at:[2021-01-01 TO 2021-12-31]`, mustDSL(`{"range":{"created_at":{"boost":1,"format":"epoch_millis","gte":1609459200000,"lte":1640908800000,"relation":"INTERSECTS"}}}`), false},
		{"ip", `ip_address:192.168.1.1`, mustDSL(`{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}}`), false},
		{"ip_range", `ip_address:[192.168.0.0 TO 192.168.255.255]`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.0","lte":"192.168.255.255","relation":"INTERSECTS"}}}`), false},
		{"byte", `level:5`, mustDSL(`{"term":{"level":{"boost":1,"value":5}}}`), false},
		{"half_float", `weight:1.5`, mustDSL(`{"term":{"weight":{"boost":1,"value":1.5}}}`), false},
		{"wildcard_type", `uuid:*abc*`, mustDSL(`{"wildcard":{"uuid":{"boost":1,"rewrite":"constant_score","value":"*abc*"}}}`), false},
		{"prefix", `status:act*`, mustDSL(`{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}`), false},
		{"wildcard", `status:act*ve`, mustDSL(`{"wildcard":{"status":{"boost":1,"rewrite":"constant_score","value":"act*ve"}}}`), false},
		{"fuzzy", `status:active~2`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
		{"regexp", `status:/act.*/`, mustDSL(`{"regexp":{"status":{"flags":"ALL","max_determinized_states":10000,"rewrite":"constant_score","value":"act.*"}}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_BooleanOperations(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"simple_or", `status:active OR status:pending`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"simple_not", `NOT status:inactive`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"term":{"status":{"boost":1,"value":"inactive"}}}}}`), false},
		{"multiple_or", `status:active OR status:pending OR status:inactive`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}},{"term":{"status":{"boost":1,"value":"inactive"}}}]}}`), false},
		{"not_with_range", `NOT count:<100`, mustDSL(`{"range":{"count":{"boost":1,"gte":100,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},
		{"not_with_prefix", `NOT status:act*`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}}}`), false},
		{"not_with_exists", `NOT _exists_:status`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"exists":{"field":"status"}}}}`), false},
		{"or_with_exists", `_exists_:status OR _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},
		{"or_with_text", `title:hello OR title:world`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"match":{"title":{"boost":1,"max_expansions":50,"query":"world"}}}]}}`), false},
		{"or_with_ip", `ip_address:192.168.1.1 OR ip_address:192.168.1.2`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}},{"term":{"ip_address":{"boost":1,"value":"192.168.1.2"}}}]}}`), false},
		{"complex_or", `status:active OR status:pending OR status:inactive OR status:deleted`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}},{"term":{"status":{"boost":1,"value":"inactive"}}},{"term":{"status":{"boost":1,"value":"deleted"}}}]}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_SpecialQueries(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"exists", `_exists_:status`, mustDSL(`{"exists":{"field":"status"}}`), false},
		{"exists_with_not", `NOT _exists_:status`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"exists":{"field":"status"}}}}`), false},
		{"exists_with_or", `_exists_:status OR _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},
		{"exists_with_and", `_exists_:status AND _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},
		{"exists_with_or_with_same_field", `_exists_:status OR _exists_:status`, mustDSL(`{"exists":{"field":"status"}}`), false},
		{"exists_with_and_with_same_field", `_exists_:status AND _exists_:status`, mustDSL(`{"exists":{"field":"status"}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_BoostParameters(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"boost_term", `status:active^2.0`, mustDSL(`{"term":{"status":{"boost":2,"value":"active"}}}`), false},
		{"boost_range", `count:[10 TO 100]^1.5`, mustDSL(`{"range":{"count":{"boost":1.5,"gte":10,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"boost_match_phrase", `title:"hello world"^1.5`, mustDSL(`{"match_phrase":{"title":{"boost":1.5,"query":"hello world"}}}`), false},
		{"boost_prefix", `status:act*^1.2`, mustDSL(`{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}`), false},
		{"boost_wildcard", `status:act*ve^1.3`, mustDSL(`{"wildcard":{"status":{"boost":1.3,"rewrite":"constant_score","value":"act*ve"}}}`), false},
		{"fuzzy", `status:active~2`, mustDSL(`{"term":{"status":{"boost":1,"value":"active"}}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_RangeQueries(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"closed_range", `count:[10 TO 100]`, mustDSL(`{"range":{"count":{"boost":1,"gte":10,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"open_range", `count:{10 TO 100}`, mustDSL(`{"range":{"count":{"boost":1,"gt":10,"lt":100,"relation":"INTERSECTS"}}}`), false},
		{"gt", `count:>10`, mustDSL(`{"range":{"count":{"boost":1,"gt":10,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},
		{"gte", `count:>=10`, mustDSL(`{"range":{"count":{"boost":1,"gte":10,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},
		{"lt", `count:<100`, mustDSL(`{"range":{"count":{"boost":1,"gt":-2147483648,"lt":100,"relation":"INTERSECTS"}}}`), false},
		{"lte", `count:<=100`, mustDSL(`{"range":{"count":{"boost":1,"gt":-2147483648,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"boost_range", `count:[10 TO 100]^1.5`, mustDSL(`{"range":{"count":{"boost":1.5,"gte":10,"lte":100,"relation":"INTERSECTS"}}}`), false},
		{"float_range", `price:[10.5 TO 100.5]`, mustDSL(`{"range":{"price":{"boost":1,"gte":10.5,"lte":100.5,"relation":"INTERSECTS"}}}`), false},
		{"date_range", `created_at:[2021-01-01 TO 2021-12-31]`, mustDSL(`{"range":{"created_at":{"boost":1,"format":"epoch_millis","gte":1609459200000,"lte":1640908800000,"relation":"INTERSECTS"}}}`), false},
		{"ip_range", `ip_address:[192.168.0.0 TO 192.168.255.255]`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.0","lte":"192.168.255.255","relation":"INTERSECTS"}}}`), false},
		{"invalid_range", `count:[100 TO 10]`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_ComplexQueries(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"or_with_exists", `_exists_:status OR _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},
		{"or_with_text", `title:hello OR title:world`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"match":{"title":{"boost":1,"max_expansions":50,"query":"world"}}}]}}`), false},
		{"or_with_ip", `ip_address:192.168.1.1 OR ip_address:192.168.1.2`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}},{"term":{"ip_address":{"boost":1,"value":"192.168.1.2"}}}]}}`), false},
		{"or_with_prefix", `status:act* OR status:pend*`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}},{"prefix":{"status":{"rewrite":"constant_score","value":"pend"}}}]}}`), false},
		{"complex_or", `status:active OR status:pending OR status:inactive OR status:deleted`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}},{"term":{"status":{"boost":1,"value":"inactive"}}},{"term":{"status":{"boost":1,"value":"deleted"}}}]}}`), false},
		{"not_with_exists", `NOT _exists_:status`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"exists":{"field":"status"}}}}`), false},
		{"not_with_range", `NOT count:<100`, mustDSL(`{"range":{"count":{"boost":1,"gte":100,"lt":2147483647,"relation":"INTERSECTS"}}}`), false},
		{"not_with_prefix", `NOT status:act*`, mustDSL(`{"bool":{"minimum_should_match":0,"must_not":{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}}}}`), false},
		{"multiple_exists_or", `_exists_:status OR _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_IPFieldBehavior(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"exact_ip", `ip_address:192.168.1.1`, mustDSL(`{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}}`), false},
		{"cidr", `ip_address:192.168.0.0/24`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.1","lte":"192.168.0.254","relation":"INTERSECTS"}}}`), false},
		{"ip_range", `ip_address:[192.168.0.0 TO 192.168.255.255]`, mustDSL(`{"range":{"ip_address":{"boost":1,"gte":"192.168.0.0","lte":"192.168.255.255","relation":"INTERSECTS"}}}`), false},
		{"ip_or", `ip_address:192.168.1.1 OR ip_address:192.168.1.2`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}},{"term":{"ip_address":{"boost":1,"value":"192.168.1.2"}}}]}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_TextFieldBehavior(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		{"single_term", `title:hello`, mustDSL(`{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}`), false},
		{"phrase", `title:"hello world"`, mustDSL(`{"match_phrase":{"title":{"boost":1,"query":"hello world"}}}`), false},
		{"with_boost", `title:hello^2.0`, mustDSL(`{"match":{"title":{"boost":2,"max_expansions":50,"query":"hello"}}}`), false},
		{"phrase_with_boost", `title:"hello world"^1.5`, mustDSL(`{"match_phrase":{"title":{"boost":1.5,"query":"hello world"}}}`), false},
		{"text_with_boolean", `title:hello AND status:active`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"term":{"status":{"boost":1,"value":"active"}}}]}}`), false},
		{"text_with_range", `title:hello AND count:>100`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}]}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}

func TestLuceneToDSL_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		opts    []Option
		wantErr bool
	}{
		{"empty_query", ``, []Option{WithMappingData(mappingJSON)}, true},
		{"invalid_syntax", `status: AND count:`, []Option{WithMappingData(mappingJSON)}, true},
		{"invalid_mapping", `status:active`, []Option{WithMappingData([]byte(`{invalid json}`))}, true},
		{"field_not_in_mapping", `unknown_field:value`, []Option{WithMappingData(mappingJSON)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LuceneToDSL(tt.query, tt.opts...)
			assert.Error(t, err)
		})
	}
}

func BenchmarkLuceneToDSL_WithMapping(b *testing.B) {
	queries := []string{
		`status:active`,
		`status:active AND count:>100`,
		`(status:active OR status:pending) AND count:>100`,
		`title:"hello world"`,
		`_exists_:status`,
	}

	for _, query := range queries {
		b.Run(query, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = LuceneToDSL(query, WithMappingData(mappingJSON))
			}
		})
	}
}

func BenchmarkLuceneToDSL_WithoutMapping(b *testing.B) {
	queries := []string{
		`status:active`,
		`status:active AND count:>100`,
		`(status:active OR status:pending) AND count:>100`,
	}

	for _, query := range queries {
		b.Run(query, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = LuceneToDSL(query)
			}
		})
	}
}

func TestLuceneToDSL_FilterContext(t *testing.T) {
	t.Run("filter_context_single_field", func(t *testing.T) {
		got, err := LuceneToDSL(
			`status:active AND count:>100`,
			WithMappingData(mappingJSON),
			WithFilterContext([]string{"status"}),
		)
		assert.NoError(t, err)
		assertDSLEqual(t, mustDSL(`{"bool":{"filter":{"term":{"status":{"boost":1,"value":"active"}}},"minimum_should_match":0,"must":{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}}}`), got)
	})

	t.Run("filter_context_multiple_fields", func(t *testing.T) {
		got, err := LuceneToDSL(
			`status:active AND count:>100 AND title:hello`,
			WithMappingData(mappingJSON),
			WithFilterContext([]string{"status", "count"}),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("filter_context_no_match", func(t *testing.T) {
		got, err := LuceneToDSL(
			`status:active AND count:>100`,
			WithMappingData(mappingJSON),
			WithFilterContext([]string{"unknown"}),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("filter_context_empty_patterns", func(t *testing.T) {
		got, err := LuceneToDSL(
			`status:active AND count:>100`,
			WithMappingData(mappingJSON),
			WithFilterContext([]string{}),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestLuceneToDSL_SubQueryCombinations(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    dsl.DSL
		wantErr bool
	}{
		// ========== Same field combinations ==========
		{"same_field_or", `status:active OR status:pending`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"same_field_and", `status:active AND status:pending`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"same_field_and_not", `status:active AND NOT status:pending`, mustDSL(`{"bool":{"minimum_should_match":0,"must":{"term":{"status":{"boost":1,"value":"active"}}},"must_not":{"term":{"status":{"boost":1,"value":"pending"}}}}}`), false},
		{"same_field_or_not", `status:active OR NOT status:pending`, mustDSL(`{"bool":{"minimum_should_match":1,"must_not":{"term":{"status":{"boost":1,"value":"pending"}}},"should":{"term":{"status":{"boost":1,"value":"active"}}}}}`), false},

		// ========== Different field combinations ==========
		{"diff_field_or", `status:active OR count:>100`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}]}}`), false},
		{"diff_field_and", `status:active AND count:>100`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}]}}`), false},
		{"diff_field_and_not", `status:active AND NOT count:>100`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":-2147483648,"lte":100,"relation":"INTERSECTS"}}}]}}`), false},
		{"diff_field_or_not", `status:active OR NOT count:>100`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":-2147483648,"lte":100,"relation":"INTERSECTS"}}}]}}`), false},

		// ========== Text + keyword combinations ==========
		{"text_and_keyword", `title:hello AND status:active`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"term":{"status":{"boost":1,"value":"active"}}}]}}`), false},
		{"text_or_keyword", `title:hello OR status:active`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"term":{"status":{"boost":1,"value":"active"}}}]}}`), false},
		{"text_and_not_keyword", `title:hello AND NOT status:active`, mustDSL(`{"bool":{"minimum_should_match":0,"must":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},"must_not":{"term":{"status":{"boost":1,"value":"active"}}}}}`), false},
		{"text_or_not_keyword", `title:hello OR NOT status:active`, mustDSL(`{"bool":{"minimum_should_match":1,"must_not":{"term":{"status":{"boost":1,"value":"active"}}},"should":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}}}`), false},

		// ========== Range + range combinations ==========
		{"range_or_range", `count:[10 TO 50] OR count:[60 TO 100]`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"range":{"count":{"boost":1,"gte":10,"lte":50,"relation":"INTERSECTS"}}},{"range":{"count":{"boost":1,"gte":60,"lte":100,"relation":"INTERSECTS"}}}]}}`), false},
		{"range_and_range", `count:[10 TO 50] AND count:[60 TO 100]`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"range":{"count":{"boost":1,"gte":10,"lte":50,"relation":"INTERSECTS"}}},{"range":{"count":{"boost":1,"gte":60,"lte":100,"relation":"INTERSECTS"}}}]}}`), false},

		// ========== Exists combinations ==========
		{"exists_and_exists", `_exists_:status AND _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},
		{"exists_or_exists", `_exists_:status OR _exists_:title`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"exists":{"field":"status"}},{"exists":{"field":"title"}}]}}`), false},

		// ========== Prefix + term combinations ==========
		{"prefix_or_term", `status:act* OR status:pending`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"prefix":{"status":{"rewrite":"constant_score","value":"act"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},

		// ========== Three-way combinations ==========
		{"three_way_and", `status:active AND count:>100 AND title:hello`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}},{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}]}}`), false},
		{"three_way_or", `status:active OR count:>100 OR title:hello`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}},{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}]}}`), false},

		// ========== Parenthesized combinations ==========
		{"paren_or_and", `(status:active OR status:pending) AND count:>100`, mustDSL(`{"bool":{"minimum_should_match":1,"must":{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}},"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"paren_and_or", `(status:active AND count:>100) OR title:hello`, mustDSL(`{"bool":{"minimum_should_match":1,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}],"should":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}}}`), false},

		// ========== AND NOT with different types ==========
		{"and_not_text", `status:active AND NOT title:hello`, mustDSL(`{"bool":{"minimum_should_match":0,"must":{"term":{"status":{"boost":1,"value":"active"}}},"must_not":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}}}`), false},
		{"and_not_range", `count:>100 AND NOT status:inactive`, mustDSL(`{"bool":{"minimum_should_match":0,"must":{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}},"must_not":{"term":{"status":{"boost":1,"value":"inactive"}}}}}`), false},

		// ========== OR NOT with different types ==========
		{"or_not_text", `status:active OR NOT title:hello`, mustDSL(`{"bool":{"minimum_should_match":1,"must_not":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},"should":{"term":{"status":{"boost":1,"value":"active"}}}}}`), false},
		{"or_not_range", `count:>100 OR NOT status:inactive`, mustDSL(`{"bool":{"minimum_should_match":1,"must_not":{"term":{"status":{"boost":1,"value":"inactive"}}},"should":{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}}}`), false},

		// ========== Complex nested combinations ==========
		{"nested_or_and", `(status:active OR status:pending) AND (count:>100 AND count:<200)`, mustDSL(`{"bool":{"minimum_should_match":1,"must":{"range":{"count":{"boost":1,"gt":100,"lt":200,"relation":"INTERSECTS"}}},"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"term":{"status":{"boost":1,"value":"pending"}}}]}}`), false},
		{"nested_and_or", `(status:active AND count:>100) OR (title:hello AND title:world)`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}]}},{"bool":{"minimum_should_match":0,"must":[{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}},{"match":{"title":{"boost":1,"max_expansions":50,"query":"world"}}}]}}]}}`), false},

		// ========== IP field combinations ==========
		{"ip_or_ip", `ip_address:192.168.1.1 OR ip_address:192.168.1.2`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}},{"term":{"ip_address":{"boost":1,"value":"192.168.1.2"}}}]}}`), false},
		{"ip_and_ip", `ip_address:192.168.1.1 AND ip_address:192.168.1.2`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"ip_address":{"boost":1,"value":"192.168.1.1"}}},{"term":{"ip_address":{"boost":1,"value":"192.168.1.2"}}}]}}`), false},

		// ========== Boolean + integer combinations ==========
		{"bool_and_integer", `is_active:true AND count:>100`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"is_active":{"boost":1,"value":true}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}]}}`), false},
		{"bool_or_integer", `is_active:true OR count:>100`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"is_active":{"boost":1,"value":true}}},{"range":{"count":{"boost":1,"gt":100,"lt":2147483647,"relation":"INTERSECTS"}}}]}}`), false},

		// ========== Date + keyword combinations ==========
		{"date_and_keyword", `created_at:[2021-01-01 TO 2021-12-31] AND status:active`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"range":{"created_at":{"boost":1,"format":"epoch_millis","gte":1609459200000,"lte":1640908800000,"relation":"INTERSECTS"}}},{"term":{"status":{"boost":1,"value":"active"}}}]}}`), false},

		// ========== Multiple NOT combinations ==========
		{"multiple_and_not", `status:active AND NOT count:<100 AND NOT title:hello`, mustDSL(`{"bool":{"minimum_should_match":0,"must":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gte":100,"lt":2147483647,"relation":"INTERSECTS"}}}],"must_not":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}}}`), false},
		{"multiple_or_not", `status:active OR NOT count:<100 OR NOT title:hello`, mustDSL(`{"bool":{"minimum_should_match":1,"should":[{"term":{"status":{"boost":1,"value":"active"}}},{"range":{"count":{"boost":1,"gte":100,"lt":2147483647,"relation":"INTERSECTS"}}},{"bool":{"minimum_should_match":0,"must_not":{"match":{"title":{"boost":1,"max_expansions":50,"query":"hello"}}}}}]}}`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LuceneToDSL(tt.query, WithMappingData(mappingJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertDSLEqual(t, tt.want, got)
			}
		})
	}
}
