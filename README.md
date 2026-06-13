# lucene-to-dsl

## Introduction

This package can parse lucene query and convert to dsl used by ES (ElasticSearch), this package is pure go package.

## Installation

```bash
go get github.com/zhuliquan/lucene-to-dsl
```

## Quick Start

```go
package main

import (
    "fmt"
    lucenedsl "github.com/zhuliquan/lucene-to-dsl"
)

func main() {
    // Load ES mapping file
    lucenedsl.LoadMappingPath("/path/to/mapping.json")

    // Convert lucene query to DSL
    dsl, err := lucenedsl.LuceneToDSL(`foo:bar AND baz:[1 TO 10]`)
    if err != nil {
        panic(err)
    }
    fmt.Println(dsl.String())
}
```

## Features

- 1、This package can convert lucene query to dsl which is used by ES.
- 2、This package can compact many leaf nodes to fewer leaf nodes (i.g. `x:>1 AND x:<10` => `{"range": {"x": {"gt": 1, "lt": 10}}}` instead of `{"bool": {"must": [{"range": {"x": {"gt": 1}}}, {"range": {"x": {"lt": 10}}}]}}`). compact dsl will be serached more faster than uncompact dsl. for example two range dsl compact to single range dsl, which can reduce a range query and two bitsect intersect.
- 3、This package can filter some wrong lucene query (i.g. `x:>1 AND x:<-1` is wrong lucene query).
- 4、This package can process wildcard field (i.e. `_exist_:fo\?bar\*`, `foo\?bar\*:bar`)

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
// LoadMappingPath loads ES mapping file path
func LoadMappingPath(path string)

// LoadCustomFuncs loads custom field value conversion functions
func LoadCustomFuncs(funcs map[string]convert.ConvertFunc)

// LuceneToDSL converts lucene query string to ES DSL
func LuceneToDSL(luceneQuery string) (dsl.DSL, error)
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
- 3、should give [index mapping of field](https://www.elastic.co/guide/en/elasticsearch/reference/7.15/mapping.html).
- 4、will ignore `boost` parameter in field mapping which using in index time boosting.
- 5、don't support alias field type (a kind of filed mapping type).

### mapping of field

In order to convert more accurately, you need the configuration of a given field, such as mapping of field.

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
