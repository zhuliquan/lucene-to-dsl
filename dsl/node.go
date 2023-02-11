package dsl

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

func (n *expandsNode) setMaxExpands(maxExpands int) {
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
	setRewrite(RewriteType)
	getRewrite() RewriteType
}

func WithRewrite(rewrite RewriteType) func(AstNode) {
	return func(n AstNode) {
		if r, ok := n.(RewriteNode); ok {
			r.setRewrite(rewrite)
		}
	}
}

// rewrite node impl
type rewriteNode struct {
	rewrite RewriteType
}

func (r *rewriteNode) setRewrite(rewrite RewriteType) {
	r.rewrite = rewrite
}

func (r *rewriteNode) getRewrite() RewriteType {
	return r.rewrite
}

// interface node which has parameter of max determinized states
type StatesNode interface {
	getMaxDeterminizedStates() int
	setMaxDeterminizedStates(int)
}

func WithMaxDeterminizedStates(states int) func(AstNode) {
	return func(n AstNode) {
		if r, ok := n.(StatesNode); ok {
			r.setMaxDeterminizedStates(states)
		}
	}
}

type statesNode struct {
	maxDeterminizedStates int
}

func (s *statesNode) getMaxDeterminizedStates() int {
	return s.maxDeterminizedStates
}

func (s *statesNode) setMaxDeterminizedStates(states int) {
	s.maxDeterminizedStates = states
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
	setFilterCtx(filterCtx bool)
	getFilterCtx() bool
}

func WithFilterCtx(filterCtx bool) func(FilterCtxNode) {
	return func(n FilterCtxNode) {
		n.setFilterCtx(filterCtx)
	}
}

type filterCtxNode struct {
	filterCtx bool // whether node is filter ctx
}

func (n *filterCtxNode) getFilterCtx() bool {
	return n.filterCtx
}

func (n *filterCtxNode) setFilterCtx(filterCtx bool) {
	n.filterCtx = filterCtx
}

type opNode struct {
	filterCtxNode
}

func NewOpNode() *opNode {
	return &opNode{
		filterCtxNode: filterCtxNode{
			filterCtx: false,
		},
	}
}

func (n *opNode) AstType() AstType {
	return OP_NODE_TYPE
}

// leaf node
type lfNode struct {
	filterCtxNode
}

func NewLfNode() *lfNode {
	var n = &lfNode{
		filterCtxNode: filterCtxNode{
			filterCtx: false,
		},
	}
	return n
}

func (n *lfNode) AstType() AstType {
	return LEAF_NODE_TYPE
}

type fieldNode struct {
	lfNode
	field string
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

type ValueNode interface {
	getValue() LeafValue
	getVType() valueType
}

type valueNode struct {
	valueType
	value LeafValue
}

func (v *valueNode) getValue() LeafValue {
	return v.value
}

func (v *valueNode) getVType() valueType {
	return v.valueType
}

func NewValueNode(value LeafValue, valueType *valueType) *valueNode {
	return &valueNode{
		valueType: *valueType,
		value:     value,
	}
}

func (v *valueNode) toPrintValue() interface{} {
	return leafValueToPrintValue(v.value, v.mType)
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
	valueType
	rValue  LeafValue
	lValue  LeafValue
	rCmpSym CompareType
	lCmpSym CompareType
}

func NewRgNode(fieldNode *fieldNode, valueType *valueType, lValue, rValue LeafValue, lCmpSym, rCmpSym CompareType) *rgNode {
	return &rgNode{
		fieldNode: *fieldNode,
		valueType: *valueType,
		lValue:    lValue,
		rValue:    rValue,
		lCmpSym:   lCmpSym,
		rCmpSym:   rCmpSym,
	}
}

type PatternMatcher interface {
	Match([]byte) bool
}

type patternNode struct {
	matcher PatternMatcher
}

func (n *patternNode) Match(text []byte) bool {
	return n.matcher.Match(text)
}

