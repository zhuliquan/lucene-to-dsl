package lucene_to_dsl

import (
	"github.com/zhuliquan/lucene-to-dsl/convert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	lucene "github.com/zhuliquan/lucene_parser"
)

func LuceneToDSL(luceneQuery string, path string, covFunc map[string]func(string) (interface{}, error)) (dsl.DSL, error) {
	var fm *mapping.Mapping
	var err error
	var qry *lucene.Lucene
	var node dsl.DSLNode

	if fm, err = mapping.LoadMapping(path); err != nil {
		return nil, err
	}
	if err = convert.InitConvert(fm, covFunc); err != nil {
		return nil, err
	}
	if qry, err = lucene.ParseLucene(luceneQuery); err != nil {
		return nil, err
	}

	if node, err = convert.LuceneToDSLNode(qry); err != nil {
		return nil, err
	}

	return node.ToDSL(), nil
}
