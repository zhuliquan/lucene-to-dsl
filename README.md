# lucene-to-dsl
## Introduction:
This package can parse lucene query and conver to dsl used by ES (ElasticSearch), this package is pure go package.
## Features
- 1、support phrase term query, for instance `x:"foo bar"`.
- 2、support regexp term query, for instance `x:/\d+\\.?\d+/`.
- 3、support bool operator （i.e. `AND`, `OR`, `NOT`, `&&`, `||`, `!`） join sub query, for instance `x:1 AND y:2`, `x:1 || y:2`.
- 4、support bound range query,  for instance `x:[1 TO 2]`, `x:[1 TO 2}`.
- 5、support side range query, for instance `x:>1` , `x:>=1` , `x:<1` , `x:<=1`.
- 6、support boost term, for instance `x:1^2` , `x:"dsada 8908"^3`
- 7、support fuzzy query, for instance `x:for~2` , `x:"foo bar"~2`

## Limitations
- 1、only support lucene query with **field name**, instead of query without **field name** (i.e. this project can't parse query like `foo OR bar`, `foo AND bar`).
