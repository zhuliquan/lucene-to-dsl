package query

type Query interface {
	String() string
	// ToASTNode() (ASTNode, error)
}
