# lucene-to-dsl
## Introduction:
This package can parse lucene query and conver to dsl used by ES (ElasticSearch), this package is pure go package.
## Features
- 1、only support field name, example `x:789`.
- 2、support phrase term query, example `x:"foo bar"`.
- 3、support regexp term query, example `x:/\d+\.\d+/`.
- 4、support bool operator （i.e. `AND` / `OR` / `NOT` / `&&` / `||` / `!`） join sub query, example `x:1 AND y:2`  /  `x:1 || y:2`.
- 5、support bound range query, example `x:[1 TO 2]` / `x:[1 TO 2}`.
- 6、support side range query, example `x:>1` / `x:>=1` / `x:<1` / `x:<=1`.
- 7、support boost term, example `x:1^2` / `x:"dsada 8908"^3`
- 8、support fuzzy query, example `x:for~2` / `x:"foo bar"~2`
