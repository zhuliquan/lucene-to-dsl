package operator

type PrefixOPType uint32

const (
	UNKNOWN_PREFIX_TYPE  PrefixOPType = 0
	SHOULD_PREFIX_TYPE   PrefixOPType = 1
	MUST_PREFIX_TYPE     PrefixOPType = 2
	MUST_NOT_PREFIX_TYPE PrefixOPType = 3
)

var prefixOPType_Values = map[PrefixOPType]string{
	UNKNOWN_PREFIX_TYPE:  "",
	SHOULD_PREFIX_TYPE:   " ",
	MUST_PREFIX_TYPE:     " +",
	MUST_NOT_PREFIX_TYPE: " -",
}

func (o PrefixOPType) String() string {
	return prefixOPType_Values[o]
}

type LogicOPType uint32

const (
	UNKNOWN_LOGIC_TYPE LogicOPType = 0
	AND_LOGIC_TYPE     LogicOPType = 1
	OR_LOGIC_TYPE      LogicOPType = 2
	NOT_LOGIC_TYPE     LogicOPType = 3
)

var LogicOPType_Values = map[LogicOPType]string{
	UNKNOWN_LOGIC_TYPE: "",
	AND_LOGIC_TYPE:     " AND ",
	OR_LOGIC_TYPE:      " OR ",
	NOT_LOGIC_TYPE:     "NOT ",
}

func (o LogicOPType) String() string {
	return LogicOPType_Values[o]
}
