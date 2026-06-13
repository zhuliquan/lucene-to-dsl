# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/lang/zh-CN/).

## [v0.1.1] - 2026-06-14

### Added

- 新增 `WithFilterContext` 选项，支持指定字段使用 filter context（不评分、可缓存），提升查询性能
- 新增 `NewConverterWithFilter` 构造函数，支持传入 filter 字段模式

### Fixed

- 修复 `_exists_:field` 查询生成错误的 `term` DSL，现正确生成 `{"exists":{"field":"xxx"}}`
- 修复 `ExistsNode.UnionJoin` / `InterSect` 在两个不同字段 ExistsNode 交互时的无限递归 bug
- 修复 `TermNode`、`RangeNode`、`PrefixNode`、`WildCardNode` 的 `UnionJoin` / `InterSect` 缺少 `NodeKey()` 一致性检查，导致不同字段节点直接交叉运算 panic

### Changed

- text 字段单次 term 查询由 `query_string` 改为 `match`，更符合 ES 全文搜索语义
- keyword / wildcard 字段的通配符模式（如 `act*`、`act*ve`）由嵌入 term 值改为生成专用 `prefix` / `wildcard` DSL 节点

## [v0.1.0] - 2026-06-13

### Features

- 支持 Lucene 查询语法解析并转换为 Elasticsearch DSL
- 支持 14 种 DSL 节点类型：term、match、match_phrase、range、bool、prefix、wildcard、fuzzy、regexp、exists、ids、query_string、match_all、empty
- 支持 30+ 种 ES 字段类型：keyword、text、integer、float、boolean、date、ip、version、wildcard 等
- 支持查询压缩/优化：同字段 range 查询合并、bool 代数简化等
- 支持无 mapping 自动类型推断
- 支持自定义值转换函数
- 支持 AND / OR / NOT 布尔逻辑
- 支持 boost 参数
- 支持正则、模糊、前缀、通配符查询
- 支持 `_exists_`、`_id` 特殊字段查询
- 支持 `*:*` match_all 查询
