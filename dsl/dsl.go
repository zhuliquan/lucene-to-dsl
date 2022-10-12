package dsl

import (
	"encoding/json"
	"reflect"
)

type DSL map[string]interface{}

func (d DSL) String() string {
	v, _ := json.Marshal(d)
	return string(v)
}

func addValueForDSL(d DSL, field string, value interface{}) {
	var v = reflect.ValueOf(value)
	if v.IsValid() && !v.IsZero() {
		d[field] = value
	}
}

type LeafValue interface{}

var EmptyDSL = DSL{}

var emptyDSL = EmptyDSL
