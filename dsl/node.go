package dsl

import "github.com/zhuliquan/lucene-to-dsl/mapping"

// define ast node of dsl
type AstNode interface {
	AstType() AstType
	DslType() DslType
	NodeKey() string

	// union_join / intersect / inverse nodes with same NodeKey (get by NodeKey() func)
	UnionJoin(AstNode) (AstNode, error)
	InterSect(AstNode) (AstNode, error)
	Inverse() (AstNode, error)
	ToDSL() DSL
}

// boost node interface
type BoostNode interface {
	getBoost() float64
	setBoost(float64)
}

func WithBoost(boost float64) func(AstNode) {
	return func(n AstNode) {
		if b, ok := n.(BoostNode); ok {
			b.setBoost(boost)
		}
	}
}

// boost node impl
type boostNode struct {
	boost float64
}

func (b *boostNode) getBoost() float64 {
	return b.boost
}

func (b *boostNode) setBoost(boost float64) {
	b.boost = boost
}

// analyzer node interface
type AnalyzerNode interface {
	setAnalyzer(analyzer string)
	getAnaLyzer() string
}

func WithAnalyzer(analyzer string) func(AstNode) {
	return func(n AstNode) {
		if a, ok := n.(AnalyzerNode); ok {
			a.setAnalyzer(analyzer)
		}
	}
}

// analyzer node impl
type analyzerNode struct {
	analyzer string
}

func (a *analyzerNode) getAnaLyzer() string {
	return a.analyzer
}

func (a *analyzerNode) setAnalyzer(analyzer string) {
	a.analyzer = analyzer
}

// rewrite node interface
type RewriteNode interface {
	setRewrite(string)
}

func WithRewrite(rewrite string) func(AstNode) {
	return func(n AstNode) {
		if r, ok := n.(RewriteNode); ok {
			r.setRewrite(rewrite)
		}
	}
}

// rewrite node impl
type rewriteNode struct {
	rewrite string
}

func (r *rewriteNode) setRewrite(rewrite string) {
	r.rewrite = rewrite
}

// indicate whether does dsl query use filter context
type FilterCtxNode interface {
	filterCtx() bool
}

// indicate whether is node array data type
type ArrayTypNode interface {
	isArray() bool
}

type opNode struct {
	filter bool
}

func NewOpNode(filter bool) *opNode {
	return &opNode{filter: filter}
}

func (n *opNode) filterCtx() bool {
	return n.filter
}

func (n *opNode) AstType() AstType {
	return OP_NODE_TYPE
}

// leaf node
type lfNode struct {
	filter bool // whether node is filter
	isList bool // whether node is list type
}

func NewLfNode(opts ...func(*lfNode)) *lfNode {
	var n = &lfNode{
		filter: false,
		isList: true,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func UseFilterCtx(filter bool) func(*lfNode) {
	return func(lf *lfNode) {
		lf.filter = filter
	}
}

func IsArrayTyp(isList bool) func(*lfNode) {
	return func(lf *lfNode) {
		lf.isList = isList
	}
}

func (n *lfNode) filterCtx() bool {
	return n.filter
}

func (n *lfNode) isArray() bool {
	return n.isList
}

func (n *lfNode) AstType() AstType {
	return LEAF_NODE_TYPE
}

type fieldNode struct {
	lfNode
	field string
}

type valueNode struct {
	value LeafValue
	mType mapping.FieldType
}

func NewValueNode(value LeafValue, mType mapping.FieldType) *valueNode {
	return &valueNode{value: value, mType: mType}
}

func (v *valueNode) toPrintValue() interface{} {
	return leafValueToPrintValue(v.value, v.mType)
}

func NewFieldNode(lfNode *lfNode, field string) *fieldNode {
	return &fieldNode{
		lfNode: *lfNode,
		field:  field,
	}
}

func (n *fieldNode) NodeKey() string {
	return n.field
}

// Key value node
type kvNode struct {
	fieldNode
	valueNode
}

func NewKVNode(fieldNode *fieldNode, value *valueNode) *kvNode {
	return &kvNode{
		fieldNode: *fieldNode,
		valueNode: *value,
	}
}

type rgNode struct {
	fieldNode
	mType   mapping.FieldType
	rValue  LeafValue
	lValue  LeafValue
	rCmpSym CompareType
	lCmpSym CompareType
}

func NewRgNode(fieldNode *fieldNode, mType mapping.FieldType, lValue, rValue LeafValue, lCmpSym, rCmpSym CompareType) *rgNode {
	return &rgNode{
		fieldNode: *fieldNode,
		mType:     mType,
		lValue:    lValue,
		rValue:    rValue,
		lCmpSym:   lCmpSym,
		rCmpSym:   rCmpSym,
	}
}

func (n *rgNode) NodeKey() string {
	return n.field
}
