# lucene-to-dsl
## Introduction:
This package can parse lucene query and convert to dsl used by ES (ElasticSearch), this package is pure go package.
## Features
- 1、This package can convert lucene query to dsl which is used by ES.
- 2、This package can compact many leaf nodes to fewer leaf nodes (i.g. `x:>1 AND x:<10` => `{"range": {"x": {"gt": 1, "lt": 10}}}` instead of `{"bool": {"must": [{"range": {"x": {"gt": 1}}}, {"range": {"x": {"lt": 10}}}]}}`). compact dsl will be serached more faster than uncompact dsl. for example two range dsl compact to single range dsl, which can reduce a range query and two bitsect intersect.
- 3、This package can filter some wrong lucene query (i.g. `x:>1 AND x:<-1` is wrong lucene query).
- 4、This package can process wildcard field (i.g. `_exist_:fo\?bar\*`, `foo\?bar\*:bar`)

## Limitations
- 1、only support lucene query with **field name**, instead of query without **field name** (i.e. this project can't parse query like `foo OR bar`, `foo AND bar`).
- 2、don't support prefix operator `'+'` / `'-'`, for instance `+fo0 -bar`.
- 3、should give [index mapping of field](https://www.elastic.co/guide/en/elasticsearch/reference/7.15/mapping.html).
- 4、 will ignore `boost` parameter in field mapping which using in index time boosting.
- 5、 don't support alias field type (a kind of filed mapping type).

### mapping of field
In order to convert more accurately, you need the configuration of a given field, such as mapping of field. 
