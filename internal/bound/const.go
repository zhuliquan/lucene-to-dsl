package bound

var (
	Inf = &RangeValue{InfinityVal: "*"}
)

type BoundType uint16

const (
	UNKNOWN_BOUND_TYPE         BoundType = iota
	LEFT_EXCLUDE_RIGHT_INCLUDE BoundType = iota
	LEFT_EXCLUDE_RIGHT_EXCLUDE
	LEFT_INCLUDE_RIGHT_INCLUDE
	LEFT_INCLUDE_RIGHT_EXCLUDE
)
