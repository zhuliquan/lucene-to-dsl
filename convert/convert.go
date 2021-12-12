package convert

import (
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/lucene"
)

type Convert func(lucene.Query) (dsl.DSLNode, error)

func ToDSLNode(lucene lucene)