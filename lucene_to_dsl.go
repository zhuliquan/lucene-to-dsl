package lucene_to_dsl

import (
	"fmt"

	mapping "github.com/zhuliquan/es-mapping"
	"github.com/zhuliquan/lucene-to-dsl/convert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	lucene "github.com/zhuliquan/lucene_parser"
)

type Config struct {
	mappingData []byte
	customFuncs map[string]convert.ConvertFunc
}

type Option func(*Config)

func WithMappingData(data []byte) Option {
	return func(o *Config) {
		o.mappingData = data
	}
}

func WithCustomConvertFunc(funcs map[string]convert.ConvertFunc) Option {
	return func(o *Config) {
		o.customFuncs = funcs
	}
}

func LuceneToDSL(
	query string,
	opts ...Option,
) (dsl.DSL, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	pm, err := mapping.LoadMappingData(cfg.mappingData)
	if err != nil {
		return nil, fmt.Errorf("failed to load mapping data, err: %v", err)
	}

	var cvt = convert.NewConverter(pm, cfg.customFuncs)
	var qry *lucene.Lucene
	var nod dsl.AstNode
	defer func() {
		if r := recover(); r != nil {
			nod = &dsl.EmptyNode{}
			err = fmt.Errorf("failed to lucene to dsl, err: %v", r)
		}
	}()

	if qry, err = lucene.ParseLucene(query); err != nil {
		return nil, err
	}

	if nod, err = cvt.LuceneToAstNode(qry); err != nil {
		return nil, err
	}

	return nod.ToDSL(), nil
}
