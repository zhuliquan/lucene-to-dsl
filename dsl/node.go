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

// expand node interface
type ExpandsNode interface {
	setMaxExpands(maxExpand int)
	getMaxExpands() int
}

func WithMaxExpands(maxExpand int) func(AstNode) {
	return func(n AstNode) {
		if e, ok := n.(ExpandsNode); ok {
			e.setMaxExpands(maxExpand)
		}
	}
}

type expandsNode struct {
	maxExpands int
}

func (n *expandsNode) setMaxExpand(maxExpands int) {
	n.maxExpands = maxExpands
}

func (n *expandsNode) getMaxExpands() int {
	return n.maxExpands
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
	getRewrite() string
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

func (r *rewriteNode) getRewrite() string {
	return r.rewrite
}

// slop node interface
type SlopNode interface {
	setSlop(slop int)
	getSlop() int
}

func WithSlop(slop int) func(AstNode) {
	return func(n AstNode) {
		if s, ok := n.(SlopNode); ok {
			s.setSlop(slop)
		}
	}
}

type slopNode struct {
	slop int
}

func (n *slopNode) getSlop() int {
	return n.slop
}

func (n *slopNode) setSlop(slop int) {
	n.slop = slop
}

// indicate whether does dsl query use filter context
type FilterCtxNode interface {
	isFilterCtx() bool
	setFilterCtx(filterCtx bool)
}

// indicate whether is node array data type
type ArrayTypeNode interface {
	isArrayType() bool
	setArrayType()
}

type opNode struct {
	filterCtx bool
}

func NewOpNode(filterCtx bool) *opNode {
	return &opNode{filterCtx: filterCtx}
}

func (n *opNode) isFilterCtx() bool {
	return n.filterCtx
}

func (n *opNode) setFilterCtx(filterCtx bool) {
	n.filterCtx = filterCtx
}

func (n *opNode) AstType() AstType {
	return OP_NODE_TYPE
}

// leaf node
type lfNode struct {
	filterCtx bool // whether node is filter ctx
}

func NewLfNode(opts ...func(*lfNode)) *lfNode {
	var n = &lfNode{
		filterCtx: false,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func WithFilterCtx(filterCtx bool) func(*lfNode) {
	return func(lf *lfNode) {
		lf.filterCtx = filterCtx
	}
}

func (n *lfNode) isFilterCtx() bool {
	return n.filterCtx
}

func (n *lfNode) setFilterCtx(filterCtx bool) {
	n.filterCtx = filterCtx
}

func (n *lfNode) AstType() AstType {
	return LEAF_NODE_TYPE
}

type fieldNode struct {
	lfNode
	field string
}

type valueNode struct {
	valueType
	value LeafValue
}

func NewValueNode(value LeafValue, valueType *valueType) *valueNode {
	return &valueNode{value: value, valueType: *valueType}
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
