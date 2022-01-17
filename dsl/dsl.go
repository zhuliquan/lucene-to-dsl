package dsl

import (
	"encoding/json"
	"net"
	"time"
)

type DSL map[string]interface{}

func (d DSL) String() string {
	v, _ := json.Marshal(d)
	return string(v)
}

type DSLTermValue struct {
	BoolTerm   bool
	IntTerm    int64
	UintTerm   uint64
	IpTerm     net.IP
	IpCidrTerm *net.IPNet
	DateTerm   time.Time
	FloatTerm  float64
	StringTerm string
}

var EmptyDSL = DSL{}
