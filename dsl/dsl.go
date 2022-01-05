package dsl

import "encoding/json"

type DSL map[string]interface{}

func (d DSL) String() string {
	v, _ := json.Marshal(d)
	return string(v)
}

var EmptyDSL = DSL{}
