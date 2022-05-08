package lucene_to_dsl

import (
	"fmt"
	"sync"

	"github.com/zhuliquan/lucene-to-dsl/convert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	lucene "github.com/zhuliquan/lucene_parser"
)

var (
	mappingPath string
	customFuncs map[string]func(interface{}, map[string]interface{}) (interface{}, error)
	onceInit    sync.Once
)

func LoadMappingPath(path string) {
	mappingPath = path
}

func LoadCustomFuncs(funcs map[string]func(interface{}, map[string]interface{}) (interface{}, error)) {
	customFuncs = funcs
}

func LuceneToDSL(luceneQuery string) (dsl.DSL, error) {
	onceInit.Do(
		func() {
			if pm, err := mapping.Init(mappingPath, customFuncs); err != nil {
				panic(err)
			} else {
				convert.Init(pm)
			}
		},
	)

	var err error
	var qry *lucene.Lucene
	var nod dsl.DSLNode
	defer func() {
		if r := recover(); r != nil {
			nod = &dsl.EmptyNode{}
			err = fmt.Errorf("failed to lucene to dsl, err: %v", r)
		}
	}()

	if qry, err = lucene.ParseLucene(luceneQuery); err != nil {
		return nil, err
	}

	if nod, err = convert.LuceneToDSLNode(qry); err != nil {
		return nil, err
	}

	return nod.ToDSL(), nil
}
