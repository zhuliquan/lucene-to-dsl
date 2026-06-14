# lucene-to-dsl

## Introduction

This package can parse lucene query and convert it to more efficient DSL queries used by ES (ElasticSearch). It automatically optimizes the generated DSL, such as compacting multiple range conditions into a single range query, which significantly improves search performance.

## Installation

```bash
go get github.com/zhuliquan/lucene-to-dsl
```

## Quick Start

### Without Mapping (Auto Type Inference)

```go
package main

import (
    "fmt"
    luceneDsl "github.com/zhuliquan/lucene-to-dsl"
)

func main() {
    // Convert lucene query without providing mapping
    // The library will automatically infer field types based on values
    dsl, err := luceneDsl.LuceneToDSL(`status:active AND views:>100`)
    if err != nil {
        panic(err)
    }
    fmt.Println(dsl.String())
    // Output: {"bool":{"must":[{"term":{"status":{"value":"active","boost":1.0}}},{"range":{"views":{"gt":100,"boost":1.0,"relation":"INTERSECTS"}}}]}}
}
```

### With Mapping (Recommended for Production)

```go
package main

import (
    "fmt"
    "os"
    luceneDsl "github.com/zhuliquan/lucene-to-dsl"
)

func main() {
    // Load ES mapping file data
    mappingData, err := os.ReadFile("/path/to/mapping.json")
    if err != nil {
        panic(err)
    }

    // Convert lucene query to DSL
    dsl, err := luceneDsl.LuceneToDSL(
        `foo:bar AND baz:[1 TO 10]`,
        luceneDsl.WithMappingData(mappingData),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(dsl.String())
}
```

## Features

- 1、This package can convert lucene query to dsl which is used by ES.
- 2、**Query optimization** - This package can compact many leaf nodes to fewer leaf nodes (i.e. `x:>1 AND x:<10` => `{"range": {"x": {"gt": 1, "lt": 10}}}` instead of `{"bool": {"must": [{"range": {"x": {"gt": 1}}}, {"range": {"x": {"lt": 10}}}]}}`). This optimization reduces the number of range queries and bitset intersections, significantly improving search performance.
- 3、**Node simplification** - Automatically simplifies complex boolean structures by removing unnecessary wrapper nodes, reducing the complexity of generated DSL queries.
- 4、**Smart deduplication** - Merges duplicate or overlapping conditions to avoid redundant queries, further optimizing search performance.
- 5、**Wildcard field support** - Handles wildcard fields (i.e. `_exist_:fo\?bar\*`, `foo\?bar*:bar`).
- 6、**No mapping required** - This package supports automatic type inference. When no mapping is provided, it will infer field types based on values (e.g., integers, dates, IP addresses).
- 7、**Filter context optimization** - Supports using filter context for non-scoring queries, which can be cached by Elasticsearch for better performance.
- 8、**Intelligent NOT operation handling** - Optimizes NOT operations by reducing the number of must_not clauses, making negation queries more efficient.

## Auto Type Inference

When no mapping is provided, the library automatically infers field types based on the query values:

| Value Pattern | Inferred Type | DSL Generated |
|---------------|---------------|---------------|
| `true` / `false` | `boolean` | `term` |
| `123`, `-456` | `keyword` | `term` |
| `3.14`, `-2.5` | `keyword` | `term` |
| `2021-01-01`, `2021-01-01T12:00:00` | `date` | `range` |
| `192.168.1.1`, `2001:db8::1` | `ip` | `term` |
| `192.168.0.0/24` | `ip` | `range` (CIDR) |
| Other strings | `keyword` | `term` |

### Examples Without Mapping

```go
// Boolean field
dsl, _ := luceneDsl.LuceneToDSL(`active:true`)
// Output: {"term":{"active":{"value":true,"boost":1.0}}}

// Numeric field
dsl, _ := luceneDsl.LuceneToDSL(`count:100`)
// Output: {"term":{"count":{"value":100,"boost":1.0}}}

// Date field
dsl, _ := luceneDsl.LuceneToDSL(`created_at:2021-01-01`)
// Output: {"range":{"created_at":{"gte":"2021-01-01T00:00:00","lte":"2021-01-01T23:59:59.999999999","boost":1.0,"relation":"INTERSECTS"}}}

// IP field
dsl, _ := luceneDsl.LuceneToDSL(`ip_address:192.168.1.1`)
// Output: {"term":{"ip_address":{"value":"192.168.1.1","boost":1.0}}}
```

## Supported Lucene Query Syntax

| Syntax | ES DSL | Description |
|--------|--------|-------------|
| `field:value` | `term` | Exact match |
| `field:"phrase"` | `match_phrase` | Phrase match |
| `field:[v1 TO v2]` | `range` | Closed interval range |
| `field:{v1 TO v2}` | `range` | Open interval range |
| `field:>v` | `range` | Greater than |
| `field:>=v` | `range` | Greater than or equal |
| `field:<v` | `range` | Less than |
| `field:<=v` | `range` | Less than or equal |
| `field:regex` | `regexp` | Regular expression |
| `field:fuzzy~2` | `fuzzy` | Fuzzy match |
| `field:prefix*` | `prefix` | Prefix query |
| `field:wild*card` | `wildcard` | Wildcard query |
| `_exists_:field` | `exists` | Field exists |
| `*:*` | `match_all` | Match all |
| `_id:xxx` | `ids` | IDs query |
| `AND` / `&&` | `bool.must` | Logical AND |
| `OR` / `\|\|` | `bool.should` | Logical OR |
| `NOT` / `-` | `bool.must_not` | Logical NOT |
| `()` | recursive | Grouping |
| `^boost` | `boost` | Field boost |

## Supported ES Field Types

| Category | Types | DSL Generated |
|----------|-------|---------------|
| Numeric | boolean, byte, short, integer, long, unsigned_long, half_float, float, double, scaled_float | `term` / `range` |
| String | keyword, constant_keyword, wildcard | `term` |
| Text | text, match_only_text | `query_string` / `match_phrase` |
| Date | date, date_range, date_nanos | `range` (epoch_millis) |
| IP | ip, ip_range | `term` / `range` (CIDR) |
| Special | version | `term` / `range` |

## API Reference

### Functions

```go
// WithMappingData provides es mapping data as []byte for the converter
func WithMappingData(data []byte) func(*Config)

// WithCustomConvertFunc provides custom field value conversion functions
func WithCustomConvertFunc(funcs map[string]convert.ConvertFunc) func(*Config)

// WithFilterContext provides convert some pattern fields with filter mode query instead must bool query
func WithFilterContext(patterns []string) func(*Config)

// LuceneToDSL converts lucene query string to ES DSL
func LuceneToDSL(query string, opts ...func(*Config)) (dsl.DSL, error)
```

### DSL Type

```go
type DSL map[string]interface{}
```

## Examples

| Lucene Query | ES DSL Output |
|--------------|---------------|
| `foo:bar` | `{"term":{"foo":{"value":"bar","boost":1.0}}}` |
| `foo:>1 AND foo:<10` | `{"range":{"foo":{"gt":1,"lt":10}}}` |
| `foo:bar OR foo:baz` | `{"bool":{"should":[{"term":{"foo":"bar"}},{"term":{"foo":"baz"}}]}}` |
| `_exists_:foo` | `{"exists":{"field":"foo"}}` |
| `*:*` | `{"match_all":{}}` |
| `_id:abc` | `{"ids":{"values":["abc"]}}` |
| `foo:bar~2` | `{"fuzzy":{"foo":{"value":"bar","fuzziness":"2"}}}` |
| `foo:/regex/` | `{"regexp":{"foo":{"value":"regex"}}}` |
| `foo:bar*` | `{"prefix":{"foo":{"value":"bar"}}}` |
| `NOT foo:bar` | `{"bool":{"must_not":{"term":{"foo":"bar"}}}}` |

## Limitations

- 1、only support lucene query with **field name**, instead of query without **field name** (i.e. this project can't parse query like `foo OR bar`, `foo AND bar`).
- 2、don't support prefix operator `'+'` / `'-'`, for instance `+foo -bar`.
- 3、without mapping, type inference is based on value patterns (may not match actual field type in ES).
- 4、will ignore `boost` parameter in field mapping which using in index time boosting.
- 5、don't support alias field type (a kind of filed mapping type).

## Field Mapping Configuration

For more accurate conversion, you can provide the configuration of a given field, such as mapping of field. The mapping file defines the types of fields in your Elasticsearch index. This is essential for accurate query conversion because different field types generate different DSL queries.

**Note:** Mapping is optional. When not provided, the library will automatically infer field types based on values.

### Mapping File Format

The mapping file follows the Elasticsearch mapping format. Here's an example:

```json
{
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
}
```

### Usage Examples

#### 1. Load Mapping and Convert Query

```go
package main

import (
    "fmt"
    "os"
    luceneDsl "github.com/zhuliquan/lucene-to-dsl"
)

func main() {
    // Load mapping file
    mappingData, _ := os.ReadFile("/path/to/mapping.json")
    
    // Convert lucene query
    dsl, err := luceneDsl.LuceneToDSL(
        `status:active AND views:>100`,
        luceneDsl.WithMappingData(mappingData),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(dsl.String())
    // Output: {"bool":{"must":[{"term":{"status":{"value":"active","boost":1.0}}},{"range":{"views":{"gt":100,"boost":1.0,"relation":"INTERSECTS"}}}]}}
}
```

#### 2. How Field Types Affect Conversion

| Field Type | Lucene Query | Generated DSL |
|------------|--------------|---------------|
| keyword | `status:active` | `{"term":{"status":{"value":"active"}}}` |
| text | `title:hello` | `{"query_string":{"query":"hello","fields":["title"]}}` |
| text (phrase) | `title:"hello world"` | `{"match_phrase":{"title":{"query":"hello world"}}}` |
| integer | `views:100` | `{"term":{"views":{"value":100}}}` |
| integer (range) | `views:>100` | `{"range":{"views":{"gt":100}}}` |
| date | `created_at:2021-01-01` | `{"range":{"created_at":{"gte":"2021-01-01T00:00:00"}}}` |
| ip | `ip_address:192.168.1.1` | `{"term":{"ip_address":{"value":"192.168.1.1"}}}` |
| ip (CIDR) | `ip_address:192.168.0.0/24` | `{"range":{"ip_address":{"gte":"192.168.0.0","lte":"192.168.0.255"}}}` |

### Custom Value Conversion

You can also define custom conversion functions for specific fields:

```go
package main

import (
    "fmt"
    "os"
    "strings"
    luceneDsl "github.com/zhuliquan/lucene-to-dsl"
    "github.com/zhuliquan/lucene-to-dsl/convert"
    "github.com/zhuliquan/es-mapping"
)

func main() {
    // Load mapping file
    mappingData, _ := os.ReadFile("/path/to/mapping.json")
    
    // Define custom conversion functions
    customFuncs := map[string]convert.ConvertFunc{
        "title": func(val interface{}, props mapping.ExtProperties) (interface{}, error) {
            // Convert value to uppercase
            if str, ok := val.(string); ok {
                return strings.ToUpper(str), nil
            }
            return val, nil
        },
    }
    
    // Convert query with mapping data and custom functions
    dsl, err := luceneDsl.LuceneToDSL(
        `title:hello`,
        luceneDsl.WithMappingData(mappingData),
        luceneDsl.WithCustomConvertFunc(customFuncs),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(dsl.String())
}
```

### Using ExtProperties

The `ExtProperties` field in mapping allows you to define custom properties that can be used in conversion functions:

```json
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "ext_properties": {
          "only_lower": true
        }
      }
    }
  }
}
```

Then in your custom function:

```go
customFuncs := map[string]convert.ConvertFunc{
    "title": func(val interface{}, props mapping.ExtProperties) (interface{}, error) {
        if onlyLower, ok := props["only_lower"].(bool); ok && onlyLower {
            if str, ok := val.(string); ok {
                return strings.ToLower(str), nil
            }
        }
        return val, nil
    },
}

dsl, err := luceneDsl.LuceneToDSL(
    `title:hello`,
    luceneDsl.WithMappingData(mappingData),
    luceneDsl.WithCustomConvertFunc(customFuncs),
)
```

## Dependencies

| Package | Version | Description |
|---------|---------|-------------|
| `lucene_parser` | v0.5.1 | Lucene query syntax parser |
| `es-mapping` | v1.1.0 | ES mapping type system |
| `datemath_parser` | v0.0.9 | ES date math expression parser |
| `go_tools` | v0.0.1 | IP CIDR utility functions |
| `go-version` | v1.6.0 | Version number comparison |
| `float16` | v0.8.4 | half_float type support |
| `scaled_float` | v0.0.2 | scaled_float type support |

## Testing

```bash
go test ./...
```

## License

MIT
