package lucene_to_dsl

import (
	"fmt"

	mapping "github.com/zhuliquan/es-mapping"
	"github.com/zhuliquan/lucene-to-dsl/convert"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	lucene "github.com/zhuliquan/lucene_parser"
)

type Config struct {
	mappingData    []byte
	customFuncs    map[string]convert.ConvertFunc
	filterPatterns []string
}

type Option func(*Config)


// WithMappingData provides es mapping data as []byte for the converter
func WithMappingData(data []byte) Option {
	return func(o *Config) {
		o.mappingData = data
	}
}

// WithCustomConvertFunc provides custom field value conversion functions
func WithCustomConvertFunc(funcs map[string]convert.ConvertFunc) Option {
	return func(o *Config) {
		o.customFuncs = funcs
	}
}

// WithFilterContext provides convert some pattern field with filter mode query instead must bool query
func WithFilterContext(patterns []string) Option {
	return func(o *Config) {
		o.filterPatterns = patterns
	}
}

// LuceneToDSL converts lucene query string to ES DSL
func LuceneToDSL(
	query string,
	opts ...Option,
) (dsl.DSL, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	var pm *mapping.PropertyMapping
	var err error
	if len(cfg.mappingData) != 0 {
		pm, err = mapping.LoadMappingData(cfg.mappingData)
		if err != nil {
			return nil, fmt.Errorf("failed to load mapping data, err: %v", err)
		}
	}

	var cvt convert.Converter
	if len(cfg.filterPatterns) > 0 {
		cvt = convert.NewConverterWithFilter(pm, cfg.customFuncs, cfg.filterPatterns)
	} else {
		cvt = convert.NewConverter(pm, cfg.customFuncs)
	}
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
