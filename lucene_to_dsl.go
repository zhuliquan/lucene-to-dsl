package lucene_to_dsl

import (
	"fmt"
	"sync"

	"github.com/zhuliquan/lucene-to-dsl/convert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	mapping "github.com/zhuliquan/es-mapping"
	lucene "github.com/zhuliquan/lucene_parser"
)

var (
	mappingPath string
	customFuncs map[string]convert.ConvertFunc
	onceInit    sync.Once
	converter   convert.Converter
)

func LoadMappingPath(path string) {
	mappingPath = path
}

func LoadCustomFuncs(funcs map[string]convert.ConvertFunc) {
	customFuncs = funcs
}

func LuceneToDSL(luceneQuery string) (dsl.DSL, error) {
	onceInit.Do(
		func() {
			if pm, err := mapping.LoadMappingFile(mappingPath); err != nil {
				panic(err)
			} else {
				converter = convert.NewConverter(pm, customFuncs)
			}
		},
	)

	var err error
	var qry *lucene.Lucene
	var nod dsl.AstNode
	defer func() {
		if r := recover(); r != nil {
			nod = &dsl.EmptyNode{}
			err = fmt.Errorf("failed to lucene to dsl, err: %v", r)
		}
	}()

	if qry, err = lucene.ParseLucene(luceneQuery); err != nil {
		return nil, err
	}

	if nod, err = converter.LuceneToAstNode(qry); err != nil {
		return nil, err
	}

	return nod.ToDSL(), nil
}
