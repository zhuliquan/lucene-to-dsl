package dsl

import "fmt"

var (
	ErrIntersectNilNode = fmt.Errorf("failed to inverse nil node")
	ErrUnionJoinNilNode = fmt.Errorf("failed to union join two nil nodes")
	ErrInverseNilNode   = fmt.Errorf("failed to inverse nil node")
)
